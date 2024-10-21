package login

import (
	"fmt"

	logger "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"

	ocmsdk "github.com/openshift-online/ocm-sdk-go"
	"github.com/openshift/backplane-cli/pkg/cli/config"
	"github.com/openshift/backplane-cli/pkg/ocm"
)

// GetRestConfig returns a client-go *rest.Config which can be used to programmatically interact with the
// Kubernetes API of a provided clusterID
func ConnectToTheCluster(bp config.BackplaneConfiguration, ocmConnection *ocmsdk.Connection, clusterID string) (*rest.Config, error) {
	cluster, err := ocm.DefaultOCMInterface.GetClusterInfoByIDWithConn(ocmConnection, clusterID)
	if err != nil {
		return nil, err
	}

	accessToken, err := ocm.DefaultOCMInterface.GetOCMAccessTokenWithConn(ocmConnection)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

//////////////////////////////////////////////////////////////////
//////////// this is working code ////////////////////////////////

// displayClusterInfo retrieves and displays basic information about the target cluster.
func displayClusterInfo(clusterID string) error {
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

// to call the current function
func calldisplayClusterInfo(clusterID string) error {
	if err := displayClusterInfo(clusterID); err != nil {
		return err
	}
	return nil
}

//////////////////////////////////////////////////////////////////
