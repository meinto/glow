package gitprovider_test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"

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
				expectedBodyString := `{"source_branch":"source-branch","target_branch":"` + targetBranchName + `","title":"Merge source-branch in ` + targetBranchName + `","remove_source_branch":true,"squash":false}`
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

})
