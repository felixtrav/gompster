package handler

// AnythingGroup returns routes that accept any HTTP method at /anything and sub-paths.
func AnythingGroup() Group {
	jsonResp := map[int]Response{
		200: {Description: "Request data echoed as JSON", ContentType: "application/json"},
	}

	return Group{
		Tag:         "Anything",
		Description: "Catch-all endpoints that echo the full request for any HTTP method.",
		Routes: []Route{
			{
				Method:      "",
				Pattern:     "/anything",
				Handler:     echoHandler,
				Summary:     "Echo any request",
				Description: "Returns the full request regardless of HTTP method.",
				Tags:        []string{"Anything"},
				Responses:   jsonResp,
			},
			{
				Method:      "",
				Pattern:     "/anything/*",
				Handler:     echoHandler,
				Summary:     "Echo any request at any sub-path",
				Description: "Returns the full request for any method at any sub-path.",
				Tags:        []string{"Anything"},
				Params: []Param{
					{Name: "path", In: "path", Description: "Arbitrary sub-path", Required: false, Schema: Schema{Type: "string"}},
				},
				Responses: jsonResp,
			},
		},
	}
}
