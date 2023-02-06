// SPDX-License-Identifier: AGPL-3.0-only
package bot

import (
	"embed"
	"fmt"
	"net/http"
	"text/template"
)

//go:embed style/*
var styles embed.FS

var staticServer = http.StripPrefix("/static/", http.FileServer(http.FS(styles)))

//go:embed text.html
var textHtml string

var textTmpl = template.Must(template.New("").Parse(textHtml))

func (bot *Bot) ShowText(w http.ResponseWriter, title, header, text string, v ...interface{}) {
	var data = struct {
		Title, Header, Text string
	}{title, header, fmt.Sprintf(text, v...)}

	if data.Header == "" {
		data.Header = data.Title
	}

	err := textTmpl.Execute(w, data)
	if err != nil {
		bot.Sugar.Errorf("error executing text template: %v", err)
	}
}
