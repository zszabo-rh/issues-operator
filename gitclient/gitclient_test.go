// Assisted by watsonx Code Assistant

package gitclient

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestGetUrl(t *testing.T) {
	repo := "git@github.com:zszabo-rh/issues-operator.git"
	expectedUrl := "https://api.github.com/repos/zszabo-rh/issues-operator/issues"
	RegisterTestingT(t)
	Expect(GetUrl(repo)).To(Equal(expectedUrl))
}

func TestGetIssues(t *testing.T) {
	repo := "git@github.com:zszabo-rh/issues-operator.git"
	RegisterTestingT(t)
	gitissues, err := GetIssues(repo)
	Expect(err).To(BeNil())
	Expect(len(gitissues)).To(BeNumerically(">", 0))
}

func TestAddIssue(t *testing.T) {
	repo := "git@github.com:zszabo-rh/issues-operator.git"
	title := "Test Issue"
	desc := "This is a test issue"
	RegisterTestingT(t)
	gitissue, err := AddIssue(repo, title, desc)
	Expect(err).To(BeNil())
	Expect(gitissue.Title).To(Equal(title))
	Expect(gitissue.Description).To(Equal(desc))
}

func TestUpdateIssue(t *testing.T) {
	repo := "git@github.com:zszabo-rh/issues-operator.git"
	title := "Updated Test Issue"
	desc := "This is an updated test issue"
	gitissue, err := AddIssue(repo, "Test Issue", "This is a test issue")
	RegisterTestingT(t)
	Expect(err).To(BeNil())
	updatedGitissue, err := UpdateIssue(repo, gitissue.Id, title, desc)
	Expect(err).To(BeNil())
	Expect(updatedGitissue.Title).To(Equal(title))
	Expect(updatedGitissue.Description).To(Equal(desc))
}
