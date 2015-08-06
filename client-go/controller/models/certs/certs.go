package certs

import (
	"encoding/json"
	"fmt"

	"github.com/brendangibat/deis/client-go/controller/api"
	"github.com/brendangibat/deis/client-go/controller/client"
)

// List certs registered with the controller.
func List(c *client.Client) ([]api.Cert, error) {
	body, err := c.BasicRequest("GET", "/v1/certs/", nil)

	if err != nil {
		return []api.Cert{}, err
	}

	res := api.Certs{}
	if err = json.Unmarshal([]byte(body), &res); err != nil {
		return []api.Cert{}, err
	}

	return res.Certs, nil
}

// New creates a new cert.
func New(c *client.Client, cert string, key string, commonName string) (api.Cert, error) {
	req := api.CertCreateRequest{Certificate: cert, Key: key, Name: commonName}
	reqBody, err := json.Marshal(req)

	if err != nil {
		return api.Cert{}, err
	}

	resBody, err := c.BasicRequest("POST", "/v1/certs/", reqBody)

	if err != nil {
		return api.Cert{}, err
	}

	resCert := api.Cert{}
	if err = json.Unmarshal([]byte(resBody), &resCert); err != nil {
		return api.Cert{}, err
	}

	return resCert, nil
}

// Delete removes a cert.
func Delete(c *client.Client, commonName string) error {
	u := fmt.Sprintf("/v1/certs/%s", commonName)

	_, err := c.BasicRequest("DELETE", u, nil)
	return err
}
