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
	"fmt"
	"html"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

var reTitle = regexp.MustCompile(`<h1>([^<]+)</h1>`)

func run() error {
	tmpl, err := template.New("template.html").Funcs(template.FuncMap{
		"noescape": func(str string) template.HTML {
			return template.HTML(str)
		},
	}).ParseFiles("template.html")
	if err != nil {
		return err
	}

	if err := filepath.Walk("contents", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".html" {
			return nil
		}
		c, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel("contents", path)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join("docs", filepath.Dir(rel)), 0755); err != nil {
			return err
		}
		w, err := os.Create(filepath.Join("docs", rel))
		if err != nil {
			return err
		}
		defer w.Close()

		title := "Ebiten - A dead simple 2D game library in Go"
		if path != filepath.Join("contents", "index.html") {
			m := reTitle.FindStringSubmatch(string(c))
			title = fmt.Sprintf("%s - Ebiten", html.UnescapeString(m[1]))
		}

		if err := tmpl.Execute(w, map[string]interface{}{
			"Title":   title,
			"Content": string(c),
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
