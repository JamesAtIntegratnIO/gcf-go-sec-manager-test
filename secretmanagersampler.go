// Package secretmanagersampler contains an HTTP Cloud Function.
package secretmanagersampler

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"

	"jamesattensure.io/secretmanagersampler/configbuilder"
)

type configuration struct {
	Data struct {
		NestOneOne string `yaml:"nest_one_one"`
		NestOneTwo string `yaml:"nest_one_two"`
		NestTwo    struct {
			NestTwoOne  string   `yaml:"nest_two_one"`
			NestTwoList []string `yaml:"nest_two_list"`
		}
	}
}

// SecretManagerSampler prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func SecretManagerSampler(w http.ResponseWriter, r *http.Request) {
	conf := configuration{}
	var d struct {
		Message string `json:"message"`
	}

	err := configbuilder.GetConfig(true, true, false, "", conf)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("%+v\n", conf.Data)

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
