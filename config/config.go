package config

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v3"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
)

// Config is a yaml struct of config values
type Config struct {
	Data struct {
		NestOneOne string `yaml:"nest_one_one"`
		NestOneTwo string `yaml:"nest_one_two"`
		NestTwo    struct {
			NestTwoOne  string   `yaml:"nest_two_one"`
			NestTwoList []string `yaml:"nest_two_list"`
		}
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
func getSecretFromGSM() ([]byte, error) {

	projectId := getEnv("PROJECT_ID", "")
	secretName := getEnv("SECRET_NAME", "")
	secretVersion := getEnv("SECRET_VERSION", "latest")

	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectId, secretName, secretVersion),
	}
	resp, err := c.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Payload.Data, nil
}

// GetConfig returns a struct of type Config
func GetConfig() (*Config, error) {
	c := &Config{}
	data, err := getSecretFromGSM()
	var err2 error
	if err != nil {
		log.Printf("Could not retrieve from Google Secret Manager %v", err)
		data, err2 = ioutil.ReadFile("./data.yaml")
		if err2 != nil {
			return nil, fmt.Errorf("first error: %v, second error: %v", err, err2)
		}
	}

	err3 := yaml.Unmarshal([]byte(data), &c)
	if err3 != nil {
		return nil, err3
	}
	return c, nil
}
