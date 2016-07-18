package opsman

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cloudfoundry-incubator/uaa-go-client"
	"github.com/pivotal-golang/lager"
	http "github.com/pivotalservices/datadog-dashboard-gen/http"
)

// Client configures an opsman connection
type Client struct {
	opsmanIP       string
	opsmanUsername string
	opsmanPassword string
	uaaClient      uaa_go_client.Client
	logger         lager.Logger
}

// New creates a new opsman Client
func New(opsmanIP, opsmanUsername, opsmanPassword string, uaaClient uaa_go_client.Client, logger lager.Logger) *Client {
	return &Client{
		opsmanIP:       opsmanIP,
		opsmanUsername: opsmanUsername,
		opsmanPassword: opsmanPassword,
		uaaClient:      uaaClient,
		logger:         logger,
	}
}

// GetAPIVersion returns the Ops Man API version
func (c *Client) GetAPIVersion() (string, error) {
	forceUpdate := true
	token, err := c.uaaClient.FetchToken(forceUpdate)
	if err != nil {
		return "", err
	}
	url := "https://" + c.opsmanIP + "/api/api_version"
	resp, err := http.SendRequest("GET", url, c.opsmanUsername, c.opsmanPassword, token.AccessToken, "")
	if err != nil {
		return "", err
	}
	res := bytes.NewBufferString(resp)
	decoder := json.NewDecoder(res)
	var ver Version
	err = decoder.Decode(&ver)
	if err != nil {
		return "", err
	}

	return ver.Version, nil
}

// GetCFDeployment returns the Elastic-Runtime deployment created by your Ops Manager
func (c *Client) GetCFDeployment(installation *InstallationSettings, products []Products) (*Deployment, error) {
	cfRelease := getProductGUID(products, "cf")
	if cfRelease == "" {
		return nil, fmt.Errorf("cf release not found")
	}

	return NewDeployment(installation, cfRelease), nil
}

// GetInstallationSettings retrieves installation settings for cf deployment
func (c *Client) GetInstallationSettings() (*InstallationSettings, error) {
	forceUpdate := true
	token, err := c.uaaClient.FetchToken(forceUpdate)
	if err != nil {
		return nil, err
	}
	url := "https://" + c.opsmanIP + "/api/installation_settings/"
	c.logger.Info("GetInstallationSettings", lager.Data{"url": url, "access_token": token.AccessToken, "opsmanUsername": c.opsmanUsername, "opsmanPassword": c.opsmanPassword})
	resp, err := http.SendRequest("GET", url, c.opsmanUsername, c.opsmanPassword, token.AccessToken, "")
	if err != nil {
		return nil, err
	}
	res := bytes.NewBufferString(resp)
	decoder := json.NewDecoder(res)
	var installation *InstallationSettings
	err = decoder.Decode(&installation)
	if err != nil {
		c.logger.Debug("Error marshalling POST body to json", lager.Data{
			"errMessage": err.Error(),
		})
		return nil, err
	}

	return installation, err
}

// GetProducts returns all the products in an OpsMan installation
func (c *Client) GetProducts() ([]Products, error) {
	forceUpdate := true
	token, err := c.uaaClient.FetchToken(forceUpdate)
	if err != nil {
		return nil, err
	}
	url := "https://" + c.opsmanIP + "/api/installation_settings/products"
	c.logger.Info("GetProducts", lager.Data{"url": url, "access_token": token.AccessToken, "opsmanUsername": c.opsmanUsername, "opsmanPassword": c.opsmanPassword})
	resp, err := http.SendRequest("GET", url, c.opsmanUsername, c.opsmanPassword, token.AccessToken, "")
	if err != nil {
		return nil, err
	}
	res := bytes.NewBufferString(resp)
	decoder := json.NewDecoder(res)
	var products []Products
	err = decoder.Decode(&products)
	if err != nil {
		c.logger.Debug("Error marshalling POST body to json", lager.Data{
			"errMessage": err.Error(),
		})
		return nil, err
	}
	return products, err
}

// gets the product GUID for a given product type
func getProductGUID(products []Products, productType string) string {
	for prod := range products {
		if products[prod].Type == productType {
			return products[prod].GUID
		}
	}
	return ""
}
