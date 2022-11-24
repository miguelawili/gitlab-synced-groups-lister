package services

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"strings"
)

func EncodeB64(data string) string {
	return b64.StdEncoding.EncodeToString([]byte(data))
}

func CsvToHtmlTable(data [][]string) string {
	var sb strings.Builder

	sb.WriteString("<table>")
	sb.WriteString("<tr>")
	sb.WriteString("<th>Gitlab Group</th>")
	sb.WriteString("<th>LDAP Group Links</th>")
	sb.WriteString("</tr>")
	for _, line := range data {
		sb.WriteString("<tr>")

		var gitlabGroup string = line[0]
		sb.WriteString("<td>")
		sb.WriteString(gitlabGroup)
		sb.WriteString("</td>")

		if len(line) < 2 {
			continue
		}
		var syncedGroups string = line[1]
		sb.WriteString("<td>")
		sb.WriteString(syncedGroups)
		sb.WriteString("</td>")

		sb.WriteString("</tr>")
	}
	sb.WriteString("</table>")

	return sb.String()
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	b := bytes.TrimRight(buffer.Bytes(), "\n")
	return b, err
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if c := s[len(s)-1]; s[0] == c && (c == '"' || c == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
