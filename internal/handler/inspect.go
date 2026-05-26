package handler

import (
	"net/http"
	"strings"
)

// InspectGroup returns routes for inspecting specific parts of the request.
func InspectGroup() Group {
	return Group{
		Tag:         "Request Inspection",
		Description: "Endpoints that return specific aspects of the incoming request.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/headers",
				Handler:     headersHandler,
				Summary:     "Return request headers",
				Description: "Returns a map of all HTTP headers sent by the client.",
				Tags:        []string{"Request Inspection"},
				Responses:   map[int]Response{200: {Description: "Header map", ContentType: "application/json"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/ip",
				Handler:     ipHandler,
				Summary:     "Return origin IP",
				Description: "Returns the IP address of the requester, respecting X-Forwarded-For.",
				Tags:        []string{"Request Inspection"},
				Responses:   map[int]Response{200: {Description: "Origin IP", ContentType: "application/json"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/ip/raw",
				Handler:     rawIPHandler,
				Summary:     "Return raw remote IP",
				Description: "Returns the direct TCP peer address, ignoring X-Forwarded-For and X-Real-Ip headers.",
				Tags:        []string{"Request Inspection"},
				Responses:   map[int]Response{200: {Description: "Raw remote IP", ContentType: "application/json"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/user-agent",
				Handler:     userAgentHandler,
				Summary:     "Return User-Agent",
				Description: "Returns the User-Agent string from the request headers.",
				Tags:        []string{"Request Inspection"},
				Responses:   map[int]Response{200: {Description: "User-Agent string", ContentType: "application/json"}},
			},
		},
	}
}

func headersHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]any{
		"headers": headersToMap(r.Header),
	})
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{
		"origin": ClientIP(r),
	})
}

// rawIPHandler intentionally bypasses ClientIP: it strips only the port from
// r.RemoteAddr without consulting X-Forwarded-For or X-Real-Ip. Note that
// middleware.RealIP (applied globally in server.go) rewrites r.RemoteAddr
// with the first value from X-Forwarded-For before any handler runs, so what
// this returns is the last-hop address as seen by the proxy layer — not
// necessarily the true TCP peer of this process.
func rawIPHandler(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr
	if i := strings.LastIndex(addr, ":"); i >= 0 {
		addr = addr[:i]
	}
	JSON(w, http.StatusOK, map[string]string{"origin": addr})
}

func userAgentHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{
		"user-agent": r.UserAgent(),
	})
}
