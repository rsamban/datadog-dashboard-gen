package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pivotal-golang/lager"

	"github.com/pivotalservices/datadog-dashboard-gen/datadog"
	"github.com/pivotalservices/datadog-dashboard-gen/opsman"
)

func main() {
	// Declare Flags
	opsmanUser := flag.String("opsman_user", "", "Ops Manager User")
	opsmanPassword := flag.String("opsman_password", "", "Ops Manager Password")
	opsmanIP := flag.String("opsman_ip", "192.168.200.10", "Ops Manager IP")
	ddAPIKey := flag.String("ddapikey", "12345-your-api-key-6789", "Datadog API Key")
	ddAppKey := flag.String("ddappkey", "12345-your-app-key-6789", "Datadog Application Key")
	useOpsMetrics := flag.Bool("use_ops_metrics", false, "Generate template from an PCF Ops Metrics deployment")
	saveFile := flag.String("save_as", "", "Save generated dashboard on local disk")
	uaaDomain := flag.String("uaa_domain", "", "UAA Domain")
	version := flag.String("version", "1.7", "Ops Manager Version")

	flag.Parse()

	logger := lager.NewLogger("Datadog Dashboard Generator")
	outSink := lager.NewWriterSink(os.Stderr, lager.DEBUG)
	logger.RegisterSink(outSink)

	if opsmanUser == nil || *opsmanUser == "" {
		log.Fatal("opsman_user must be provided and not empty")
	}

	if opsmanPassword == nil || *opsmanPassword == "" {
		log.Fatal("opsman_password must be provided and not empty")
	}

	if opsmanIP == nil || *opsmanIP == "" {
		log.Fatal("opsman_ip must be provided and not empty")
	}

	if ddAPIKey == nil || *ddAPIKey == "" {
		log.Fatal("ddapikey must be provided and not empty")
	}

	if ddAppKey == nil || *ddAppKey == "" {
		log.Fatal("ddappkey must be provided and not empty")
	}

	if uaaDomain == nil || *uaaDomain == "" {
		log.Fatal("uaaDomain must be provided and not empty")
	}

	opsmanClient := opsman.New(*version, *opsmanIP, *opsmanUser, *opsmanPassword, *uaaDomain, logger)
	logger.Info("opsmanClient:", lager.Data{"client": opsmanClient})

	// Check we are using a supported Ops Man
	err := opsman.ValidateAPIVersion(opsmanClient)
	if err != nil {
		logger.Error("ValidateAPIVersion", err)
	}

	// Get installation settings from Ops Man foundation
	installation, err := opsmanClient.GetInstallationSettings()
	if err != nil {
		log.Fatal(err)
	}

	products, err := opsmanClient.GetProducts()
	if err != nil {
		log.Fatal(err)
	}

	deployment, err := opsmanClient.GetCFDeployment(installation, products)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	if *useOpsMetrics {
		err = datadog.StopLightsOpsMetricsTemplate(&buf, deployment)
	} else {
		err = datadog.StopLightsTemplate(&buf, deployment)
	}
	if err != nil {
		log.Fatal(err)
	}

	if *saveFile != "" {
		err := ioutil.WriteFile(*saveFile, buf.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
	dashboardJSON := buf.String()

	if _, err := datadog.CreateStoplightDashboard(*ddAPIKey, *ddAppKey, dashboardJSON); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Your PCF Stoplights Datadog dashboard has been published ... Go Fetch @ https://app.datadoghq.com/dash/list :)")
}
