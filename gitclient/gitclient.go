package gitclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type GitClient struct {
	repo  string
	token string
}

type GitIssue struct {
	Title       string `json:"title"`
	Description string `json:"body"`
	Status      string `json:"state"`
	Id          int    `json:"number"`
	LastUpdated string `json:"updated_at"`
}

type Env struct {
	GitToken string `required:"true" envconfig:"gittoken"`
}

func GetToken() (string, error) {
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		return "", err
	}
	return env.GitToken, nil
}

func BuildUrl(repo string) (string, error) {
	split1 := strings.Split(repo, ":")
	if len(split1) != 2 {
		return "", fmt.Errorf("invalid repo format")
	}
	split2 := strings.Split(split1[1], ".")
	if len(split2) != 2 {
		return "", fmt.Errorf("invalid repo format")
	}
	url := "https://api.github.com/repos/" + split2[0] + "/issues"
	return url, nil
}

func NewGitClient(repo string) (*GitClient, error) {
	httprepo, err := BuildUrl(repo)
	if err != nil {
		return nil, err
	}
	token, err := GetToken()
	if err != nil {
		return nil, err
	}
	g := GitClient{httprepo, token}
	return &g, nil
}

func (g *GitClient) GetIssues() ([]GitIssue, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", g.repo, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+g.token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return nil, fmt.Errorf("%v", http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var gitissues []GitIssue
	err = json.Unmarshal([]byte(body), &gitissues)
	if err != nil {
		return nil, err
	}
	return gitissues, nil
}

func (g *GitClient) AddIssue(title string, desc string) (GitIssue, error) {
	gitissue := GitIssue{
		Title:       title,
		Description: desc}

	gitissueJson, err := json.Marshal(gitissue)
	if err != nil {
		return GitIssue{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", g.repo, bytes.NewBuffer(gitissueJson))
	if err != nil {
		return GitIssue{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+g.token)

	resp, err := client.Do(req)
	if err != nil {
		return GitIssue{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GitIssue{}, err
	}

	err = json.Unmarshal([]byte(body), &gitissue)
	return gitissue, err
}

func (g *GitClient) UpdateIssue(Id int, title string, desc string) (GitIssue, error) {
	gitissue := GitIssue{
		Title:       title,
		Description: desc}

	gitissueJson, err := json.Marshal(gitissue)
	if err != nil {
		return GitIssue{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", g.repo+"/"+fmt.Sprint(Id), bytes.NewBuffer(gitissueJson))

	if err != nil {
		return GitIssue{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+g.token)

	resp, err := client.Do(req)
	if err != nil {
		return GitIssue{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return GitIssue{}, fmt.Errorf("%v", http.StatusText(resp.StatusCode))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GitIssue{}, err
	}

	err = json.Unmarshal([]byte(body), &gitissue)
	return gitissue, err
}
