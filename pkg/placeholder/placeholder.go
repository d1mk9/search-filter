package placeholder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"search-filter/pkg/types"
	"text/template"
	"time"
)

var (
	reTodayMinus = regexp.MustCompile(`\{\{\s*today-(\d+)d\s*\}\}`)
	reCurrUser   = regexp.MustCompile(`\{\{\s*current_user\s*\}\}`)
)

func normalizePlaceholders(in string) string {
	in = reTodayMinus.ReplaceAllString(in, "{{ daysAgo $1 }}")
	in = reCurrUser.ReplaceAllString(in, "{{ currentUser }}")
	return in
}

func RenderTemplate(input string, now time.Time, loc *time.Location, currentUser int64) ([]byte, error) {
	funcs := template.FuncMap{
		"today":       func() string { return now.Format("2006-01-02") },
		"daysAgo":     func(n int) string { return now.AddDate(0, 0, -n).Format("2006-01-02") },
		"currentUser": func() int64 { return currentUser },
	}

	tpl, err := template.New("query").
		Funcs(funcs).
		Option("missingkey=error").
		Parse(normalizePlaceholders(input))
	if err != nil {
		return nil, fmt.Errorf("template parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, nil); err != nil {
		return nil, fmt.Errorf("template execute: %w", err)
	}

	return buf.Bytes(), nil
}

func RenderQuery(q types.Query, now time.Time, loc *time.Location, currentUser int64) (types.Query, error) {
	raw, err := json.Marshal(q)
	if err != nil {
		return nil, fmt.Errorf("marshal query: %w", err)
	}

	out, err := RenderTemplate(string(raw), now, loc, currentUser)
	if err != nil {
		return nil, err
	}

	var res types.Query
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, fmt.Errorf("unmarshal rendered query: %w", err)
	}
	return res, nil
}
