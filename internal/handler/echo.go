package handler

import "net/http"

// EchoGroup returns routes that echo the full HTTP request back to the caller.
func EchoGroup() Group {
	jsonResp := map[int]Response{
		200: {Description: "Request data echoed as JSON", ContentType: "application/json"},
	}

	return Group{
		Tag:         "HTTP Methods",
		Description: "Endpoints that echo back the full request for each HTTP method.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/get",
				Handler:     echoHandler,
				Summary:     "Echo a GET request",
				Description: "Returns query parameters, headers, and request metadata.",
				Tags:        []string{"HTTP Methods"},
				Responses:   jsonResp,
			},
			{
				Method:      http.MethodPost,
				Pattern:     "/post",
				Handler:     echoHandler,
				Summary:     "Echo a POST request",
				Description: "Returns the request body, form data, headers, and metadata.",
				Tags:        []string{"HTTP Methods"},
				Responses:   jsonResp,
			},
			{
				Method:      http.MethodPut,
				Pattern:     "/put",
				Handler:     echoHandler,
				Summary:     "Echo a PUT request",
				Description: "Returns the request body, headers, and metadata.",
				Tags:        []string{"HTTP Methods"},
				Responses:   jsonResp,
			},
			{
				Method:      http.MethodPatch,
				Pattern:     "/patch",
				Handler:     echoHandler,
				Summary:     "Echo a PATCH request",
				Description: "Returns the request body, headers, and metadata.",
				Tags:        []string{"HTTP Methods"},
				Responses:   jsonResp,
			},
			{
				Method:      http.MethodDelete,
				Pattern:     "/delete",
				Handler:     echoHandler,
				Summary:     "Echo a DELETE request",
				Description: "Returns query parameters, headers, and metadata.",
				Tags:        []string{"HTTP Methods"},
				Responses:   jsonResp,
			},
		},
	}
}

func echoHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, CaptureRequest(r))
}
