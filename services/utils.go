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
	for idx, line := range lines {
		sb.WriteString("<tr>")
		sb.WriteString()

		if idx == len(lines)-1 {
			sb.WriteString("</table>")
		}
	}
}
