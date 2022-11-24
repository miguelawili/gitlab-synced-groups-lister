package services

import (
	b64 "encoding/base64"
	"strings"
)

func EncodeB64(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

func CsvToHtmlTable(data string) string {
	var sb strings.Builder

	lines := strings.Split(data, "\n")

	sb.WriteString("<table>")
	sb.WriteString("<tr>")
	sb.WriteString("<th>Gitlab Group</th>")
	sb.WriteString("<th>LDAP Group Links</th>")
	sb.WriteString("</tr>")
	for idx, line := range lines {
		lineItems := strings.Split(line, ",")
		var gitlabGroup string = lineItems[0]
		var groupsSynced string = strings.Join(lineItems[1:len(lineItems)-1], ",")

		sb.WriteString("<tr>")
		sb.WriteString("<td>")
		sb.WriteString(gitlabGroup)
		sb.WriteString("</td>")
		sb.WriteString("<td>")
		sb.WriteString(groupsSynced)
		sb.WriteString("</td>")
		sb.WriteString("</tr>")

		if idx == len(lines)-1 {
			sb.WriteString("</table>")
		}
	}

	return sb.String()
}
