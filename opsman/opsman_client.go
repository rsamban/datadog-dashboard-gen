package opsman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pivotal-golang/lager"
	http "github.com/pivotalservices/datadog-dashboard-gen/http"
)

// Client configures an opsman connection
type Client struct {
	opsmanVersion  string
	opsmanIP       string
	opsmanUsername string
	opsmanPassword string
	uaaDomain      string
	logger         lager.Logger
}

// New creates a new opsman Client
func New(opsmanVersion, opsmanIP, opsmanUsername, opsmanPassword, uaaDomain string, logger lager.Logger) *Client {
	return &Client{
		opsmanVersion:  opsmanVersion,
		opsmanIP:       opsmanIP,
		opsmanUsername: opsmanUsername,
		opsmanPassword: opsmanPassword,
		uaaDomain:      uaaDomain,
		logger:         logger,
	}
}

// GetAPIVersion returns the Ops Man API version
func (c *Client) GetAPIVersion() (string, error) {

	//TODO what exactly this API version meaning here. Ops Man 1.7 showing 2.0 version here
	// if c.opsmanVersion != "1.6" {
	// 	return "", fmt.Errorf("https://%s/api/api_version is no longer supported", c.opsmanIP)
	// }

	token, err := c.fetchToken()
	if err != nil {
		return "", err
	}
	url := "https://" + c.opsmanIP + "/api/api_version"
	resp, err := http.SendRequest("GET", url, c.opsmanUsername, c.opsmanPassword, token, "")
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

	c.GetIps(installation)
	fmt.Printf("RAMESH:::installation-Begin\n")
	fmt.Printf("%v\n", installation)
	fmt.Printf("RAMESH:::installation-End\n")
	deployment := NewDeployment(installation, cfRelease)
	return deployment, nil

}

// GetInstallationSettings retrieves installation settings for cf deployment
func (c *Client) GetInstallationSettings() (*InstallationSettings, error) {
	token, err := c.fetchToken()
	if err != nil {
		return nil, err
	}
	url := "https://" + c.opsmanIP + "/api/installation_settings/"
	c.logger.Info("GetInstallationSettings", lager.Data{"url": url, "access_token": token, "opsmanUsername": c.opsmanUsername, "opsmanPassword": c.opsmanPassword})
	resp, err := http.SendRequest("GET", url, c.opsmanUsername, c.opsmanPassword, token, "")
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

// GetIps retrieves installation settings for cf deployment
func (c *Client) GetIps(installation *InstallationSettings) error {
	token, err := c.fetchToken()
	if err != nil {
		return err
	}

	for i, product := range installation.Products {
		url := "https://" + c.opsmanIP + "/api/v0/deployed/products/" + product.GUID + "/static_ips"
		c.logger.Info("GetIps", lager.Data{"url": url, "access_token": token, "opsmanUsername": c.opsmanUsername, "opsmanPassword": c.opsmanPassword})
		resp, err := http.SendRequest("GET", url, c.opsmanUsername, c.opsmanPassword, token, "")
		if err != nil {
			return err
		}
		res := bytes.NewBufferString(resp)
		decoder := json.NewDecoder(res)
		var ips []Ip
		err = decoder.Decode(&ips)
		if err != nil {
			c.logger.Debug("Error marshalling POST body to json", lager.Data{
				"errMessage": err.Error(),
			})
			return err
		}
		for j, job := range product.Jobs {
			for k, partition := range job.Partition {
				s := strings.Split(partition.InstallationName, "-")
				s[0] = job.GUID
				ipName := strings.Join(s, "-")
				for _, ip := range ips {
					if ipName == ip.Name {
						fmt.Printf("RAMESH::: Assigning:::%s, %v \n", ipName, ip.Ips)
						installation.Products[i].Jobs[j].Partition[k].Ips = ip.Ips
						//temp := installation.Products[i].Jobs[j].Partition[k].Ips
						//fmt.Printf("RAMESH:::temp %v \n", temp)
						//partition.Ips = ip.Ips
					}
				}
			}
		}
	}
	return err
}

// GetProducts returns all the products in an OpsMan installation
func (c *Client) GetProducts() ([]Products, error) {
	token, err := c.fetchToken()
	if err != nil {
		return nil, err
	}
	url := "https://" + c.opsmanIP + "/api/installation_settings/products"
	c.logger.Info("GetProducts", lager.Data{"url": url, "access_token": token, "opsmanUsername": c.opsmanUsername, "opsmanPassword": c.opsmanPassword})
	resp, err := http.SendRequest("GET", url, c.opsmanUsername, c.opsmanPassword, token, "")
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

func (c *Client) fetchToken() (string, error) {
	var clientID = "opsman"
	var clientSecret = ""
	c.logger.Debug("oauthHTTPGet", lager.Data{"uaaDomain": "https://" + c.uaaDomain})

	return http.GetToken("https://"+c.uaaDomain, c.opsmanUsername, c.opsmanPassword, clientID, clientSecret)
}
