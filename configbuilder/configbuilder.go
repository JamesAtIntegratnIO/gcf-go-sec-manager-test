package configbuilder

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1beta1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	"gopkg.in/yaml.v2"
)

type gcloudVars struct {
	projectID     string
	secretName    string
	secretVersion string
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func setGcloudVars() gcloudVars {
	return gcloudVars{
		getEnv("PROJECT_ID", ""),
		getEnv("SECRET_NAME", ""),
		getEnv("SECRET_VERSION", "latest"),
	}
}

func (g gcloudVars) getSecretFromGSM() ([]byte, error) {
	if g.projectID == "" || g.secretName == "" {
		return nil, errors.New(`
		environment variables for gcp are not set
		please set 'PROJECT_ID', 'SECRET_NAME', and 'SECRET_VERSION'
		`)
	}
	ctx := context.Background()
	c, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s",
			g.projectID,
			g.secretName,
			g.secretVersion),
	}
	resp, err := c.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Payload.Data, nil
}

// GetConfig returns a struct of type Config
func GetConfig(useGCPSecrets bool,
	yamlConfig bool, jsonConfig bool,
	filePath string, config interface{}) error {
	var data []byte
	var err error
	if yamlConfig && jsonConfig {
		return errors.New("yamlconfig and jsonconfig cannot both be true")
	}

	if useGCPSecrets {
		gcloudVars := setGcloudVars()
		data, err = gcloudVars.getSecretFromGSM()
		if err != nil {
			return err
		}
	} else {
		data, err = ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
	}
	if yamlConfig {
		err = yaml.Unmarshal([]byte(data), &config)
		if err != nil {
			return err
		}
	} else if jsonConfig {
		err = json.Unmarshal([]byte(data), &config)
		if err != nil {
			return err
		}
	}

	return nil
}
