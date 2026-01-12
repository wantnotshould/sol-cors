// Package cors
// Copyright 2026 wantnotshould. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.
package cors

import (
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/wantnotshould/sol"
)

// Config defines the configuration for CORS handling
type Config struct {
	AllowOrigins     []string
	AllowAllOrigins  bool
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

// CORS is a middleware handler for handling Cross-Origin Resource Sharing (CORS).
// It checks the request headers and sets the appropriate CORS headers in the response.
func CORS(cfg Config) sol.HandlerFunc {
	// Validate conflict between AllowAllOrigins and AllowCredentials
	if cfg.AllowAllOrigins && cfg.AllowCredentials {
		panic("cors: AllowAllOrigins=true conflicts with AllowCredentials=true")
	}

	// Set default methods and max-age if not provided
	if len(cfg.AllowMethods) == 0 {
		cfg.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodHead}
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 86400 // Default max-age is 1 day
	}

	// Precompute headers for CORS responses
	methods := strings.Join(cfg.AllowMethods, ", ")
	headers := strings.Join(cfg.AllowHeaders, ", ")
	maxAge := strconv.Itoa(cfg.MaxAge)

	return func(c *sol.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin == "" {
			// No Origin header, skip CORS handling
			c.Next()
			return
		}

		// Check if the origin is allowed
		if !cfg.AllowAllOrigins && !slices.Contains(cfg.AllowOrigins, origin) {
			// Origin not allowed, skip handling
			c.Next()
			return
		}

		// Set CORS headers
		setCORSHeaders(c, cfg, origin, methods, headers, maxAge)

		if c.Method() == http.MethodOptions &&
			c.Header("Access-Control-Request-Method") != "" {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed to the next handler
		c.Next()
	}
}

// setCORSHeaders sets the necessary CORS headers to the response.
func setCORSHeaders(c *sol.Context, cfg Config, origin, methods, headers, maxAge string) {
	// Set Access-Control-Allow-Origin header
	if cfg.AllowAllOrigins {
		c.SetHeader("Access-Control-Allow-Origin", "*")
	} else {
		c.SetHeader("Access-Control-Allow-Origin", origin)
		c.SetHeader("Vary", "Origin") // This helps caching servers differentiate responses by Origin
	}

	// Set AllowCredentials header if true
	if cfg.AllowCredentials {
		c.SetHeader("Access-Control-Allow-Credentials", "true")
	}

	// Set AllowMethods header
	if methods != "" {
		c.SetHeader("Access-Control-Allow-Methods", methods)
	}

	// Set AllowHeaders header
	if headers != "" {
		c.SetHeader("Access-Control-Allow-Headers", headers)
	}

	// Set MaxAge header
	c.SetHeader("Access-Control-Max-Age", maxAge)
}

// Default is a default CORS handler with common settings for most APIs
func Default() sol.HandlerFunc {
	return CORS(Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowHeaders:    []string{"Content-Type", "Authorization"},
		MaxAge:          86400, // 1 day
	})
}
