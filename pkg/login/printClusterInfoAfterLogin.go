package login

import (
	"fmt"

	"github.com/openshift/backplane-cli/pkg/ocm"
	logger "github.com/sirupsen/logrus"
)

//displayClusterInfo retrieves and displays basic information about the target cluster.

func PrintClusterInfo(clusterID string) error {
	logger := logger.WithField("clusterID", clusterID)

	// Retrieve cluster information
	clusterInfo, err := ocm.DefaultOCMInterface.GetClusterInfoByID(clusterID)
	if err != nil {
		return fmt.Errorf("error retrieving cluster info: %v", err)
	}

	// Display cluster information
	fmt.Printf("\n%-25s %s\n", "Cluster ID:", clusterInfo.ID())
	fmt.Printf("%-25s %s\n", "Cluster Name:", clusterInfo.Name())
	fmt.Printf("%-25s %s\n", "Cluster Status:", clusterInfo.Status().State())
	fmt.Printf("%-25s %s\n", "Cluster Region:", clusterInfo.Region().ID())
	fmt.Printf("%-25s %s\n", "Cluster Provider:", clusterInfo.CloudProvider().ID())
	fmt.Printf("%-25s %t\n", "Hypershift Enabled:", clusterInfo.Hypershift().Enabled())

	PrintOpenshiftVersion(clusterID)
	PrintAccessProtectionStatus(clusterID)
	// Display access protection status
	logger.Info("Basic cluster information displayed.")
	return nil

}

// PrintAccessProtectionStatus retrieves and displays the access protection status of the target cluster.

func PrintAccessProtectionStatus(clusterID string) {
	ocmConnection, err := ocm.DefaultOCMInterface.SetupOCMConnection()
	if err != nil {
		fmt.Printf("Error setting up OCM connection: %v\n", err)
		return
	}
	defer ocmConnection.Close()
	enabled, _ := ocm.DefaultOCMInterface.IsClusterAccessProtectionEnabled(ocmConnection, clusterID)
	if enabled {
		fmt.Printf("%-25s %s", "Access protection:", "Enabled")
	} else {
		fmt.Printf("%-25s %s", "Access protection:", "Disabled\n")
	}

}

func PrintOpenshiftVersion(clusterID string) {
	clusterInfo, err := ocm.DefaultOCMInterface.GetClusterInfoByID(clusterID)
	if err != nil {
		fmt.Println("Error retrieving cluster info: ", err)
		return
	}
	openshiftVersion, _ := clusterInfo.GetOpenshiftVersion()
	fmt.Printf("%-25s %s\n", "Openshift Version: ", openshiftVersion)
}
