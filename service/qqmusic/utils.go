// /qqmusic/utils.go
package qqmusic

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

// checkNum 检查字符串是否为纯数字
var isNumeric = regexp.MustCompile(`^\d+$`).MatchString

// resolveRespJson 处理 JSONP 响应
func resolveRespJson(callBackSign, val string) string {
	if !strings.HasPrefix(val, callBackSign) {
		return ""
	}
	jsonStr := strings.TrimPrefix(val, callBackSign+"(")
	if strings.HasSuffix(jsonStr, ")") {
		return jsonStr[:len(jsonStr)-1]
	}
	return jsonStr
}

// 包含歌词内容的 XML 结构
type LyricContentXML struct {
	LyricContent string `xml:",chardata"`
}

// recursionFindElement 递归查找 XML 元素并返回其内部文本
// Go 的标准库没有直接的 DOM 式 API，所以我们用流式解析器来实现类似功能。
func recursionFindElement(xmlData string, mapping map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	decoder := xml.NewDecoder(strings.NewReader(strings.TrimSpace(xmlData)))

	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if se, ok := token.(xml.StartElement); ok {
			if goKey, found := mapping[se.Name.Local]; found {
				var content string
				if err := decoder.DecodeElement(&content, &se); err != nil {
					return nil, err
				}
				result[goKey] = content
			}
		}
	}
	return result, nil
}

// parseLyricContentXML 解析内嵌的歌词XML
func parseLyricContentXML(xmlData string) (string, error) {
	if !strings.Contains(xmlData, "<?xml") {
		return xmlData, nil
	}
	decoder := xml.NewDecoder(strings.NewReader(xmlData))
	for {
		token, err := decoder.Token()
		if err != nil {
			if err == io.EOF {
				return "", errors.New("lyric_1 tag not found")
			}
			return "", err
		}
		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "Lyric_1" {
			for _, attr := range se.Attr {
				if attr.Name.Local == "LyricContent" {
					return attr.Value, nil
				}
			}
		}
	}
}

func TimeFormat(ms int) string {
	seconds := ms / 1000
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	milliseconds := ms % 1000
	return fmt.Sprintf("%02d:%02d.%02d", minutes, remainingSeconds, milliseconds/10)
}

func FormatLyric(lrc string) string {
	if lrc == "" {
		return ""
	}
	var formattedLyrics strings.Builder

	// Regex to capture lines with start time and duration
	lineRegex := regexp.MustCompile(`\[(\d+),(\d+)\](.*)`)
	// Regex to capture individual words with their start and end times within a line
	wordRegex := regexp.MustCompile(`(.*?)\((\d+),(\d+)\)`)

	lines := strings.Split(lrc, "\n")

	currentTime := 0

	for _, line := range lines {
		if line == "" {
			continue
		}

		matches := lineRegex.FindStringSubmatch(line)
		if len(matches) == 4 {
			startTime, _ := strconv.Atoi(matches[1])
			duration, _ := strconv.Atoi(matches[2])
			content := matches[3]

			// If the current line's startTime is greater than currentTime,
			// it means there's a gap, so we update currentTime to the line's startTime.
			// This handles cases where a line starts later than the previous one ended.
			if startTime > currentTime {
				currentTime = startTime
			}

			// Process the words within the content
			wordMatches := wordRegex.FindAllStringSubmatch(content, -1)
			var words []string
			for _, wm := range wordMatches {
				if len(wm) == 4 {
					words = append(words, wm[1])
				}
			}

			// Construct the formatted lyric line
			formattedLine := fmt.Sprintf("[%s]%s\n", TimeFormat(currentTime), strings.Join(words, ""))
			formattedLyrics.WriteString(formattedLine)

			// Update currentTime by adding the duration of the current line
			currentTime += duration
		} else {
			// If it's not a timed line (like ti, ar, al, by, offset), print it directly
			formattedLyrics.WriteString(line + "\n")
		}
	}

	return formattedLyrics.String()
}
