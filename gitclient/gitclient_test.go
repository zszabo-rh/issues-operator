package gitclient

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestTestGetUrl(t *testing.T) {
	g := NewWithT(t)
	b := GetUrl("git@github.com:zszabo-rh/issues-operator.git")
	g.Expect(b).To(Equal("https://api.github.com/repos/zszabo-rh/issues-operator/issues"))
}

func TestGetToken(t *testing.T) {
	g := NewWithT(t)
	b := GetToken()
	g.Expect(b).To(Equal("blabla"))
}

/*
	func TestCreateNew(t *testing.T) {
		g := NewWithT(t)
		b, err := AddIssue("git@github.com:zszabo-rh/issues-operator.git", "New Issue", "Created by gomega")
		if err != nil {
			panic(err)
		}
		g.Expect(b > 0).To(BeTrue())
	}
*/
func TestGetIssues(t *testing.T) {
	g := NewWithT(t)
	b, err := GetIssues("git@github.com:zszabo-rh/issues-operator.git")
	if err != nil {
		panic(err)
	}
	g.Expect(len(b) == 7).To(BeTrue())
}

func TestUpdateIssue(t *testing.T) {
	g := NewWithT(t)
	title := "Updated title"
	b, err := UpdateIssue("git@github.com:zszabo-rh/issues-operator.git", 3, title, "Updated description")
	if err != nil {
		panic(err)
	}
	g.Expect(b.Title).To(Equal(title))
}
