package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pivotal-golang/lager"

	"github.com/cloudfoundry-incubator/uaa-go-client"
	uaa_config "github.com/cloudfoundry-incubator/uaa-go-client/config"
	"github.com/pivotal-golang/clock"
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
	uaaAdminClientUsername := flag.String("uaa_admin_client_username", "", "UAA Admin Client username")
	uaaAdminClientSecret := flag.String("uaa_admin_client_secret", "", "UAA Admin Client secret")

	flag.Parse()

	logger := lager.NewLogger("Datadog Dashboard Generator")
	outSink := lager.NewWriterSink(os.Stderr, lager.DEBUG)
	logger.RegisterSink(outSink)

	if uaaAdminClientUsername == nil || *uaaAdminClientUsername == "" {
		log.Fatal("uaa_admin_client_username must be provided and not empty")
	}

	// if uaaAdminClientSecret == nil || *uaaAdminClientSecret == "" {
	// 	log.Fatal("uaaAdminClientSecret must be provided and not empty")
	// }

	if uaaDomain == nil || *uaaDomain == "" {
		log.Fatal("uaaDomain must be provided and not empty")
	}

	adminUaaConfig := &uaa_config.Config{
		ClientName:       *uaaAdminClientUsername,
		ClientSecret:     *uaaAdminClientSecret,
		UaaEndpoint:      fmt.Sprintf("https://%s", *uaaDomain),
		SkipVerification: true,
	}

	logger.Info("adminUAACONFIG:", lager.Data{"config": adminUaaConfig})

	adminUaaClient, err := uaa_go_client.NewClient(logger, adminUaaConfig, clock.NewClock())
	if err != nil {
		logger.Fatal("Failed to generate a UAA client", err)
	}

	opsmanClient := opsman.New(*opsmanIP, *opsmanUser, *opsmanPassword, adminUaaClient, logger)
	logger.Info("opsmanClient:", lager.Data{"client": opsmanClient})

	// Check we are using a supported Ops Man
	// err = opsman.ValidateAPIVersion(opsmanClient)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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
