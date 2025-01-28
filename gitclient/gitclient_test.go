package gitclient_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/zszabo-rh/issues-operator/gitclient"
)

var _ = Describe("Gitclient", func() {

	var (
		originalToken string
		issueId       int
	)

	Context("when the environment is properly configured", func() {
		BeforeEach(func() {
			originalToken = os.Getenv("GITTOKEN")
			os.Setenv("GITTOKEN", "abc123")
		})

		AfterEach(func() {
			if originalToken == "" {
				os.Unsetenv("GITTOKEN")
			} else {
				os.Setenv("GITTOKEN", originalToken)
			}
		})

		It("should return the git token", func() {
			token, err := gitclient.GetToken()
			Expect(err).ToNot(HaveOccurred())
			Expect(token).To(Equal("abc123"))
		})
	})

	Context("when the environment is not properly configured", func() {
		BeforeEach(func() {
			originalToken = os.Getenv("GITTOKEN")
			os.Unsetenv("GITTOKEN")
		})

		AfterEach(func() {
			if originalToken != "" {
				os.Setenv("GITTOKEN", originalToken)
			}
		})

		It("should return an error", func() {
			token, err := gitclient.GetToken()
			Expect(err).To(HaveOccurred())
			Expect(token).To(Equal(""))
		})
	})

	Context("when proper input url is provided", func() {
		It("should return the converted url", func() {
			url, err := gitclient.BuildUrl("git@github.com:myrepo/myuser.git")
			Expect(err).ToNot(HaveOccurred())
			Expect(url).To(Equal("https://api.github.com/repos/myrepo/myuser/issues"))
		})
	})

	Context("when wrong input url is provided", func() {
		It("should return an error", func() {
			url, err := gitclient.BuildUrl("https://api.github.com/repos/myrepo/myuser/issues")
			Expect(err).To(HaveOccurred())
			Expect(url).To(Equal(""))
		})
	})

	Context("when proper input url is provided", func() {
		It("should return the issue list", func() {
			client, err := gitclient.NewGitClient("git@github.com:zszabo-rh/issues-operator.git")
			Expect(err).ToNot(HaveOccurred())
			gitissues, err := client.GetIssues()
			Expect(err).ToNot(HaveOccurred())
			Expect(len(gitissues)).To(BeNumerically(">", 0))
		})
	})

	Context("when wrong url is provided", func() {
		It("should return an error", func() {
			client, err := gitclient.NewGitClient("git@github.com:zszabo/issues-operator.git")
			Expect(err).ToNot(HaveOccurred())
			gitissues, err := client.GetIssues()
			Expect(err).To(Equal(fmt.Errorf("%v", "Not Found")))
			Expect(gitissues).To(BeNil())
		})
	})

	Context("when no token is set", func() {
		BeforeEach(func() {
			originalToken = os.Getenv("GITTOKEN")
			os.Setenv("GITTOKEN", "abc123")
		})

		AfterEach(func() {
			if originalToken == "" {
				os.Unsetenv("GITTOKEN")
			} else {
				os.Setenv("GITTOKEN", originalToken)
			}
		})

		It("should return an error", func() {
			client, err := gitclient.NewGitClient("git@github.com:zszabo-rh/issues-operator.git")
			Expect(err).ToNot(HaveOccurred())
			gitissues, err := client.GetIssues()
			Expect(err).To(Equal(fmt.Errorf("%v", "Unauthorized")))
			Expect(gitissues).To(BeNil())
		})
	})

	Context("when new issue title is provided for AddIssue", func() {
		BeforeEach(func() {
		})

		It("should create the github issue", func() {
			client, err := gitclient.NewGitClient("git@github.com:zszabo-rh/issues-operator.git")
			Expect(err).ToNot(HaveOccurred())
			issueTitle := fmt.Sprintf("Generated_test_issue_%v", time.Now().Format("2006-01-02T15:04:05Z"))
			gitissue, err := client.AddIssue(issueTitle, "description")
			issueId = gitissue.Id
			Expect(err).ToNot(HaveOccurred())
			Expect(gitissue.Title).To(Equal(issueTitle))
		})
	})

	Context("when existing issue ID is provided for UpdateIssue", func() {
		It("should update the title and description", func() {
			client, err := gitclient.NewGitClient("git@github.com:zszabo-rh/issues-operator.git")
			Expect(err).ToNot(HaveOccurred())
			issueTitle := fmt.Sprintf("Generated_test_issue_%v", time.Now().Format("2006-01-02T15:04:05Z"))
			gitissue, err := client.UpdateIssue(issueId, issueTitle, "new description")
			Expect(err).ToNot(HaveOccurred())
			Expect(gitissue.Title).To(Equal(issueTitle))
			Expect(gitissue.Description).To(Equal("new description"))
		})
	})

	Context("when non-existing issue ID is provided for UpdateIssue", func() {
		It("should return an error", func() {
			client, err := gitclient.NewGitClient("git@github.com:zszabo-rh/issues-operator.git")
			Expect(err).ToNot(HaveOccurred())
			issueTitle := fmt.Sprintf("Generated_test_issue_%v", time.Now().Format("2006-01-02T15:04:05Z"))
			gitissue, err := client.UpdateIssue(999, issueTitle, "new description")
			Expect(err).To(Equal(fmt.Errorf("%v", "Not Found")))
			Expect(gitissue.Title).To(Equal(""))
		})
	})

})
