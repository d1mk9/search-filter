package placeholder

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"text/template"
	"time"
)

type TemplateCtx struct {
	Now         time.Time
	Loc         *time.Location
	CurrentUser int64
}

var (
	reTodayMinus = regexp.MustCompile(`\{\{\s*today-(\d+)d\s*\}\}`)
	reCurrUser   = regexp.MustCompile(`\{\{\s*current_user\s*\}\}`)
)

func normalizePlaceholders(in []byte) []byte {
	s := string(in)
	s = reTodayMinus.ReplaceAllString(s, "{{ daysAgo $1 }}")
	s = reCurrUser.ReplaceAllString(s, "{{ currentUser }}")
	return []byte(s)
}

func RenderTemplate(input []byte, ctx TemplateCtx) (json.RawMessage, error) {
	now := ctx.Now
	if ctx.Loc != nil {
		now = now.In(ctx.Loc)
	}

	funcs := template.FuncMap{
		"today": func() string {
			return now.Format("2006-01-02")
		},
		"daysAgo": func(n int) string {
			return now.AddDate(0, 0, -n).Format("2006-01-02")
		},
		"currentUser": func() int64 {
			return ctx.CurrentUser
		},
	}

	src := normalizePlaceholders(input)

	tpl, err := template.New("query").
		Funcs(funcs).
		Option("missingkey=error").
		Parse(string(src))
	if err != nil {
		return nil, fmt.Errorf("template parse: %w", err)
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, nil); err != nil {
		return nil, fmt.Errorf("template execute: %w", err)
	}

	var compact bytes.Buffer
	if err := json.Compact(&compact, out.Bytes()); err != nil {
		return nil, fmt.Errorf("template produced invalid JSON: %w", err)
	}
	return json.RawMessage(compact.Bytes()), nil
}
