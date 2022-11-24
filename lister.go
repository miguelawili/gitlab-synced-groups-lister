package main

import (
	logger "gitlab-synced-groups-lister/logging"
	services "gitlab-synced-groups-lister/services"
	confluence "gitlab-synced-groups-lister/services/confluence"
	gitlab "gitlab-synced-groups-lister/services/gitlab"
	settings "gitlab-synced-groups-lister/settings"

	"strconv"
)

func main() {
	logger.InitLogger("./config/settings.toml")
	conf := settings.LoadConfig("./config/settings.toml")
	logger.Log().Debugf("Configuration loaded!\n%v", conf)

	outputFileName := conf.Output.FileName

	gitlabBaseUrl := conf.Gitlab.Url
	gitlabApiVersion := conf.Gitlab.ApiVersion
	gitlabToken := conf.Gitlab.Token

	csvContent := gitlab.GetSyncedGroups(gitlabBaseUrl, gitlabApiVersion, gitlabToken, outputFileName)

	if csvContent == "" {
		logger.Log().Fatal("CSV content is empty!")
	}

	// Confluence
	var confluenceBaseHeaders []services.Header
	confluenceToken := services.EncodeB64(conf.Confluence.User + ":" + conf.Confluence.Token)
	confluenceBaseHeaders = append(confluenceBaseHeaders, services.NewHeader("Authorization", "Basic "+confluenceToken))

	var confluenceBaseQueries []services.Query

	confluenceService := confluence.New(
		conf.Confluence.Url,
		confluenceBaseHeaders,
		confluenceBaseQueries,
	)

	content := confluenceService.GetPage(
		[]services.Header{},
		[]services.Query{
			{
				Key:   "limit",
				Value: strconv.Itoa(2),
			},
			{
				Key:   "expand",
				Value: "version,body.storage,space",
			},
		},
		conf.Confluence.PageId,
	)
	currentVersion := content.Version.Number

	confluenceService.UpdatePageContent(
		[]services.Header{},
		[]services.Query{
			{
				Key:   "limit",
				Value: strconv.Itoa(2),
			},
			{
				Key:   "expand",
				Value: "version,body.storage,space",
			},
		},
		conf.Confluence.PageId,
		currentVersion+1,
		csvContent,
	)

	logger.Log().Debugf("content:\n%+v", content)
}
