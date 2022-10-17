package scraper

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/margostino/metabook/common"
	"strings"
)

var STOPWORDS = []string{
	"facebook",
	"twitter",
	"instagram",
}

func Collect() map[string]string {
	var data = make(map[string]string)
	url := "https://en.wikipedia.org/wiki/The_Lord_of_the_Rings"
	//url := "https://en.wikipedia.org/wiki/Dune_(novel)"
	browser := rod.New().MustConnect()
	defer browser.Close()
	page := browser.MustPage(url).MustWaitLoad()
	bodyElements, err := page.Elements(".infobox.vcard")
	common.Check(err)

	for _, element := range bodyElements {
		text, _ := element.Text()
		lines := common.NewString(text).
			ReplaceAll("\n\n", "\n").
			Split("\n").
			Value()

		normalizedLines := normalizeLines(lines)

		for _, line := range normalizedLines {
			parts := strings.Split(line, "\t")
			key := common.NewString(parts[0]).
				ToLower().
				ReplaceAll("\u00a0", " ").
				ReplaceAll(" ", "_").
				Value()
			value := common.NewString(parts[1]).
				ReplaceAll("[1]", "").
				Value()

			data[key] = value
		}
	}

	return data
}

func normalizeLines(lines []string) []string {
	var normalizedLines = make([]string, 0)
	var key, value string

	for i, line := range lines {
		if i == 0 || strings.Contains(line, "\t") {
			parts := strings.Split(line, "\t")

			if len(parts) == 2 && parts[1] == "" {
				key = parts[0]
			} else if len(parts) == 2 {
				normalizedLines = append(normalizedLines, line)
			} else if len(parts) == 1 {
				titleLine := fmt.Sprintf("title\t%s", line)
				normalizedLines = append(normalizedLines, titleLine)
			}

			if key != "" && value != "" {
				keyValue := fmt.Sprintf("%s\t%s", key, value)
				normalizedLines = append(normalizedLines, keyValue)
				key = ""
				value = ""
			}

		} else if i != 1 && !strings.Contains(line, "\t") {
			if value == "" {
				value = line
			} else {
				value = fmt.Sprintf("%s, %s", value, line)
			}
		}
	}

	return normalizedLines
}
