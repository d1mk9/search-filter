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

type TemplateCtx struct {
	Now         time.Time
	Loc         *time.Location
	CurrentUser int64
}

var (
	reTodayMinus = regexp.MustCompile(`\{\{\s*today-(\d+)d\s*\}\}`)
	reCurrUser   = regexp.MustCompile(`\{\{\s*current_user\s*\}\}`)
)

func normalizePlaceholders(in string) string {
	in = reTodayMinus.ReplaceAllString(in, "{{ daysAgo $1 }}")
	in = reCurrUser.ReplaceAllString(in, "{{ currentUser }}")
	return in
}

func RenderTemplate(input string, ctx TemplateCtx) ([]byte, error) {
	now := ctx.Now
	if ctx.Loc != nil {
		now = now.In(ctx.Loc)
	}

	funcs := template.FuncMap{
		"today":       func() string { return now.Format("2006-01-02") },
		"daysAgo":     func(n int) string { return now.AddDate(0, 0, -n).Format("2006-01-02") },
		"currentUser": func() int64 { return ctx.CurrentUser },
	}

	src := normalizePlaceholders(input)

	tpl, err := template.New("query").
		Funcs(funcs).
		Option("missingkey=error").
		Parse(src)
	if err != nil {
		return nil, fmt.Errorf("template parse: %w", err)
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, nil); err != nil {
		return nil, fmt.Errorf("template execute: %w", err)
	}

	if !json.Valid(buf.Bytes()) {
		return nil, fmt.Errorf("template produced invalid JSON")
	}

	return buf.Bytes(), nil
}

func RenderQuery(q types.Query, ctx TemplateCtx) (types.Query, error) {
	raw, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	out, err := RenderTemplate(string(raw), ctx)
	if err != nil {
		return nil, err
	}
	var res types.Query
	if err := json.Unmarshal([]byte(out), &res); err != nil {
		return nil, err
	}
	return res, nil
}
