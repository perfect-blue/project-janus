package gateways

import (
	"fmt"
	"strings"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Println("this is test")
}

func TestRouter(t *testing.T) {
	trie := NewRouter()
	fmt.Println("1. ADDING ROUTES WITH METHODS AND QUERY RULES:")
	
	// Simple routes without queries
	trie.AddRoute(GET, "/", nil, "Home page")
	trie.AddRoute(GET, "/users", nil, "List all users")
	trie.AddRoute(POST, "/users", nil, "Create new user")
	
	// // Routes with query parameters
	trie.AddRoute(GET, "/users", map[string]QueryRule{
		"page":  {Required: false, AllowedValues: nil, Description: "Page number"},
		"limit": {Required: false, AllowedValues: []string{"10", "20", "50"}, Description: "Items per page"},
	}, "List users with pagination")
	
	trie.AddRoute(GET, "/search", map[string]QueryRule{
		"q":    {Required: true, AllowedValues: nil, Description: "Search query"},
		"type": {Required: false, AllowedValues: []string{"user", "post", "comment"}, Description: "Search type"},
	}, "Search endpoint")
	
	trie.AddRoute(GET, "/api/v1/posts", map[string]QueryRule{
		"status":   {Required: false, AllowedValues: []string{"draft", "published", "archived"}, Description: "Post status"},
		"author":   {Required: false, AllowedValues: nil, Description: "Author ID"},
		"category": {Required: true, AllowedValues: []string{"tech", "business", "lifestyle"}, Description: "Post category"},
	}, "Get posts with filters")
	
	// // Same path, different methods
	trie.AddRoute(PUT, "/api/v1/posts", nil, "Update post")
	trie.AddRoute(DELETE, "/api/v1/posts", map[string]QueryRule{
		"confirm": {Required: true, AllowedValues: []string{"yes"}, Description: "Confirmation"},
	}, "Delete post with confirmation")
	
	// Display the trie
	trie.PrintTrie()

	// // Test route matching
	fmt.Println("2. TESTING ROUTE MATCHING:")
	
	testCases := []struct {
		method HTTPMethod
		url    string
	}{
		{GET, "/"},
		{GET, "/users"},
		{POST, "/users"},
		{GET, "/users?page=1&limit=20"},
		{GET, "/users?page=1&limit=100"}, // Invalid limit
		{GET, "/search?q=golang"},
		{GET, "/search"}, // Missing required query
		{GET, "/search?q=golang&type=user"},
		{GET, "/search?q=golang&type=invalid"}, // Invalid type
		{GET, "/api/v1/posts?category=tech&status=published"},
		{GET, "/api/v1/posts?status=published"}, // Missing required category
		{PUT, "/api/v1/posts"},
		{DELETE, "/api/v1/posts?confirm=yes"},
		{DELETE, "/api/v1/posts"}, // Missing required confirmation
	}
	
	for _, testCase := range testCases {
		fmt.Printf("Testing: %s %s\n", testCase.method, testCase.url)
		result := trie.FindRoute(testCase.method, testCase.url)
		
		if result.Found {
			fmt.Printf("✅ Match found: %s\n", result.Route.Description)
			if len(result.QueryParams) > 0 {
				fmt.Printf("   Query params: %v\n", result.QueryParams)
			}
			if len(result.QueryErrors) > 0 {
				fmt.Printf("❌ Query errors: %v\n", result.QueryErrors)
			}
		} else {
			fmt.Printf("❌ No match found\n")
		}
		fmt.Println(strings.Repeat("-", 50))
	}
}