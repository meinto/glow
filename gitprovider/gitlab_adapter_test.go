package gitprovider_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"

	"github.com/meinto/glow"
	mockb "github.com/meinto/glow/__mock__"
	mockg "github.com/meinto/glow/git/__mock__"
	"github.com/meinto/glow/gitprovider"
	. "github.com/meinto/glow/gitprovider"
	. "github.com/onsi/gomega"
)

var _ = Describe("Branch", func() {
	var gp Service
	var mockBranch mockb.MockBranchInterface
	var mockGitService mockg.MockNativeServiceInterface
	var mockHTTPClient HttpClient

	BeforeEach(func() {
		mockCtrl := gomock.NewController(GinkgoT())
		newMockGitService := mockg.NewMockNativeService(mockCtrl)
		mockBranch = mockb.NewMockBranch(mockCtrl)
		mockHTTPClient = &MockHTTPClient{}
		gp = NewGitlabService(gitprovider.Options{
			"https://mock.endpoint",
			"mock-namespace",
			"mock-project",
			"mock-token",
			false,
			mockHTTPClient,
		})
		gp.SetGitService(newMockGitService)
		mockGitService = gp.GitService().(mockg.MockNativeServiceInterface)
		log.Println(mockBranch, mockGitService)
	})

	AfterEach(func() {
		os.Setenv("CI_COMMIT_REF_NAME", "")
		os.Setenv("CI_GIT_USER", "")
		os.Setenv("CI_GIT_TOKEN", "")
	})

	Context("close command", func() {
		TestCreateMergeRequest := func(targetBranches []glow.Branch) {
			for range targetBranches {
				mockBranch.EXPECT().ShortBranchName().Return("source-branch")
			}
			mockHTTPClient.(*MockHTTPClient).SetRequestIntercaptionCallback(func(req *http.Request) {
				expectedURL, _ := url.Parse("https://mock.endpoint/api/v4/projects/mock-namespace%2Fmock-project/merge_requests")
				Expect(req.URL).To(Equal(expectedURL))
				Expect(req.Header).To(Equal(http.Header{
					"Content-Type":  []string{"application/json"},
					"Private-Token": []string{"mock-token"},
				}))
				bb, _ := ioutil.ReadAll(req.Body)
				targetBranchName := targetBranches[mockHTTPClient.(*MockHTTPClient).RequestCounter].ShortBranchName()
				expectedBodyString := `{"source_branch":"source-branch","target_branch":"` + targetBranchName + `","title":"Merge source-branch in ` + targetBranchName + `","remove_source_branch":false,"squash":false}`
				Expect(string(bb)).To(Equal(expectedBodyString))
			})
		}

		It("detects the close branches and creates merge request for them", func() {
			mockBranch.EXPECT().ShortBranchName().Return("mock-branch")
			mockGitService.EXPECT().RemoteBranchExists("mock-branch").Return("", "", nil)
			mockBranch.EXPECT().CanBeClosed().Return(true)
			branchList := []glow.Branch{
				glow.NewBranch("target-1"),
				glow.NewBranch("target-2"),
				glow.NewBranch("branch-3"),
			}
			mockGitService.EXPECT().BranchList().Return(branchList, "", "", nil)
			targetBranches := []glow.Branch{
				glow.NewBranch("target-1"),
				glow.NewBranch("target-2"),
			}
			mockBranch.EXPECT().CloseBranches(branchList).Return(targetBranches)
			TestCreateMergeRequest(targetBranches)
			gp.Close(mockBranch)
			Expect(mockHTTPClient.(*MockHTTPClient).RequestCounter).To(Equal(len(targetBranches)))
		})
	})

	Context("publish command", func() {
		TestCreateMergeRequest := func(targetBranch glow.Branch) {
			mockBranch.EXPECT().ShortBranchName().Return("source-branch")
			mockHTTPClient.(*MockHTTPClient).SetRequestIntercaptionCallback(func(req *http.Request) {
				expectedURL, _ := url.Parse("https://mock.endpoint/api/v4/projects/mock-namespace%2Fmock-project/merge_requests")
				Expect(req.URL).To(Equal(expectedURL))
				Expect(req.Header).To(Equal(http.Header{
					"Content-Type":  []string{"application/json"},
					"Private-Token": []string{"mock-token"},
				}))
				bb, _ := ioutil.ReadAll(req.Body)
				targetBranchName := targetBranch.ShortBranchName()
				expectedBodyString := `{"source_branch":"source-branch","target_branch":"` + targetBranchName + `","title":"Merge source-branch in ` + targetBranchName + `","remove_source_branch":false,"squash":false}`
				Expect(string(bb)).To(Equal(expectedBodyString))
			})
		}

		It("detects the publish branch and creates merge request", func() {
			mockBranch.EXPECT().ShortBranchName().Return("mock-branch")
			mockGitService.EXPECT().RemoteBranchExists("mock-branch").Return("", "", nil)
			mockBranch.EXPECT().CanBePublished().Return(true)
			publishBranch := glow.NewBranch("mock-branch")
			mockBranch.EXPECT().PublishBranch().Return(publishBranch)
			TestCreateMergeRequest(publishBranch)
			gp.Publish(mockBranch)
			Expect(mockHTTPClient.(*MockHTTPClient).RequestCounter).To(Equal(1))
		})
	})

	It("gets the cicd branch", func() {
		os.Setenv("CI_COMMIT_REF_NAME", "feature-branch")
		branch, _ := gp.GetCIBranch()
		Expect(branch.ShortBranchName()).To(Equal("feature-branch"))
	})

	It("detects the cicd origin", func() {
		os.Setenv("CI_GIT_USER", "username")
		os.Setenv("CI_GIT_TOKEN", "password")
		origin, _ := gp.DetectCICDOrigin()
		Expect(origin).To(Equal("https://username:password@mock.endpoint/mock-namespace/mock-project.git"))
	})
})
