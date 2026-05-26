// Package handler contains HTTP handlers and route registration types.
package handler

import "net/http"

// Schema describes the type and constraints of a parameter value.
type Schema struct {
	Type    string
	Format  string
	Default any
	Minimum *float64
	Maximum *float64
}

// Param describes a URL path, query string, or header parameter.
type Param struct {
	Name        string
	In          string // "path", "query", "header"
	Description string
	Required    bool
	Schema      Schema
}

// Response describes a single HTTP response variant for a route.
type Response struct {
	Description string
	ContentType string // empty means no body schema
}

// Route is a single HTTP endpoint with full metadata for the OpenAPI spec.
type Route struct {
	Method      string // empty means all methods
	Pattern     string
	Handler     http.HandlerFunc
	Summary     string
	Description string
	Tags        []string
	Params      []Param
	Responses   map[int]Response
}

// Group is a named, tagged collection of related routes.
// Adding a new Group to the server is all that is needed to extend the API.
type Group struct {
	Tag         string
	Description string
	Routes      []Route
}

func f64p(f float64) *float64 { return &f }
