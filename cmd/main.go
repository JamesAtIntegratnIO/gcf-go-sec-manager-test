package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"jamesattensure.io/secretmanagersampler"
)

func main() {

	funcframework.RegisterHTTPFunction("/", secretmanagersampler.SecretManagerSampler)
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start %v", err)
	}
}
