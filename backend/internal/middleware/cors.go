package middleware

import (
	"github.com/go-chi/cors"
)

func (mw *MiddlewareManager) Cors() cors.Options {
	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	return cors.Options{
		// AllowedOrigins: []string{"*"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:3000"},
		// AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		OptionsPassthrough: false,
		MaxAge:             300, // Maximum value not ignored by any of major browsers
		Debug:              false,
	}
}
