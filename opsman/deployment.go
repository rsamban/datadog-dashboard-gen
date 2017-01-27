package opsman

import "fmt"

// NewDeployment creates a new installation deployment
func NewDeployment(installation *InstallationSettings, cfRelease string) *Deployment {

	var uaaJobParts []string
	for _, p := range getPartitions(installation, cfRelease, "uaa") {
		if p.InstanceCount > 0 {
			uaaJobParts = append(uaaJobParts, p.InstallationName)
		}
	}
	var routerJobParts []string
	var routerIps []string
	for _, p := range getPartitions(installation, cfRelease, "router") {
		if p.InstanceCount > 0 {
			routerJobParts = append(routerJobParts, p.InstallationName)
			for _, ip := range p.Ips {
				routerIps = append(routerIps, ip)
			}
		}
	}
	var ccJobParts []string
	for _, p := range getPartitions(installation, cfRelease, "cloud_controller") {
		if p.InstanceCount > 0 {
			ccJobParts = append(ccJobParts, p.InstallationName)
		}
	}
	var diegoBrainParts []string
	var diegoBrainIps []string
	for _, p := range getPartitions(installation, cfRelease, "diego_brain") {
		if p.InstanceCount > 0 {
			diegoBrainParts = append(diegoBrainParts, p.InstallationName)
			for _, ip := range p.Ips {
				diegoBrainIps = append(diegoBrainIps, ip)
			}
		}
	}
	var diegoCellParts []string
	var diegoCellIps []string
	for _, p := range getPartitions(installation, cfRelease, "diego_cell") {
		if p.InstanceCount > 0 {
			diegoCellParts = append(diegoCellParts, p.InstallationName)
			for _, ip := range p.Ips {
				diegoCellIps = append(diegoCellIps, ip)
			}
		}
	}
	var diegoDatabaseParts []string
	var diegoDatabaseIps []string
	for _, p := range getPartitions(installation, cfRelease, "diego_database") {
		if p.InstanceCount > 0 {
			diegoDatabaseParts = append(diegoDatabaseParts, p.InstallationName)
			for _, ip := range p.Ips {
				diegoDatabaseIps = append(diegoDatabaseIps, ip)
			}
		}
	}
	var uaaDatabaseParts []string
	for _, p := range getPartitions(installation, cfRelease, "uaadb") {
		if p.InstanceCount > 0 {
			uaaDatabaseParts = append(uaaDatabaseParts, p.InstallationName)
		}
	}
	var ccJobDatabaseParts []string
	for _, p := range getPartitions(installation, cfRelease, "ccdb") {
		if p.InstanceCount > 0 {
			ccJobDatabaseParts = append(ccJobDatabaseParts, p.InstallationName)
		}
	}

	deployment := &Deployment{
		Release:                     cfRelease,
		UaaJobs:                     uaaJobParts,
		RouterJobs:                  routerJobParts,
		RouterIps:                   routerIps,
		CloudControllerJobs:         ccJobParts,
		CloudControllerDatabaseJobs: ccJobDatabaseParts,
		DiegoBrainJobs:              diegoBrainParts,
		DiegoBrainIps:               diegoBrainIps,
		DiegoCellJobs:               diegoCellParts,
		DiegoCellIps:                diegoCellIps,
		DiegoDatabaseJobs:           diegoDatabaseParts,
		DiegoDatabaseIps:            diegoDatabaseIps,
		UaaDatabaseJobs:             uaaDatabaseParts,
	}

	return deployment
}

// ValidateAPIVersion checks for a supported API version
func ValidateAPIVersion(client *Client) error {
	version, err := client.GetAPIVersion()
	if err != nil {
		return err
	}

	if version != "2.0" {
		return fmt.Errorf("This version of Ops Manager: '" + version + "' is not supported")
	}

	return nil
}

func getPartitions(installation *InstallationSettings, productName, jobInstallationName string) []Partition {
	for _, product := range installation.Products {
		if product.Name == productName {
			for _, job := range product.Jobs {
				if job.InstallationName == jobInstallationName {
					return job.Partition
				}
			}
		}
	}
	return nil
}

func getIps(installation *InstallationSettings, productName, jobInstallationName string) []Partition {
	for _, product := range installation.Products {
		if product.Name == productName {

			for _, job := range product.Jobs {
				if job.InstallationName == jobInstallationName {
					return job.Partition
				}
			}
		}
	}
	return nil
}
