package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/felixtrav/gompster/internal/samples/images"

	"github.com/go-chi/chi/v5"
)

// OriginalWork holds attribution for a work that a derivative was based on.
// Only required when the distributed file is a derivative of a separately-licensed original.
type OriginalWork struct {
	Title      string `json:"title"`
	Author     string `json:"author"`
	AuthorURL  string `json:"author_url,omitempty"`
	Source     string `json:"source"`
	License    string `json:"license"`
	LicenseURL string `json:"license_url"`
}

// ImageAttribution holds the full attribution fields for a single image.
// Edit the placeholder values in supportedImages with the real details for each image.
type ImageAttribution struct {
	Title         string        `json:"title"`
	Author        string        `json:"author"`
	AuthorURL     string        `json:"author_url,omitempty"`
	Source        string        `json:"source"`
	License       string        `json:"license"`
	LicenseURL    string        `json:"license_url"`
	Modifications string        `json:"modifications,omitempty"`
	OriginalWork  *OriginalWork `json:"original_work,omitempty"`
}

type imageType struct {
	contentType string
	data        func() []byte
	attribution ImageAttribution
}

// supportedImages is the single source of truth for formats, data, and attribution.
// Replace the placeholder strings with real values before distributing.
var supportedImages = map[string]imageType{
	"jpeg": {
		contentType: "image/jpeg",
		data:        func() []byte { return images.JPEG },
		attribution: ImageAttribution{
			Title:         "gopher.svg",
			Author:        "Takuya Ueda",
			AuthorURL:     "https://twitter.com/tenntenn",
			Source:        "https://github.com/golang-samples/gopher-vector/blob/master/gopher.svg",
			Modifications: "Converted from SVG to JPEG using ImageMagick",
			License:       "CC BY 3.0",
			LicenseURL:    "https://creativecommons.org/licenses/by/3.0/",
			OriginalWork: &OriginalWork{
				Title:      "Go gopher",
				Author:     "Renee French",
				AuthorURL:  "https://reneefrench.blogspot.com/",
				Source:     "https://go.dev/blog/gopher",
				License:    "CC BY 4.0",
				LicenseURL: "https://creativecommons.org/licenses/by/4.0/",
			},
		},
	},
	"png": {
		contentType: "image/png",
		data:        func() []byte { return images.PNG },
		attribution: ImageAttribution{
			Title:      "gopher.png",
			Author:     "Takuya Ueda",
			AuthorURL:  "https://twitter.com/tenntenn",
			Source:     "https://github.com/golang-samples/gopher-vector/blob/master/gopher.png",
			License:    "CC BY 3.0",
			LicenseURL: "https://creativecommons.org/licenses/by/3.0/",
			OriginalWork: &OriginalWork{
				Title:      "Go gopher",
				Author:     "Renee French",
				AuthorURL:  "https://reneefrench.blogspot.com/",
				Source:     "https://go.dev/blog/gopher",
				License:    "CC BY 4.0",
				LicenseURL: "https://creativecommons.org/licenses/by/4.0/",
			},
		},
	},
	"svg": {
		contentType: "image/svg+xml",
		data:        func() []byte { return images.SVG },
		attribution: ImageAttribution{
			Title:      "gopher.svg",
			Author:     "Takuya Ueda",
			AuthorURL:  "https://twitter.com/tenntenn",
			Source:     "https://github.com/golang-samples/gopher-vector/blob/master/gopher.svg",
			License:    "CC BY 3.0",
			LicenseURL: "https://creativecommons.org/licenses/by/3.0/",
			OriginalWork: &OriginalWork{
				Title:      "Go gopher",
				Author:     "Renee French",
				AuthorURL:  "https://reneefrench.blogspot.com/",
				Source:     "https://go.dev/blog/gopher",
				License:    "CC BY 4.0",
				LicenseURL: "https://creativecommons.org/licenses/by/4.0/",
			},
		},
	},
	"webp": {
		contentType: "image/webp",
		data:        func() []byte { return images.WebP },
		attribution: ImageAttribution{
			Title:         "gopher.svg",
			Author:        "Takuya Ueda",
			AuthorURL:     "https://twitter.com/tenntenn",
			Source:        "https://github.com/golang-samples/gopher-vector/blob/master/gopher.svg",
			Modifications: "Converted from SVG to WebP using ImageMagick",
			License:       "CC BY 3.0",
			LicenseURL:    "https://creativecommons.org/licenses/by/3.0/",
			OriginalWork: &OriginalWork{
				Title:      "Go gopher",
				Author:     "Renee French",
				AuthorURL:  "https://reneefrench.blogspot.com/",
				Source:     "https://go.dev/blog/gopher",
				License:    "CC BY 4.0",
				LicenseURL: "https://creativecommons.org/licenses/by/4.0/",
			},
		},
	},
	"avif": {
		contentType: "image/avif",
		data:        func() []byte { return images.AVIF },
		attribution: ImageAttribution{
			Title:         "gopher.svg",
			Author:        "Takuya Ueda",
			AuthorURL:     "https://twitter.com/tenntenn",
			Source:        "https://github.com/golang-samples/gopher-vector/blob/master/gopher.svg",
			Modifications: "Converted from SVG to AVIF using ImageMagick",
			License:       "CC BY 3.0",
			LicenseURL:    "https://creativecommons.org/licenses/by/3.0/",
			OriginalWork: &OriginalWork{
				Title:      "Go gopher",
				Author:     "Renee French",
				AuthorURL:  "https://reneefrench.blogspot.com/",
				Source:     "https://go.dev/blog/gopher",
				License:    "CC BY 4.0",
				LicenseURL: "https://creativecommons.org/licenses/by/4.0/",
			},
		},
	},
	"gif": {
		contentType: "image/gif",
		data:        func() []byte { return images.GIF },
		attribution: ImageAttribution{
			Title:      "gopher-dance-long-3x.gif",
			Author:     "Egon Elbre",
			AuthorURL:  "https://egonelbre.com/",
			Source:     "https://github.com/egonelbre/gophers/blob/master/.thumb/animation/gopher-dance-long-3x.gif",
			License:    "CC0 1.0",
			LicenseURL: "https://creativecommons.org/publicdomain/zero/1.0/",
			OriginalWork: &OriginalWork{
				Title:      "Go gopher",
				Author:     "Renee French",
				AuthorURL:  "https://reneefrench.blogspot.com/",
				Source:     "https://go.dev/blog/gopher",
				License:    "CC BY 4.0",
				LicenseURL: "https://creativecommons.org/licenses/by/4.0/",
			},
		},
	},
	"apng": {
		contentType: "image/apng",
		data:        func() []byte { return images.APNG },
		attribution: ImageAttribution{
			Title:         "gopher-dance-long-3x.gif",
			Author:        "Egon Elbre",
			AuthorURL:     "https://egonelbre.com/",
			Source:        "https://github.com/egonelbre/gophers/blob/master/.thumb/animation/gopher-dance-long-3x.gif",
			License:       "CC0 1.0",
			LicenseURL:    "https://creativecommons.org/publicdomain/zero/1.0/",
			Modifications: "Converted from GIF to APNG using ImageMagick",
			OriginalWork: &OriginalWork{
				Title:      "Go gopher",
				Author:     "Renee French",
				AuthorURL:  "https://reneefrench.blogspot.com/",
				Source:     "https://go.dev/blog/gopher",
				License:    "CC BY 4.0",
				LicenseURL: "https://creativecommons.org/licenses/by/4.0/",
			},
		},
	},
}

func ImageGroup() Group {
	return Group{
		Tag:         "Images",
		Description: "Returns sample images in various formats.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/image",
				Summary:     "Image content negotiation",
				Description: "Returns a sample image based on the `Accept` header. Supports `image/jpeg`, `image/png`, `image/svg+xml`, `image/webp`, `image/avif`, `image/gif`, and `image/apng`. Falls back to PNG for `*/*`. Returns `406` if no supported type is acceptable. Attribution details are returned in response headers.",
				Responses: map[int]Response{
					200: {Description: "Image in the negotiated format"},
					406: {Description: "None of the accepted types are supported", ContentType: "application/json"},
				},
				Handler: imageNegotiateHandler,
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/image/attribution",
				Summary:     "Image attributions",
				Description: "Returns attribution details for all sample images.",
				Responses: map[int]Response{
					200: {Description: "Attribution map keyed by image type", ContentType: "application/json"},
				},
				Handler: imageAttributionHandler,
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/image/{type}",
				Summary:     "Image by type",
				Description: "Returns a sample image of the requested type. Valid values: `jpeg`, `png`, `svg`, `webp`, `avif`, `gif`, `apng`. Attribution details are returned in response headers.",
				Params: []Param{
					{Name: "type", In: "path", Required: true, Description: "Image format (jpeg, png, svg, webp, avif, gif, apng)", Schema: Schema{Type: "string"}},
				},
				Responses: map[int]Response{
					200: {Description: "Requested image"},
					404: {Description: "Unknown image type", ContentType: "application/json"},
				},
				Handler: imageTypeHandler,
			},
		},
	}
}

func imageNegotiateHandler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")

	// Preference order for wildcard or empty Accept.
	preferred := []string{"png", "jpeg", "webp", "avif", "gif", "apng", "svg"}

	if accept == "" || strings.Contains(accept, "*/*") {
		serveImage(w, supportedImages["png"])
		return
	}

	for _, pref := range preferred {
		img := supportedImages[pref]
		if strings.Contains(accept, img.contentType) {
			serveImage(w, img)
			return
		}
	}

	JSON(w, http.StatusNotAcceptable, map[string]any{
		"error":     "not acceptable",
		"supported": []string{"image/jpeg", "image/png", "image/svg+xml", "image/webp", "image/avif", "image/gif", "image/apng"},
	})
}

func imageTypeHandler(w http.ResponseWriter, r *http.Request) {
	typeName := strings.ToLower(chi.URLParam(r, "type"))
	img, ok := supportedImages[typeName]
	if !ok {
		JSON(w, http.StatusNotFound, map[string]any{
			"error":     "unknown image type",
			"supported": []string{"jpeg", "png", "svg", "webp", "avif", "gif", "apng"},
		})
		return
	}
	serveImage(w, img)
}

func imageAttributionHandler(w http.ResponseWriter, r *http.Request) {
	out := make(map[string]ImageAttribution, len(supportedImages))
	for k, v := range supportedImages {
		out[k] = v.attribution
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(out)
}

func serveImage(w http.ResponseWriter, img imageType) {
	a := img.attribution
	w.Header().Set("Content-Type", img.contentType)
	w.Header().Set("X-Attribution-Title", a.Title)
	w.Header().Set("X-Attribution-Author", a.Author)
	if a.AuthorURL != "" {
		w.Header().Set("X-Attribution-Author-URL", a.AuthorURL)
	}
	w.Header().Set("X-Attribution-Source", a.Source)
	w.Header().Set("X-Attribution-License", fmt.Sprintf("%s <%s>", a.License, a.LicenseURL))
	if a.Modifications != "" {
		w.Header().Set("X-Attribution-Modifications", a.Modifications)
	}
	if ow := a.OriginalWork; ow != nil {
		w.Header().Set("X-Attribution-Original-Title", ow.Title)
		w.Header().Set("X-Attribution-Original-Author", ow.Author)
		if ow.AuthorURL != "" {
			w.Header().Set("X-Attribution-Original-Author-URL", ow.AuthorURL)
		}
		w.Header().Set("X-Attribution-Original-Source", ow.Source)
		w.Header().Set("X-Attribution-Original-License", fmt.Sprintf("%s <%s>", ow.License, ow.LicenseURL))
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(img.data())
}
