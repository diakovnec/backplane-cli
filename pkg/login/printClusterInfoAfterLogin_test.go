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

// Mocking the necessary methods that will be called within PrintClusterInfo
type MockClusterInfo struct {
	id                string
	name              string
	status            string
	region            string
	cloudProvider     string
	hypershiftEnabled bool
}

func (m *MockClusterInfo) ID() string {
	return m.id
}

func (m *MockClusterInfo) Name() string {
	return m.name
}

func (m *MockClusterInfo) Status() *MockClusterStatus {
	return &MockClusterStatus{state: m.status}
}

func (m *MockClusterInfo) Region() *MockClusterRegion {
	return &MockClusterRegion{id: m.region}
}

func (m *MockClusterInfo) CloudProvider() *MockCloudProvider {
	return &MockCloudProvider{id: m.cloudProvider}
}

func (m *MockClusterInfo) Hypershift() *MockHypershift {
	return &MockHypershift{enabled: m.hypershiftEnabled}
}

func (m *MockClusterInfo) GetOpenshiftVersion() (string, error) {
	return "4.8.0", nil
}

// Mocked nested structs
type MockClusterStatus struct {
	state string
}

func (s *MockClusterStatus) State() string {
	return s.state
}

type MockClusterRegion struct {
	id string
}

func (r *MockClusterRegion) ID() string {
	return r.id
}

type MockCloudProvider struct {
	id string
}

func (p *MockCloudProvider) ID() string {
	return p.id
}

type MockHypershift struct {
	enabled bool
}

func (h *MockHypershift) Enabled() bool {
	return h.enabled
}

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
			Status(cmv1.NewClusterStatus().State("Running")).
			Region(cmv1.NewCloudRegion().ID("us-east-1")).
			OpenshiftVersion("4.14.8").Build()

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
		Expect(output).To(ContainSubstring("Cluster Status:           Running\n"))
		Expect(output).To(ContainSubstring("Cluster Region:           us-east-1\n"))
		Expect(output).To(ContainSubstring("Cluster Provider:         aws\n"))
		Expect(output).To(ContainSubstring("Hypershift Enabled:       false\n"))
		Expect(output).To(ContainSubstring("Version:                  4.14.8\n"))
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
