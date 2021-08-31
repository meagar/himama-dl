package himama

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"
)

type Activity struct {
	AddedBy  string
	Date     string
	Title    string
	MediaURL string
}

var nameRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)

func (a *Activity) SuggestedLocalFilename() string {
	dateParts := strings.Split(a.Date, "/")
	date := fmt.Sprintf("20%s-%s-%s", dateParts[2], zeroPad(dateParts[0]), zeroPad(dateParts[1]))

	hash := sha1.New()
	hash.Write([]byte(a.MediaURL))
	hashStr := hex.EncodeToString(hash.Sum(nil))

	nameParts := []string{date, a.AddedBy, a.Title, hashStr[0:8]}
	for i := range nameParts {
		nameParts[i] = strings.TrimSpace(nameRegex.ReplaceAllString(nameParts[i], " "))
	}

	parsedURL, err := url.Parse(a.MediaURL)
	if err != nil {
		panic(err)
	}

	ext := strings.ToLower(filepath.Ext(parsedURL.Path))
	return strings.Join(nameParts, " - ") + ext
}

func zeroPad(str string) string {
	if len(str) < 2 {
		return "0" + str
	}
	return str
}
