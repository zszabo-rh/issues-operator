package gitclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

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

func GetToken() string {
	var s Env
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
		return ""
	}
	return s.GitToken
}

func GetUrl(repo string) string {
	// Example repo: git@github.com:zszabo-rh/issues-operator.git
	// Path:         zszabo-rh/issues-operator
	// Output:       https://api.github.com/repos/zszabo-rh/issues-operator/issues
	split1 := strings.Split(repo, ":")
	split2 := strings.Split(split1[1], ".")
	url := "https://api.github.com/repos/" + split2[0] + "/issues"
	return url
}

func GetIssues(repo string) ([]GitIssue, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", GetUrl(repo), nil)

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	var gitissues []GitIssue
	err = json.Unmarshal([]byte(body), &gitissues)
	defer resp.Body.Close()
	return gitissues, err
}

func AddIssue(repo string, title string, desc string) (GitIssue, error) {
	gitissue := GitIssue{
		Title:       title,
		Description: desc}

	gitissueJson, err := json.Marshal(gitissue)
	if err != nil {
		log.Fatal(err.Error())
		return GitIssue{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", GetUrl(repo), bytes.NewBuffer(gitissueJson))

	if err != nil {
		log.Fatal(err.Error())
		return GitIssue{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+GetToken())

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err.Error())
		return GitIssue{}, err
	}
	defer resp.Body.Close()

	err = parseResponse(resp, &gitissue)
	if err != nil {
		log.Fatal(err.Error())
		return GitIssue{}, err
	}

	return gitissue, err
}

func parseResponse(resp *http.Response, output *GitIssue) error {

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	//return json.Unmarshal([]byte(body), output)

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	structValue := reflect.ValueOf(output).Elem()
	for _, fieldName := range []string{"Title", "Description", "Id", "Status", "LastUpdated"} {
		if field, ok := structValue.Type().FieldByName(fieldName); ok && field.Tag != "" {
			jsonTag := field.Tag.Get("json")
			if outputValue, ok := data[jsonTag]; ok {
				outputType := structValue.FieldByName(fieldName).Type()
				outputValue := reflect.ValueOf(outputValue)
				if outputValue.Type().ConvertibleTo(outputType) {
					structValue.FieldByName(fieldName).Set(outputValue.Convert(outputType))
				}
			}
		}
	}
	return nil
}

func UpdateIssue(repo string, Id int, title string, desc string) (GitIssue, error) {
	gitissue := GitIssue{
		Title:       title,
		Description: desc}

	gitissueJson, err := json.Marshal(gitissue)
	if err != nil {
		log.Fatal(err.Error())
		return GitIssue{}, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PATCH", GetUrl(repo)+"/"+fmt.Sprint(Id), bytes.NewBuffer(gitissueJson))

	if err != nil {
		log.Fatal(err.Error())
		return GitIssue{}, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+GetToken())

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err.Error())
		return GitIssue{}, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err.Error())
		return GitIssue{}, err
	}

	err = json.Unmarshal([]byte(body), &gitissue)
	defer resp.Body.Close()
	return gitissue, err
}
