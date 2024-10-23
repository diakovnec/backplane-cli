package login

import (
	"context"
	"fmt"

	//	"github.com/openshift-online/ocm-sdk-go"
	accountsmgmtv1 "github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1" // Import the package
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
	fmt.Printf("Cluster ID: %s\n", clusterInfo.ID())
	fmt.Printf("Cluster Name: %s\n", clusterInfo.Name())
	fmt.Printf("Cluster Status: %s\n", clusterInfo.Status().State())
	fmt.Printf("Cluster Region: %s\n", clusterInfo.Region().ID())
	fmt.Printf("Cluster Provider: %s\n", clusterInfo.CloudProvider().ID())

	logger.Info("Basic cluster information displayed.")
	return nil

}

// PrintAccessProtectionStatus retrieves and displays the access protection status of the target cluster.

func PrintAccessProtectionStatus(clusterID string) {
	ocmConnection, _ := ocm.DefaultOCMInterface.SetupOCMConnection()

	defer ocmConnection.Close()
	enabled, _ := ocm.DefaultOCMInterface.IsClusterAccessProtectionEnabled(ocmConnection, clusterID)
	if enabled {
		fmt.Println("Access protection: Enabled")
	} else {
		fmt.Println("Access protection: Disabled")
	}
}

//Print Customer information

func PrintOrgName(id string) (*accountsmgmtv1.Organization, error) {
	// Set up the OCM connection
	ocmConnection, _ := ocm.DefaultOCMInterface.SetupOCMConnection()

	defer ocmConnection.Close()

	// Fetch the organization details
	orgClient := ocmConnection.AccountsMgmt().V1().Organizations()
	orgResponse, err := orgClient.Organization(id).Get().SendContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get organization: %v", err)
	}

	fmt.Printf("Organization ID: %s\n", orgResponse.Body().ID())
	fmt.Printf("Organization Name: %s\n", orgResponse.Body().Name())
	return orgResponse.Body(), nil
}
