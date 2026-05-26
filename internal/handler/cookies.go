package handler

import "net/http"

// CookiesGroup returns routes for reading, setting, and deleting cookies.
func CookiesGroup() Group {
	return Group{
		Tag:         "Cookies",
		Description: "Endpoints for inspecting and manipulating HTTP cookies.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/cookies",
				Handler:     cookiesHandler,
				Summary:     "Return cookies",
				Description: "Returns all cookies present in the request.",
				Tags:        []string{"Cookies"},
				Responses:   map[int]Response{200: {Description: "Cookie map", ContentType: "application/json"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/cookies/set",
				Handler:     cookiesSetHandler,
				Summary:     "Set cookies",
				Description: "Sets one or more cookies from query parameters and redirects to /cookies.",
				Tags:        []string{"Cookies"},
				Params: []Param{
					{Name: "name", In: "query", Description: "Cookie name=value pairs to set", Required: false, Schema: Schema{Type: "string"}},
				},
				Responses: map[int]Response{302: {Description: "Redirect to /cookies"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/cookies/delete",
				Handler:     cookiesDeleteHandler,
				Summary:     "Delete cookies",
				Description: "Expires named cookies (from query params) and redirects to /cookies.",
				Tags:        []string{"Cookies"},
				Params: []Param{
					{Name: "name", In: "query", Description: "Cookie names to delete", Required: false, Schema: Schema{Type: "string"}},
				},
				Responses: map[int]Response{302: {Description: "Redirect to /cookies"}},
			},
		},
	}
}

func cookiesHandler(w http.ResponseWriter, r *http.Request) {
	cookies := make(map[string]string)
	for _, c := range r.Cookies() {
		cookies[c.Name] = c.Value
	}
	JSON(w, http.StatusOK, map[string]any{"cookies": cookies})
}

func cookiesSetHandler(w http.ResponseWriter, r *http.Request) {
	for name, vals := range r.URL.Query() {
		http.SetCookie(w, &http.Cookie{
			Name:  name,
			Value: vals[0],
			Path:  "/",
		})
	}
	http.Redirect(w, r, "/cookies", http.StatusFound)
}

func cookiesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	for name := range r.URL.Query() {
		// MaxAge -1 signals net/http to emit "Max-Age=0" on the wire, which
		// instructs the browser to delete the cookie immediately.
		http.SetCookie(w, &http.Cookie{
			Name:   name,
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
	}
	http.Redirect(w, r, "/cookies", http.StatusFound)
}
