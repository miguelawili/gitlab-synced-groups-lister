package confluence

import (
	"encoding/json"

	logger "gitlab-synced-groups-lister/logging"
	"gitlab-synced-groups-lister/services"
)

type Confluence struct {
	service     *services.HTTPService
	baseHeaders []services.Header
	baseQueries []services.Query
}

func New(url string, headers []services.Header, queries []services.Query) *Confluence {
	return &Confluence{
		service:     services.NewHttpService(url),
		baseHeaders: headers,
		baseQueries: queries,
	}
}

func (c *Confluence) GetPage(headers []services.Header, queries []services.Query, pageId string) GetPageContent {
	requestHeaders := append(c.baseHeaders, headers...)
	requestParams := append(c.baseQueries, queries...)

	req := services.ConstructRequest("GET", services.BuildUrl(c.service.Url, []string{"content", pageId}), requestHeaders, requestParams, nil)
	logger.Log().Debugf("req:\n%v", req)

	res, err := c.service.Client.Do(req)
	if err != nil && res.StatusCode != 200 {
		logger.Log().Errorf("ERR:\n%v")
	}

	_, resBody := services.ParseHttpResponse(res)

	var content GetPageContent
	err = json.Unmarshal([]byte(resBody), &content)
	if err != nil {
		logger.Log().Fatalf("Error deserializing JSON output:\n%v", err)
	}

	return content
}

func (c *Confluence) UpdatePageContent(headers []services.Header, queries []services.Query, pageId string, pageTitle string, currentVersion int, records [][]string) bool {
	requestHeaders := append(c.baseHeaders, headers...)
	requestParams := append(c.baseQueries, queries...)

	var htmlTable string = services.CsvToHtmlTable(records)

	input := &UpdatePagePayload{
		Title: pageTitle,
		Type:  "page",
		Version: VersionPayload{
			Number:    currentVersion + 1,
			MinorEdit: false,
		},
		Body: BodyPayload{
			Storage: ViewPayload{
				Value:          htmlTable,
				Representation: "storage",
			},
		},
	}
	payload, err := services.JSONMarshal(input)
	if err != nil {
		logger.Log().Errorf("Error converting to JSON str:\n%v", err)
	}

	logger.Log().Debugf("htmlTable:\n%s", htmlTable)

	req := services.ConstructRequest("PUT", services.BuildUrl(c.service.Url, []string{"content", pageId}), requestHeaders, requestParams, payload)
	logger.Log().Debugf("req:\n%v", req)

	res, err := c.service.Client.Do(req)
	if err != nil && res.StatusCode != 200 {
		logger.Log().Errorf("ERR:\n%v")
	}

	_, resBody := services.ParseHttpResponse(res)

	logger.Log().Debugf("resBody:\n%v", resBody)

	return true
}
