// package login

// import (
// 	"bytes"
// 	"fmt"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// 	"github.com/openshift/backplane-cli/pkg/ocm"
// 	ocmMock "github.com/openshift/backplane-cli/pkg/ocm/mocks"
// 	"github.com/sirupsen/logrus"
// )

// // Mocking the necessary methods that will be called within PrintClusterInfo
// type MockClusterInfo struct {
// 	id                string
// 	name              string
// 	status            string
// 	region            string
// 	cloudProvider     string
// 	hypershiftEnabled bool
// }

// func (m *MockClusterInfo) ID() string {
// 	return m.id
// }

// func (m *MockClusterInfo) Name() string {
// 	return m.name
// }

// func (m *MockClusterInfo) Status() *MockClusterStatus {
// 	return &MockClusterStatus{state: m.status}
// }

// func (m *MockClusterInfo) Region() *MockClusterRegion {
// 	return &MockClusterRegion{id: m.region}
// }

// func (m *MockClusterInfo) CloudProvider() *MockCloudProvider {
// 	return &MockCloudProvider{id: m.cloudProvider}
// }

// func (m *MockClusterInfo) Hypershift() *MockHypershift {
// 	return &MockHypershift{enabled: m.hypershiftEnabled}
// }

// func (m *MockClusterInfo) GetOpenshiftVersion() (string, error) {
// 	return "4.8.0", nil
// }

// type MockClusterStatus struct {
// 	state string
// }

// func (s *MockClusterStatus) State() string {
// 	return s.state
// }

// func (s *MockClusterStatus) LimitedSupportReasonCount() int {
// 	return 0
// }

// type MockClusterRegion struct {
// 	id string
// }

// func (r *MockClusterRegion) ID() string {
// 	return r.id
// }

// type MockCloudProvider struct {
// 	id string
// }

// func (p *MockCloudProvider) ID() string {
// 	return p.id
// }

// type MockHypershift struct {
// 	enabled bool
// }

// func (h *MockHypershift) Enabled() bool {
// 	return h.enabled
// }

// var _ = Describe("PrintClusterInfo", func() {
// 	var (
// 		clusterID        string
// 		mockCluster      *MockClusterInfo
// 		buf              *bytes.Buffer
// 		mockOcmInterface *ocmMock.MockOCMInterface
// 		mockCtrl         *gomock.Controller
// 	)

// 	BeforeEach(func() {
// 		mockCtrl = gomock.NewController(GinkgoT())
// 		mockOcmInterface = ocmMock.NewMockOCMInterface(mockCtrl)
// 		ocm.DefaultOCMInterface = mockOcmInterface

// 		clusterID = "test-cluster-id"
// 		mockCluster = &MockClusterInfo{
// 			id:                clusterID,
// 			name:              "Test Cluster",
// 			status:            "Running",
// 			region:            "us-east-1",
// 			cloudProvider:     "AWS",
// 			hypershiftEnabled: true,
// 		}
// 		buf = new(bytes.Buffer)
// 		logrus.SetOutput(buf)

// 		// Mocking the ocm.DefaultOCMInterface
// 		mockOcmInterface.EXPECT().GetClusterInfoByID(clusterID).Return(mockCluster, nil).AnyTimes()
// 		mockOcmInterface.EXPECT().SetupOCMConnection().Return(nil, nil).AnyTimes()
// 		mockOcmInterface.EXPECT().IsClusterAccessProtectionEnabled(nil, clusterID).Return(true, nil).AnyTimes()
// 	})

// 	It("should print cluster information", func() {
// 		err := PrintClusterInfo(clusterID)
// 		Expect(err).To(BeNil())

// 		output := buf.String()
// 		Expect(output).To(ContainSubstring(fmt.Sprintf("Cluster ID:                 %s\n", clusterID)))
// 		Expect(output).To(ContainSubstring("Cluster Name:               Test Cluster\n"))
// 		Expect(output).To(ContainSubstring("Cluster Status:             Running\n"))
// 		Expect(output).To(ContainSubstring("Cluster Region:             us-east-1\n"))
// 		Expect(output).To(ContainSubstring("Cluster Provider:           AWS\n"))
// 		Expect(output).To(ContainSubstring("Hypershift Enabled:         true\n"))
// 		Expect(output).To(ContainSubstring("Openshift Version:          4.8.0\n"))
// 		Expect(output).To(ContainSubstring("Access Protection:          Enabled"))
// 		Expect(output).To(ContainSubstring("Limited Support Status:     Fully Supported\n"))
// 	})

// 	AfterEach(func() {
// 		// Reset the ocm.DefaultOCMInterface to avoid side effects in other tests
// 		ocm.DefaultOCMInterface = nil
// 		mockCtrl.Finish()
// 	})
// })

// func TestLogin(t *testing.T) {
// 	RegisterFailHandler(Fail)
// 	RunSpecs(t, "Login Suite")
// }
