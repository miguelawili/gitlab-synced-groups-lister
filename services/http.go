package services

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	logger "gitlab-synced-groups-lister/logging"
)

type HTTPService struct {
	Url    string
	Client *http.Client
}

func NewHttpService(url string) *HTTPService {
	client := &http.Client{}

	return &HTTPService{
		Url:    url,
		Client: client,
	}
}

func ConstructRequest(requestMethod string, baseUrl string, headers []Header, queries []Query, body []byte) *http.Request {
	requestUrl, err := url.Parse(baseUrl)
	if err != nil {
		logger.Log().Fatal(err)
	}

	if len(queries) > 0 {
		q := requestUrl.Query()
		for _, query := range queries {
			q.Set(query.GetKey(), query.GetValue())
		}
		requestUrl.RawQuery = q.Encode()
	}

	var reqBody *bytes.Buffer
	if body == nil {
		reqBody = nil
	}
	reqBody = bytes.NewBuffer(body)

	req, err := http.NewRequest(requestMethod, requestUrl.String(), reqBody)
	if err != nil {
		logger.Log().Fatal(err)
	}

	if len(headers) > 0 {
		for _, header := range headers {
			req.Header.Set(header.GetKey(), header.GetValue())
		}
	}

	return req
}

func BuildUrl(baseUrl string, path []string) string {
	var sb strings.Builder

	sb.WriteString(baseUrl)
	for i, entry := range path {
		sb.WriteString(entry)

		if i < len(path)-1 {
			sb.WriteString("/")
		}
	}

	return sb.String()
}

func ParseHttpResponse(response *http.Response) (map[string][]string, string) {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Log().Fatal(err)
	}

	return response.Header, string(body)
}
