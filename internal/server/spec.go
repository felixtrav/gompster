package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/felixtrav/gompster/internal/handler"
	"github.com/felixtrav/gompster/internal/meta"
)

// --- OpenAPI 3.0 type definitions ---

type openAPISpec struct {
	OpenAPI string                     `json:"openapi"`
	Info    openAPIInfo                `json:"info"`
	Tags    []openAPITag               `json:"tags,omitempty"`
	Paths   map[string]openAPIPathItem `json:"paths"`
}

type openAPIInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type openAPITag struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

// openAPIPathItem maps HTTP method -> operation.
type openAPIPathItem map[string]openAPIOperation

type openAPIOperation struct {
	Summary     string                     `json:"summary,omitempty"`
	Description string                     `json:"description,omitempty"`
	Tags        []string                   `json:"tags,omitempty"`
	Parameters  []openAPIParameter         `json:"parameters,omitempty"`
	Responses   map[string]openAPIResponse `json:"responses"`
}

type openAPIParameter struct {
	Name        string        `json:"name"`
	In          string        `json:"in"`
	Description string        `json:"description,omitempty"`
	Required    bool          `json:"required"`
	Schema      openAPISchema `json:"schema"`
}

type openAPISchema struct {
	Type    string   `json:"type,omitempty"`
	Format  string   `json:"format,omitempty"`
	Default any      `json:"default,omitempty"`
	Minimum *float64 `json:"minimum,omitempty"`
	Maximum *float64 `json:"maximum,omitempty"`
}

type openAPIResponse struct {
	Description string                      `json:"description"`
	Content     map[string]openAPIMediaType `json:"content,omitempty"`
}

type openAPIMediaType struct {
	Schema openAPISchema `json:"schema"`
}

// generateSpec builds and marshals an OpenAPI 3.0.3 spec from the registered groups.
func generateSpec(groups []handler.Group) []byte {
	spec := openAPISpec{
		OpenAPI: "3.0.3",
		Info: openAPIInfo{
			Title:       meta.Name,
			Description: meta.Description + "\n\nRun locally: `docker run -p 80:80 felixtrav/gompster`",
			Version:     "1.0.0",
		},
		Paths: make(map[string]openAPIPathItem),
	}

	seenTags := map[string]bool{}
	for _, g := range groups {
		if !seenTags[g.Tag] {
			spec.Tags = append(spec.Tags, openAPITag{Name: g.Tag, Description: g.Description})
			seenTags[g.Tag] = true
		}

		for _, route := range g.Routes {
			item, ok := spec.Paths[route.Pattern]
			if !ok {
				item = make(openAPIPathItem)
			}

			tags := route.Tags
			if len(tags) == 0 {
				tags = []string{g.Tag}
			}
			op := openAPIOperation{
				Summary:     route.Summary,
				Description: route.Description,
				Tags:        tags,
				Responses:   make(map[string]openAPIResponse),
			}

			for _, p := range route.Params {
				if p.Name == "*" {
					// Skip wildcard catch-all path placeholders.
					continue
				}
				op.Parameters = append(op.Parameters, openAPIParameter{
					Name:        p.Name,
					In:          p.In,
					Description: p.Description,
					Required:    p.Required,
					Schema: openAPISchema{
						Type:    p.Schema.Type,
						Format:  p.Schema.Format,
						Default: p.Schema.Default,
						Minimum: p.Schema.Minimum,
						Maximum: p.Schema.Maximum,
					},
				})
			}

			for code, resp := range route.Responses {
				r := openAPIResponse{Description: resp.Description}
				if resp.ContentType != "" {
					r.Content = map[string]openAPIMediaType{
						resp.ContentType: {Schema: openAPISchema{Type: "object"}},
					}
				}
				op.Responses[fmt.Sprintf("%d", code)] = r
			}

			method := strings.ToLower(route.Method)
			// A blank Method means the route accepts any HTTP method (e.g. /anything).
			// Register the same operation under all five standard methods so Swagger
			// UI shows each one as a usable endpoint.
			if method == "" {
				for _, m := range []string{"get", "post", "put", "patch", "delete"} {
					item[m] = op
				}
			} else {
				item[method] = op
			}
			spec.Paths[route.Pattern] = item
		}
	}

	b, _ := json.MarshalIndent(spec, "", "  ")
	return b
}

// swaggerUIHTML serves Swagger UI (loaded from CDN) pointing at /openapi.json.
const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>gompster — API Reference</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    body { margin: 0; }
    .topbar { background-color: #1b1b2f !important; }
    .topbar-wrapper a span { display: none; }
    .topbar-wrapper::after {
      content: "gompster";
      color: #e2e2e2;
      font-size: 1.3rem;
      font-weight: 700;
      padding-left: 1.2rem;
    }
  </style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
  SwaggerUIBundle({
    url: "/openapi.json",
    dom_id: "#swagger-ui",
    deepLinking: true,
    tagsSorter: "alpha",
    tryItOutEnabled: true,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIBundle.SwaggerUIStandalonePreset],
    layout: "BaseLayout",
  });
</script>
</body>
</html>`

func swaggerUIHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, swaggerUIHTML)
}
