package gateways

import (
	"fmt"
	"net/url"
	"strings"
)

type HTTPMethod string
const(
	GET HTTPMethod = "GET"
	POST HTTPMethod = "POST"
	PUT HTTPMethod = "PUT"
	DELETE HTTPMethod = "DELETE"
)

type Route struct {
	Method 			HTTPMethod
	Path			string
	QueryRules		map[string]QueryRule
	Description		string
}

type QueryRule struct {
	Required		bool
	AllowedValues	[]string
	Description		string
}

type RouteNode struct {
	children 	map[string]*RouteNode 	// child node for each path segment
	routes		map[HTTPMethod]*Route	// Routes for each HTTP method at this node
	path		string					// The path to this node
}

type Router struct {
	root *RouteNode
}

type MatchResult struct {
	Found        bool
	Route        *Route
	QueryParams  map[string]string // Parsed query parameters
	QueryErrors  []string          // Query validation errors
}

func NewRouteNode() *RouteNode {
	return &RouteNode{
		children: make(map[string]*RouteNode),
		routes: make(map[HTTPMethod]*Route),
	}
}

func NewRouter() *Router {
	return &Router{
		root: &RouteNode{
			children: make(map[string]*RouteNode),
			routes: make(map[HTTPMethod]*Route),
			path: "/",
		},
	}
}

func (t *Router) AddRoute(
	method HTTPMethod,
	path string,
	queryRules map[string]QueryRule,
	description string,
) {
	cleanPath := strings.Trim(path, "/")
	var segments []string
	if cleanPath != ""{
		segments = strings.Split(cleanPath, "/")
	}

	current := t.root
	currentPath := ""

	for _, segment := range segments {
		currentPath += "/" + segment
		if current.children[segment] == nil {
			current.children[segment] = NewRouteNode()
			current.children[segment].path = currentPath
		}
		current = current.children[segment]
	}

	route := &Route{
		Method:		method,
		Path:		path,
		QueryRules: queryRules,
		Description: description,
	}

	current.routes[method] = route
}

func (t *Router) FindRoute(method HTTPMethod, fullUrl string) *MatchResult {
	parsedURL, err := url.Parse(fullUrl)
	if err != nil {
		return &MatchResult{Found: false}
	}

	path := parsedURL.Path
	queryParams := make(map[string]string)

	for key, values := range parsedURL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0] // Take first value if multiple
		}
	}

	cleanPath := strings.Trim(path, "/")
	var segments []string
	if cleanPath != "" {
		segments = strings.Split(cleanPath, "/")
	}

	current := t.root
	for _, segment := range segments {
		if current.children[segment] == nil {
			return &MatchResult{Found: false}
		}
		current = current.children[segment]
	}

	route, methodExsists := current.routes[method]
	if !methodExsists {
		return &MatchResult{Found: false}
	}

	queryErrors := t.validateQueryParams(route.QueryRules, queryParams)
	result := &MatchResult{
		Found:			true,
		Route:			route,
		QueryParams:	queryParams,
		QueryErrors:	queryErrors,
	}

	return result
}

func (t *Router) validateQueryParams(rules map[string]QueryRule, params map[string]string) []string{
	var errors []string
	for paramName, rule := range rules {
		value, exists := params[paramName]
		if rule.Required && !exists {
			error := fmt.Sprintf("Required query parameter '%s' is missing", paramName)
			errors = append(errors, error)
			continue
		}

		if exists && value != "" && len(rule.AllowedValues) > 0 {
			allowed := false
			for _, allowedValue := range rule.AllowedValues {
				if value == allowedValue {
					break
				}
			}

			if !allowed {
				error := fmt.Sprintf("Parameter '%s' has invalid value '%s'. Allowed values: %v", 
					paramName, value, rule.AllowedValues)
				errors = append(errors, error)
			}
		}
	}

	return errors
}

func (t *Router) printNode(node *RouteNode, segment string, level int) {
	indent := strings.Repeat("  ", level)
	fmt.Printf("%s%s", indent, segment)
	
	if len(node.routes) > 0 {
		fmt.Printf(" -> Methods: ")
		for method, route := range node.routes {
			fmt.Printf("%s", method)
			if len(route.QueryRules) > 0 {
				fmt.Printf("(+queries)")
			}
			fmt.Printf(" ")
		}
	}
	fmt.Println()
	
	// Show query rules if any
	for method, route := range node.routes {
		if len(route.QueryRules) > 0 {
			fmt.Printf("%s  [%s Query Rules]:\n", indent, method)
			for paramName, rule := range route.QueryRules {
				fmt.Printf("%s    %s: required=%t", indent, paramName, rule.Required)
				if len(rule.AllowedValues) > 0 {
					fmt.Printf(", allowed=%v", rule.AllowedValues)
				}
				if rule.Description != "" {
					fmt.Printf(" (%s)", rule.Description)
				}
				fmt.Println()
			}
		}
	}
	
	// Print children
	for childSegment, childNode := range node.children {
		t.printNode(childNode, childSegment, level+1)
	}
}


func (t *Router) PrintTrie() {
	fmt.Println("=== Trie Structure ===")
	t.printNode(t.root, "", 0)
	fmt.Println()
}
