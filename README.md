## Datadog Dashboard Generator

This is a command line utility that can be used to deploy a MVP "StopLights" Dashboard to a given Users Datadog Subscription from a template.

## Logical Flow

1. Tool queries PCF Ops Manager for the following:

        cf-release string (e.g. cf-76553423523ab)
        job indexes/partitions from CF manifest

2. Reads template of PCF Stoplights Dashboard

3. Uploads a generated Datadog dashboard combining Ops Man vars & template

## Example usage

```
datadog-dashboard-gen \
  -opsman_user=<REPLACE-WITH-OPS-MANAGER-USERNAME> \
  -opsman_passwd=<REPLACE-WITH-OPS-MANAGER-PASSWORD> \
  -opsman_ip=<REPLACE-WITH-OPS-MANAGER-IP> \
  -uaa_domain=<REPLACE-WITH-UAA-DOMAIN> \
  -ddapikey=<REPLACE-WITH-DATADOG-API-KEY> \
  -ddappkey=<REPLACE-WITH-DATADOG-APP-KEY>
```

## Build & Run

1. Clone repo

        git clone https://github.com/pivotalservices/datadog-dashboard-gen.git

1. Build binary

        cd datadog-dashboard-gen
        go install

1. Run program to upload the Stoplights dashboard

        $GOPATH/bin/datadog-dashboard-gen -opsman_user=<OPSMAN_USER> -opsman_password=<OPSMAN_PASSWORD> -opsman_ip=<OPSMAN_IP> \
        -ddapikey=<DATADOG_API_KEY> -ddappkey=<DATADOG_APP_KEY>

## Generate code from template

1. Install [ego](https://github.com/benbjohnson/ego)

1. Run `ego` (template is located under `templates/screen`)
```bash
ego -package datadog -o datadog/stoplights.go
```
