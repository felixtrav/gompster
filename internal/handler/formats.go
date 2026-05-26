package handler

import (
	"compress/flate"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/andybalholm/brotli"

	"github.com/felixtrav/gompster/internal/meta"
	"github.com/felixtrav/gompster/internal/samples"
)

// FormatsGroup returns routes that respond with various encodings and content types.
func FormatsGroup() Group {
	return Group{
		Tag:         "Response Formats",
		Description: "Endpoints that return responses in different encodings and content types.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/gzip",
				Handler:     gzipHandler,
				Summary:     "Gzip-compressed response",
				Description: "Returns a gzip-compressed JSON body with `Content-Encoding: gzip`.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "Gzip-compressed JSON", ContentType: "application/json"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/deflate",
				Handler:     deflateHandler,
				Summary:     "Deflate-compressed response",
				Description: "Returns a deflate-compressed JSON body with `Content-Encoding: deflate`.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "Deflate-compressed JSON", ContentType: "application/json"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/brotli",
				Handler:     brotliHandler,
				Summary:     "Brotli-compressed response",
				Description: "Returns a brotli-compressed JSON body with `Content-Encoding: br`.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "Brotli-compressed JSON", ContentType: "application/json"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/utf8",
				Handler:     utf8Handler,
				Summary:     "UTF-8 encoded response",
				Description: "Returns a UTF-8 encoded plain text body containing Unicode sample text.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "UTF-8 text", ContentType: "text/plain"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/html",
				Handler:     htmlFormatHandler,
				Summary:     "HTML response",
				Description: "Returns a simple HTML document.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "HTML document", ContentType: "text/html"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/json",
				Handler:     jsonFormatHandler,
				Summary:     "JSON response",
				Description: "Returns a sample JSON document.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "Sample JSON", ContentType: "application/json"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/xml",
				Handler:     xmlFormatHandler,
				Summary:     "XML response",
				Description: "Returns a sample XML document.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "Sample XML", ContentType: "application/xml"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/robots.txt",
				Handler:     robotsTxtHandler,
				Summary:     "Robots.txt",
				Description: "Returns a robots.txt that disallows /deny.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "robots.txt", ContentType: "text/plain"}},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/deny",
				Handler:     denyHandler,
				Summary:     "Denied by robots.txt",
				Description: "Returns a plain text page that is disallowed by /robots.txt.",
				Tags:        []string{"Response Formats"},
				Responses:   map[int]Response{200: {Description: "Denied content", ContentType: "text/plain"}},
			},
		},
	}
}

// --- compressed handlers ---

func gzipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "gzip")
	gz := gzip.NewWriter(w)
	writeProjectPayload(gz, "gzip", ClientIP(r))
	_ = gz.Close()
}

func deflateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "deflate")
	fw, _ := flate.NewWriter(w, flate.DefaultCompression)
	writeProjectPayload(fw, "deflate", ClientIP(r))
	_ = fw.Close()
}

func brotliHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Encoding", "br")
	bw := brotli.NewWriter(w)
	writeProjectPayload(bw, "brotli", ClientIP(r))
	_ = bw.Close()
}

// writeProjectPayload writes the gompster project JSON payload to any writer,
// annotated with the active content encoding and the caller's origin IP.
func writeProjectPayload(w interface{ Write([]byte) (int, error) }, encoding, origin string) {
	_, _ = fmt.Fprintf(w, `{
  "project": %q,
  "description": %q,
  "encoding": %q,
  "origin": %q,
  "features": [
    "Echo any HTTP method",
    "Inspect headers, IP, and User-Agent",
    "Generate UUIDs, random bytes, and delays",
    "Test redirects, auth, cookies, and streaming",
    "Multiple response formats and encodings",
    "Request forwarding with response capture",
    "OpenAPI spec and Swagger UI built-in"
  ],
  "links": {
    "spec": "/openapi.json",
    "swagger": "/"
  }
}`, meta.Name, meta.Description, encoding, origin)
}

// --- content type handlers ---

func utf8Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(samples.UTF8Demo)
}

func htmlFormatHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>gompster</title></head>
<body>
  <h1>gompster</h1>
  <p>` + meta.Description + `</p>
  <ul>
    <li>Echo any HTTP method</li>
    <li>Inspect headers, IP, and User-Agent</li>
    <li>Generate UUIDs, random bytes, and delays</li>
    <li>Test redirects, auth, cookies, and streaming</li>
    <li>Multiple response formats and encodings</li>
    <li>Request forwarding with response capture</li>
    <li>OpenAPI spec and Swagger UI built-in</li>
  </ul>
  <p><a href="/openapi.json">OpenAPI spec</a> &mdash; <a href="/">Swagger UI</a></p>
</body>
</html>`)
}

func jsonFormatHandler(w http.ResponseWriter, _ *http.Request) {
	JSON(w, http.StatusOK, map[string]any{
		"project":     meta.Name,
		"description": meta.Description,
		"features": []string{
			"Echo any HTTP method",
			"Inspect headers, IP, and User-Agent",
			"Generate UUIDs, random bytes, and delays",
			"Test redirects, auth, cookies, and streaming",
			"Multiple response formats and encodings",
			"Request forwarding with response capture",
			"OpenAPI spec and Swagger UI built-in",
		},
		"links": map[string]string{
			"spec":    "/openapi.json",
			"swagger": "/",
		},
	})
}

type xmlProject struct {
	XMLName     xml.Name     `xml:"project"`
	Name        string       `xml:"name,attr"`
	Description string       `xml:"description"`
	Features    []xmlFeature `xml:"features>feature"`
}

type xmlFeature struct {
	Name string `xml:",chardata"`
}

func xmlFormatHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprint(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	_ = enc.Encode(xmlProject{
		Name:        meta.Name,
		Description: meta.Description,
		Features: []xmlFeature{
			{"Echo any HTTP method"},
			{"Inspect headers, IP, and User-Agent"},
			{"Generate UUIDs, random bytes, and delays"},
			{"Test redirects, auth, cookies, and streaming"},
			{"Multiple response formats and encodings"},
			{"Request forwarding with response capture"},
			{"OpenAPI spec and Swagger UI built-in"},
		},
	})
}

// --- robots / deny ---

func robotsTxtHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, "User-agent: *\nDisallow: /deny\n")
}

func denyHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, `
                                  ....
                                .'' .'''
.                             .'   :
\\                          .:    :
 \\                        _:    :       ..----.._
  \\                    .:::.....:::.. .'         ''.
   \\                 .'  #-. .-######'     #        '.
    \\                 '.##'/ ' ################       :
     \\                  #####################         :
      \\               ..##.-.#### .''''###'.._        :
       \\             :--:########:            '.    .' :
        \\..__...--.. :--:#######.'   '.         '.     :
        :     :  : : '':'-:'':'::        .         '.  .'
        '---'''..: :    ':    '..'''.      '.        :'
           \\  :: : :     '      ''''''.     '.      .:
            \\ ::  : :     '            '.      '      :
             \\::   : :           ....' ..:       '     '.
              \\::  : :    .....####\\ .~~.:.             :
               \\':.:.:.:'#########.===. ~ |.'-.   . '''.. :
                \\    .'  ########## \ \ _.' '. '-.       '''.
                :\\  :     ########   \ \      '.  '-.        :
               :  \\'    '   #### :    \ \      :.    '-.      :
              :  .'\\   :'  :     :     \ \       :      '-.    :
             : .'  .\\  '  :      :     :\ \       :        '.   :
             ::   :  \\'  :.      :     : \ \      :          '. :
             ::. :    \\  : :      :    ;  \ \     :           '.:
              : ':    '\\ :  :     :     :  \:\     :        ..'
                 :    ' \\ :        :     ;  \|      :   .'''
                 '.   '  \\:                         :.''
                  .:..... \\:       :            ..''
                 '._____|'.\\......'''''''.:..'''
                            \\

YOU SHALL NOT PASS

This page has been denied by robots.txt.
`)
}
