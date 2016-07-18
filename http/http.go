package http

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
)

// SendRequest sends http requests
func SendRequest(method, url, user, passwd, clientToken, data string) (string, error) {
	//Ignore Self Signed SSL
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	//make Request Object
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	if err != nil {
		return "", err
	}

	//Set Auth
	if clientToken != "" {
		req.Header.Add("Authorization", "Bearer "+clientToken)
	} else if user != "" && passwd != "" {
		req.SetBasicAuth(user, passwd)
	}

	//If POST set header
	if method == "POST" {
		req.Header.Add("Content-type", "application/json")
	}
	dumpHttp := true
	if dumpHttp {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	//Make Client http Request
	client := http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if dumpHttp {
		dump, _ := httputil.DumpResponse(res, true)
		fmt.Println(string(dump))
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	//If POST verify Dashboard was published
	if method == "POST" && res.Status != "200 OK" {
		return "", fmt.Errorf("got " + res.Status + " when sending dashboard to datadog; expecting 200")
	}

	return string(body), nil
}
