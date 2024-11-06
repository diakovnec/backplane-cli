package login

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/backplane-cli/pkg/ocm"
	ocmMock "github.com/openshift/backplane-cli/pkg/ocm/mocks"

	cmv1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

var _ = Describe("PrintClusterInfo", func() {
	var (
		clusterID        string
		buf              *bytes.Buffer
		mockOcmInterface *ocmMock.MockOCMInterface
		mockCtrl         *gomock.Controller
		oldStdout        *os.File
		r, w             *os.File
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockOcmInterface = ocmMock.NewMockOCMInterface(mockCtrl)
		ocm.DefaultOCMInterface = mockOcmInterface

		clusterID = "test-cluster-id"
		buf = new(bytes.Buffer)
		log.SetOutput(buf)

		// Redirect standard output to the buffer
		oldStdout = os.Stdout
		r, w, _ = os.Pipe()
		os.Stdout = w

		clusterInfo, _ := cmv1.NewCluster().
			ID(clusterID).
			Name("Test Cluster").
			CloudProvider(cmv1.NewCloudProvider().ID("aws")).
			State(cmv1.ClusterState("ready")).
			//Status(cmv1.NewClusterStatus().State("Running")).
			Region(cmv1.NewCloudRegion().ID("us-east-1")).
			Hypershift(cmv1.NewHypershift().Enabled(false)).
			OpenshiftVersion("4.14.8").
			Status(cmv1.NewClusterStatus().LimitedSupportReasonCount(0)).
			Build()

		mockOcmInterface.EXPECT().GetClusterInfoByID(clusterID).Return(clusterInfo, nil).AnyTimes()
	})

	It("should print cluster information", func() {
		err := PrintClusterInfo(clusterID)
		Expect(err).To(BeNil())

		// Capture the output
		w.Close()
		os.Stdout = oldStdout
		_, _ = buf.ReadFrom(r)

		output := buf.String()
		Expect(output).To(ContainSubstring(fmt.Sprintf("Cluster ID:               %s\n", clusterID)))
		Expect(output).To(ContainSubstring("Cluster Name:             Test Cluster\n"))
		Expect(output).To(ContainSubstring("Cluster Status:           ready\n"))
		Expect(output).To(ContainSubstring("Cluster Region:           us-east-1\n"))
		Expect(output).To(ContainSubstring("Cluster Provider:         aws\n"))
		Expect(output).To(ContainSubstring("Hypershift Enabled:       false\n"))
		Expect(output).To(ContainSubstring("Version:                  4.14.8\n"))
		Expect(output).To(ContainSubstring("Limited Support Status:   Fully Supported\n"))
	})

	AfterEach(func() {
		// Reset the ocm.DefaultOCMInterface to avoid side effects in other tests
		ocm.DefaultOCMInterface = nil
	})
})

func TestLogin(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Login Suite")
}
