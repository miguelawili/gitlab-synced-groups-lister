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

	req := services.ConstructRequest("GET", services.BuildUrl(c.service.Url, []string{"content", pageId}), requestHeaders, requestParams)
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

func (c *Confluence) UpdatePageContent(headers []services.Header, queries []services.Query, pageId string, currentVersion int, body string) {
	// requestHeaders := append(c.baseHeaders, headers...)
	// requestParams := append(c.baseQueries, queries...)

	htmlTable := services.CsvToHtmlTable(body)

	logger.Log().Debugf("htmlTable:\n%s", htmlTable)
}
