package handler

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

// uuidEpoch is the number of 100-ns intervals between the UUID epoch
// (15 Oct 1582) and the Unix epoch (1 Jan 1970).
const uuidEpoch = 0x01b21dd213814000

// Standard v3/v5 namespaces defined in RFC 4122 §4.3.
var uuidNamespaces = map[string][16]byte{
	"dns":  {0x6b, 0xa7, 0xb8, 0x10, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
	"url":  {0x6b, 0xa7, 0xb8, 0x11, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
	"oid":  {0x6b, 0xa7, 0xb8, 0x12, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
	"x500": {0x6b, 0xa7, 0xb8, 0x14, 0x9d, 0xad, 0x11, 0xd1, 0x80, 0xb4, 0x00, 0xc0, 0x4f, 0xd4, 0x30, 0xc8},
}

// UUIDGroup returns routes for generating UUIDs of each RFC 4122 / draft version.
func UUIDGroup() Group {
	jsonResp := map[int]Response{200: {Description: "Generated UUID", ContentType: "application/json"}}
	nameParams := []Param{
		{Name: "namespace", In: "query", Description: "Namespace: dns, url, oid, x500 (default: dns)", Required: false, Schema: Schema{Type: "string", Default: "dns"}},
		{Name: "name", In: "query", Description: "Name to hash", Required: true, Schema: Schema{Type: "string"}},
	}

	return Group{
		Tag:         "UUIDs",
		Description: "Endpoints for generating UUIDs of specific versions.",
		Routes: []Route{
			{
				Method:      http.MethodGet,
				Pattern:     "/uuid",
				Handler:     uuidV4Handler,
				Summary:     "UUID v4 (alias)",
				Description: "Alias for /uuid/v4. Generates a cryptographically random UUID v4.",
				Tags:        []string{"UUIDs"},
				Responses:   jsonResp,
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/uuid/v1",
				Handler:     uuidV1Handler,
				Summary:     "UUID v1 — time-based",
				Description: "Generates a time-based UUID v1 using the current UTC time and a random node ID.",
				Tags:        []string{"UUIDs"},
				Responses:   jsonResp,
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/uuid/v3",
				Handler:     uuidV3Handler,
				Summary:     "UUID v3 — name-based (MD5)",
				Description: "Generates a name-based UUID v3 by hashing the namespace UUID and `name` with MD5.",
				Tags:        []string{"UUIDs"},
				Params:      nameParams,
				Responses: map[int]Response{
					200: {Description: "UUID v3", ContentType: "application/json"},
					400: {Description: "Missing name or unknown namespace"},
				},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/uuid/v4",
				Handler:     uuidV4Handler,
				Summary:     "UUID v4 — random",
				Description: "Generates a cryptographically random UUID v4 per RFC 4122 §4.4.",
				Tags:        []string{"UUIDs"},
				Responses:   jsonResp,
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/uuid/v5",
				Handler:     uuidV5Handler,
				Summary:     "UUID v5 — name-based (SHA-1)",
				Description: "Generates a name-based UUID v5 by hashing the namespace UUID and `name` with SHA-1.",
				Tags:        []string{"UUIDs"},
				Params:      nameParams,
				Responses: map[int]Response{
					200: {Description: "UUID v5", ContentType: "application/json"},
					400: {Description: "Missing name or unknown namespace"},
				},
			},
			{
				Method:      http.MethodGet,
				Pattern:     "/uuid/v7",
				Handler:     uuidV7Handler,
				Summary:     "UUID v7 — time-ordered",
				Description: "Generates a UUID v7 with a Unix millisecond timestamp prefix for natural sort order (draft-peabody-dispatch-new-uuid-format).",
				Tags:        []string{"UUIDs"},
				Responses:   jsonResp,
			},
		},
	}
}

// formatUUID formats a 16-byte array as a canonical UUID string.
func formatUUID(b [16]byte) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		hex.EncodeToString(b[0:4]),
		hex.EncodeToString(b[4:6]),
		hex.EncodeToString(b[6:8]),
		hex.EncodeToString(b[8:10]),
		hex.EncodeToString(b[10:]),
	)
}

func uuidResponse(w http.ResponseWriter, version int, id string) {
	JSON(w, http.StatusOK, map[string]any{
		"version": version,
		"uuid":    id,
	})
}

// v1: time-based with random node (RFC 4122 §4.2).
func uuidV1Handler(w http.ResponseWriter, _ *http.Request) {
	var b [16]byte
	now := uint64(time.Now().UTC().UnixNano()/100) + uuidEpoch

	binary.BigEndian.PutUint32(b[0:4], uint32(now&0xffffffff))        // time_low
	binary.BigEndian.PutUint16(b[4:6], uint16((now>>32)&0xffff))      // time_mid
	binary.BigEndian.PutUint16(b[6:8], uint16((now>>48)&0x0fff)|0x1000) // time_hi_and_version

	_, _ = rand.Read(b[8:])    // clock_seq + node (random)
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant

	uuidResponse(w, 1, formatUUID(b))
}

// v3: MD5 name-based (RFC 4122 §4.3).
func uuidV3Handler(w http.ResponseWriter, r *http.Request) {
	ns, name, ok := resolveNameParams(w, r)
	if !ok {
		return
	}
	h := md5.New()
	h.Write(ns[:])
	h.Write([]byte(name))
	sum := h.Sum(nil)

	var b [16]byte
	copy(b[:], sum)
	b[6] = (b[6] & 0x0f) | 0x30 // version 3
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant

	uuidResponse(w, 3, formatUUID(b))
}

// v4: random (RFC 4122 §4.4).
func uuidV4Handler(w http.ResponseWriter, _ *http.Request) {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant
	uuidResponse(w, 4, formatUUID(b))
}

// v5: SHA-1 name-based (RFC 4122 §4.3).
func uuidV5Handler(w http.ResponseWriter, r *http.Request) {
	ns, name, ok := resolveNameParams(w, r)
	if !ok {
		return
	}
	h := sha1.New()
	h.Write(ns[:])
	h.Write([]byte(name))
	sum := h.Sum(nil)

	var b [16]byte
	copy(b[:], sum[:16])
	b[6] = (b[6] & 0x0f) | 0x50 // version 5
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant

	uuidResponse(w, 5, formatUUID(b))
}

// v7: Unix-ms time-ordered (draft-peabody-dispatch-new-uuid-format-04 §5.2).
func uuidV7Handler(w http.ResponseWriter, _ *http.Request) {
	var b [16]byte
	_, _ = rand.Read(b[:])

	ms := uint64(time.Now().UnixMilli())
	b[0] = byte(ms >> 40)
	b[1] = byte(ms >> 32)
	b[2] = byte(ms >> 24)
	b[3] = byte(ms >> 16)
	b[4] = byte(ms >> 8)
	b[5] = byte(ms)

	b[6] = (b[6] & 0x0f) | 0x70 // version 7
	b[8] = (b[8] & 0x3f) | 0x80 // RFC 4122 variant

	uuidResponse(w, 7, formatUUID(b))
}

// resolveNameParams validates and resolves the namespace and name query params
// used by v3 and v5 handlers.
func resolveNameParams(w http.ResponseWriter, r *http.Request) (ns [16]byte, name string, ok bool) {
	name = r.URL.Query().Get("name")
	if name == "" {
		JSON(w, http.StatusBadRequest, map[string]string{"error": "name query parameter is required"})
		return
	}
	nsKey := r.URL.Query().Get("namespace")
	if nsKey == "" {
		nsKey = "dns"
	}
	ns, exists := uuidNamespaces[nsKey]
	if !exists {
		JSON(w, http.StatusBadRequest, map[string]string{
			"error":      "unknown namespace, must be one of: dns, url, oid, x500",
			"namespace":  nsKey,
		})
		return ns, name, false
	}
	return ns, name, true
}
