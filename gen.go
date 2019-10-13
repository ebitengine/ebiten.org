// Copyright 2019 Hajime Hoshi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	url         = "https://ebiten.org"
	description = "Ebiten is an open-source game library, with which you can develop 2D games with simple API for multi platforms in the Go programming language."
)

var reTitle = regexp.MustCompile(`<h1>([^<]+)</h1>`)

func cleanup() error {
	return filepath.Walk("docs", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}
		if strings.HasPrefix(filepath.Base(path), "_") {
			return nil
		}
		if err := os.Remove(path); err != nil {
			return err
		}
		return nil
	})
}

func run() error {
	tmpl, err := template.New("tmpl.html").Funcs(template.FuncMap{
		"noescape": func(str string) template.HTML {
			return template.HTML(str)
		},
	}).ParseFiles("tmpl.html")
	if err != nil {
		return err
	}

	if err := cleanup(); err != nil {
		return err
	}

	if err := filepath.Walk("contents", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		content := ""

		switch filepath.Ext(path) {
		case ".html":
			if filepath.Base(path) == "tmpl.html" {
				return nil
			}

			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			content = string(c)

		case ".json":
			t, err := template.ParseFiles(filepath.Join(filepath.Dir(path), "tmpl.html"))
			if err != nil {
				return err
			}
			var j map[string]interface{}

			c, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			if err := json.Unmarshal(c, &j); err != nil {
				return err
			}
			j["Base"] = filepath.Base(path[:len(path)-len(filepath.Ext(path))])

			b := &bytes.Buffer{}
			if err := t.Execute(b, j); err != nil {
				return err
			}
			content = string(b.Bytes())

		default:
			return nil
		}

		rel, err := filepath.Rel("contents", path[:len(path)-len(filepath.Ext(path))]+".html")
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join("docs", filepath.Dir(rel)), 0755); err != nil {
			return err
		}
		// TODO: What if the file already exists?
		w, err := os.Create(filepath.Join("docs", rel))
		if err != nil {
			return err
		}
		defer w.Close()

		title := "Ebiten - A dead simple 2D game library in Go"
		if path != filepath.Join("contents", "index.html") {
			m := reTitle.FindStringSubmatch(content)
			title = fmt.Sprintf("%s - Ebiten", html.UnescapeString(m[1]))
		}

		nav := false
		feedback := false
		if path != filepath.Join("contents", "404.html") {
			nav = true
			feedback = true
		}

		canonical := ""
		switch rel {
		case "404.html":
			// No canonical URL
		case "index.html":
			canonical = url
		default:
			canonical = url + "/" + rel
			canonical = strings.TrimSuffix(canonical, "index.html")
		}

		if err := tmpl.Execute(w, map[string]interface{}{
			"Title":     title,
			"Desc":      description,
			"Content":   content,
			"Canonical": canonical,
			"Nav":       nav,
			"Feedback":  feedback,
		}); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
