// Package secretmanagersampler contains an HTTP Cloud Function.
package secretmanagersampler

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"jamesattensure.io/secretmanagersampler/config"
)

// SecretManagerSampler prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func SecretManagerSampler(w http.ResponseWriter, r *http.Request) {

	var d struct {
		Message string `json:"message"`
	}

	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("%+v\n", config.Data)

	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Hello World")

		return
	}
	if d.Message == "" {
		fmt.Fprint(w, "Hello World!")
		return
	}
	fmt.Fprint(w, html.EscapeString(d.Message))

}
