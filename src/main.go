package main

import (
	"fmt"
	"janus/src/gateways"
	"net/http"
)

func main() {
	cfg, err := gateways.LoadConfig("config/config.dev.json")
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	server := http.Server{
		Addr: fmt.Sprintf("localhost:%d", cfg.Port),
		Handler: setupEndpoints(cfg),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok"))
}

func setupEndpoints(cfg gateways.Config) *http.ServeMux {
	router := gateways.NewRouter()
	for _, endpoint := range cfg.Endpoints {
		queryRules := make(map[string]gateways.QueryRule)
		for _, queryParam := range endpoint.Query {
			queryRules[queryParam] = gateways.QueryRule{
				Required:    false, // Set based on your requirements
				AllowedValues: []string{}, // Set based on your requirements
				Description: fmt.Sprintf("Query parameter for %s", queryParam),
			}
		}
		
		router.AddRoute(
			endpoint.Method,
			endpoint.Endpoint,
			queryRules,
			fmt.Sprintf(
				"Endpoint for %s %s", endpoint.Method, endpoint.Endpoint,
			),
		)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Skip health check and root path
		if r.URL.Path == "/health" || r.URL.Path == "/" {
			return
		}
		
		// Convert string method to HTTPMethod
		var method gateways.HTTPMethod
		switch r.Method {
		case "GET":
			method = gateways.GET
		case "POST":
			method = gateways.POST
		case "PUT":
			method = gateways.PUT
		case "DELETE":
			method = gateways.DELETE
		default:
			http.Error(w, "Method not supported", http.StatusMethodNotAllowed)
			return
		}
		
		// Use trie to find the route
		fullURL := r.URL.String()
		matchResult := router.FindRoute(method, fullURL)
		
		if !matchResult.Found {
			http.Error(w, "Route not found", http.StatusNotFound)
			return
		}
		
		// Find the corresponding config endpoint to get backend info
		var configEndpoint *gateways.Endpoint
		for _, ep := range cfg.Endpoints {
			if ep.Endpoint == matchResult.Route.Path && ep.Method == method {
				configEndpoint = &ep
				break
			}
		}
		
		if configEndpoint == nil {
			http.Error(w, "Config not found for route", http.StatusInternalServerError)
			return
		}
		
		// Display route information and backend services
		fmt.Fprintf(w, "Method: %s\n", method)
		fmt.Fprintf(w, "Path: %s\n", matchResult.Route.Path)
		
		// Show query validation errors if any
		if len(matchResult.QueryErrors) > 0 {
			fmt.Fprintf(w, "\nQuery Validation Errors:\n")
			for _, err := range matchResult.QueryErrors {
				fmt.Fprintf(w, "  - %s\n", err)
			}
		}
		
		// Show backend services
		fmt.Fprintf(w, "\nBackend Services:\n")
		for i, backend := range configEndpoint.Backend {
			fmt.Fprintf(w, "  Backend %d:\n", i+1)
			fmt.Fprintf(w, "    URL Pattern: %s\n", backend.URLPattern)
			fmt.Fprintf(w, "    Hosts: %v\n", backend.Host)
			if backend.Port > 0 {
				fmt.Fprintf(w, "    Port: %d\n", backend.Port)
			}
		}
		
		// Show query parameters received
		if len(matchResult.QueryParams) > 0 {
			fmt.Fprintf(w, "\nReceived Query Parameters:\n")
			for key, value := range matchResult.QueryParams {
				fmt.Fprintf(w, "  %s: %s\n", key, value)
			}
		}
		
		// Show query mappings if any
		if len(configEndpoint.QueryMapping) > 0 {
			fmt.Fprintf(w, "\nQuery Parameter Mappings:\n")
			for from, to := range configEndpoint.QueryMapping {
				fmt.Fprintf(w, "  %s -> %s\n", from, to)
			}
		}
	})

	return mux
}