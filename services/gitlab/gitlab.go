package gitlab

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	logger "gitlab-synced-groups-lister/logging"
	services "gitlab-synced-groups-lister/services"

	"github.com/tomnomnom/linkheader"
)

type ldapGroupLink struct {
	Cn          string `json:"cn"`
	GroupAccess int    `json:"group_access"`
	Provider    string `json:"provider"`
	Filter      string `json:"filter"`
}

type gitlabGroup struct {
	Id                             int             `json:"id"`
	WebUrl                         string          `json:"web_url"`
	Name                           string          `json:"name"`
	Path                           string          `json:"path"`
	Description                    string          `json:"description"`
	Visbility                      string          `json:"visibility"`
	ShareWithGroupLock             bool            `json:"share_with_group_lock"`
	RequireTwoFactorAuthentication bool            `json:"require_two_factor_authentication"`
	TwoFactorGracePeriod           int             `json:"two_factor_grace_period"`
	ProjectCreationLevel           string          `json:"project_creation_level"`
	AutoDevopsEnabled              bool            `json:"auto_devops_enabled"`
	SubgroupCreationLevel          string          `json:"subgroup_creation_level"`
	EmailsDisabled                 bool            `json:"emails_disabled"`
	MentionsDisabled               bool            `json:"mentions_disabled"`
	LfsEnabled                     bool            `json:"lfs_enabled"`
	DefaultBranchProtection        int             `json:"default_branch_protection"`
	AvatarUrl                      string          `json:"avatar_url"`
	RequestAccessEnabled           bool            `json:"request_access_enabled"`
	FullName                       string          `json:"full_name"`
	FullPath                       string          `json:"full_path"`
	CreatedAt                      string          `json:"created_at"`
	ParentId                       int             `json:"parent_id"`
	LdapCn                         string          `json:"ldap_cn"`
	LdapAccess                     int             `json:"ldap_access"`
	LdapGroupLinks                 []ldapGroupLink `json:"ldap_group_links"`
	MarkedForDeletion              bool            `json:"marked_for_deletion"`
}

func buildGitlabUrl(baseUrl string, apiVersion string, resource string) string {
	var stringBuilder strings.Builder

	stringBuilder.WriteString(baseUrl + "/api/" + apiVersion + "/" + resource)

	logger.Log().Debugf("buildGitlabUrl: %v\n", stringBuilder.String())
	return stringBuilder.String()
}

func constructRequest(requestMethod string, baseUrl string, headers []services.Header, queries []services.Query) *http.Request {
	requestUrl, err := url.Parse(baseUrl)
	if err != nil {
		logger.Log().Fatal(err)
	}

	q := requestUrl.Query()
	for _, query := range queries {
		q.Set(query.GetKey(), query.GetValue())
	}
	requestUrl.RawQuery = q.Encode()

	req, err := http.NewRequest(requestMethod, requestUrl.String(), nil)

	if err != nil {
		logger.Log().Fatal(err)
	}
	for _, header := range headers {
		req.Header.Set(header.GetKey(), header.GetValue())
	}
	return req
}

func executeRequest(request *http.Request) *http.Response {
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		logger.Log().Fatal(err)
	}

	return response
}

func parseResponse(response *http.Response) (map[string][]string, string) {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Log().Fatal(err)
	}

	return response.Header, string(body)
}

func buildHeaders(token string) []services.Header {
	headers := []services.Header{
		{
			Key:   "PRIVATE-TOKEN",
			Value: token,
		},
	}

	return headers
}

func buildQueries(pagination string, perPage string, orderBy string, sort string) []services.Query {
	queries := []services.Query{
		{
			Key:   "pagination",
			Value: "keyset",
		},
		{
			Key:   "per_page",
			Value: "100",
		},
		{
			Key:   "order_by",
			Value: "id",
		},
		{
			Key:   "sort",
			Value: "asc",
		},
	}

	return queries
}

func GetSyncedGroups(baseUrl string, apiVersion string, token string, outputFileName string) (csvContent string) {
	headers := buildHeaders(token)
	queries := buildQueries(
		"keyset",
		"100",
		"id",
		"asc",
	)
	request := constructRequest("GET", buildGitlabUrl(baseUrl, apiVersion, "groups"), headers, queries)
	logger.Log().Debugf("request: %v\n", request)
	response := executeRequest(request)
	logger.Log().Debugf("response: %v\n", response)

	var responses []gitlabGroup

	initialRespHeaders, initialRespBody := parseResponse(response)
	json.Unmarshal([]byte(initialRespBody), &responses)
	logger.Log().Debugf("responses: %v\n", responses)

	responseHeaders := initialRespHeaders
	for {
		if linkHeader, ok := responseHeaders["Link"]; ok {
			var nextLink string
			nextLinks := linkheader.Parse(strings.Join(linkHeader[:], ","))
			for _, link := range nextLinks {
				if link.Rel == "next" {
					nextLink = link.URL
				}
			}
			logger.Log().Debugf("nextLink: %s\n", nextLink)
			if nextLink == "" {
				break
			}

			request := constructRequest("GET", nextLink, headers, nil)
			response := executeRequest(request)

			tempHeaders, responseBody := parseResponse(response)
			responseHeaders = tempHeaders

			var gitlabGroups []gitlabGroup
			json.Unmarshal([]byte(responseBody), &gitlabGroups)
			responses = append(responses, gitlabGroups...)

			if _, ok := responseHeaders["Link"]; ok {
				continue
			} else {
				break
			}
		} else {
			break
		}
	}

	csvFile, err := os.Create(outputFileName)
	if err != nil {
		logger.Log().Fatal(err)
	}
	defer csvFile.Close()

	records := [][]string{
		{"gitlab_group", "ldap_group_links"},
	}

	for _, entry := range responses {
		fullPath := entry.FullPath
		var stringBuilder strings.Builder
		for idx, ldapGroupLink := range entry.LdapGroupLinks {
			cn := ldapGroupLink.Cn
			if cn != "" {
				stringBuilder.WriteString(cn)
			} else {
				continue
			}

			if idx < len(entry.LdapGroupLinks)-1 {
				stringBuilder.WriteString(",")
			}
		}
		ldapGroupLinks := stringBuilder.String()
		record := []string{fullPath, ldapGroupLinks}
		records = append(records, record)
		logger.Log().Infof("%s", record)
	}

	buf := new(bytes.Buffer)
	stringWriter := csv.NewWriter(buf)

	err = stringWriter.WriteAll(records)
	if err != nil {
		logger.Log().Fatal(err)
	}
	csvContent = buf.String()

	return csvContent
}
