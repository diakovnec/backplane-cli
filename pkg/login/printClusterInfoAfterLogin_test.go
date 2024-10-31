package login

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"

	// "github.com/openshift/backplane-cli/pkg/login" // Removed to avoid import cycle
	"github.com/openshift/backplane-cli/pkg/ocm/mocks" // Adjust the import path as necessary
)

// Initialize Ginkgo and Gomega
func TestLogin(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Login Suite")
}

var _ = ginkgo.Describe("Cluster Information", func() {
	var (
		ctrl             *gomock.Controller
		mockOCMInterface *mocks.MockOCMInterface
		clusterID        string
	)

	ginkgo.BeforeEach(func() {
		ctrl = gomock.NewController(ginkgo.GinkgoT())
		mockOCMInterface = mocks.NewMockOCMInterface(ctrl)
		clusterID = "test-cluster-id"
	})

	ginkgo.AfterEach(func() {
		ctrl.Finish()
	})

	ginkgo.Context("PrintClusterInfo", func() {
		ginkgo.It("should return an error when cluster info cannot be retrieved", func() {
			mockOCMInterface.EXPECT().GetClusterInfoByID(clusterID).Return(nil, fmt.Errorf("not found"))

			err := PrintClusterInfo(clusterID)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(err.Error()).To(gomega.Equal("error retrieving cluster info: not found"))
		})

		ginkgo.It("should display cluster information correctly", func() {
			// Create a mock clusterInfo with necessary methods
			clusterInfo := mocks.NewMockClusterInfo(ctrl)

			clusterInfo.EXPECT().ID().Return(clusterID)
			clusterInfo.EXPECT().Name().Return("Test Cluster")
			mockClusterStatus := mocks.NewMockClusterStatus(ctrl)
			clusterInfo.EXPECT().Status().Return(mockClusterStatus)
			clusterInfo.EXPECT().Region().Return(mocks.NewMockRegion(ctrl))
			clusterInfo.EXPECT().CloudProvider().Return(mocks.NewMockCloudProvider(ctrl))
			clusterInfo.EXPECT().Hypershift().Return(mocks.NewMockHypershift(ctrl))

			mockOCMInterface.EXPECT().GetClusterInfoByID(clusterID).Return(clusterInfo, nil)

			// Mock methods for the returned status, region, and cloud provider
			// Add similar expectations for Status, Region, CloudProvider, etc.

			err := PrintClusterInfo(clusterID)
			gomega.Expect(err).ToNot(gomega.HaveOccurred())
			// You can verify the output using logrus or a similar approach
		})
	})

	ginkgo.Context("PrintAccessProtectionStatus", func() {
		ginkgo.It("should display enabled access protection", func() {
			ocmConnection := mocks.NewMockOCMConnection(ctrl)
			mockOCMInterface.EXPECT().SetupOCMConnection().Return(ocmConnection, nil)
			defer ocmConnection.Close()

			mockOCMInterface.EXPECT().IsClusterAccessProtectionEnabled(ocmConnection, clusterID).Return(true, nil)

			PrintAccessProtectionStatus(clusterID)
			// Check the output for "Access Protection: Enabled"
		})

		ginkgo.It("should display disabled access protection", func() {
			ocmConnection := mocks.NewMockOCMConnection(ctrl)
			mockOCMInterface.EXPECT().SetupOCMConnection().Return(ocmConnection, nil)
			defer ocmConnection.Close()

			mockOCMInterface.EXPECT().IsClusterAccessProtectionEnabled(ocmConnection, clusterID).Return(false, nil)

			PrintAccessProtectionStatus(clusterID)
			// Check the output for "Access Protection: Disabled"
		})
	})

	ginkgo.Context("PrintOpenshiftVersion", func() {
		ginkgo.It("should print the OpenShift version", func() {
			clusterInfo := mocks.NewMockClusterInfo(ctrl)
			mockOCMInterface.EXPECT().GetClusterInfoByID(clusterID).Return(clusterInfo, nil)

			clusterInfo.EXPECT().GetOpenshiftVersion().Return("4.10", nil)

			PrintOpenshiftVersion(clusterID)
			// Check the output for "Openshift Version: 4.10"
		})
	})

	ginkgo.Context("GetLimitedSupportStatus", func() {
		ginkgo.It("should return fully supported status", func() {
			clusterInfo := mocks.NewMockClusterInfo(ctrl)
			mockOCMInterface.EXPECT().GetClusterInfoByID(clusterID).Return(clusterInfo, nil)

			mockClusterStatus := mocks.NewMockClusterStatus(ctrl)
			clusterInfo.EXPECT().Status().Return(mockClusterStatus)
			mockClusterStatus.EXPECT().LimitedSupportReasonCount().Return(0)

			status := GetLimitedSupportStatus(clusterID)
			gomega.Expect(status).To(gomega.Equal("0"))
			// Check the output for "Limited Support Status: Fully Supported"
		})

		ginkgo.It("should return limited support status", func() {
			clusterInfo := mocks.NewMockClusterInfo(ctrl)
			mockOCMInterface.EXPECT().GetClusterInfoByID(clusterID).Return(clusterInfo, nil)

			mockClusterStatus := mocks.NewMockClusterStatus(ctrl)
			clusterInfo.EXPECT().Status().Return(mockClusterStatus)
			mockClusterStatus.EXPECT().LimitedSupportReasonCount().Return(1)

			status := GetLimitedSupportStatus(clusterID)
			gomega.Expect(status).To(gomega.Equal("1"))
			// Check the output for "Limited Support Status: Limited Support"
		})
	})
})
