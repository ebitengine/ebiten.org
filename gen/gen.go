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

package gen

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

func walkHTML(node *html.Node, f func(node *html.Node) error) error {
	if err := f(node); err != nil {
		return err
	}
	for n := node.FirstChild; n != nil; n = n.NextSibling {
		if err := walkHTML(n, f); err != nil {
			return err
		}
	}
	return nil
}

func findFirstElementByName(node *html.Node, name string) (*html.Node, error) {
	t := errors.New("regular termination")
	var found *html.Node
	if err := walkHTML(node, func(node *html.Node) error {
		if node.Type == html.ElementNode && node.Data == name {
			found = node
			return t
		}
		return nil
	}); err != nil && err != t {
		return nil, err
	}
	return found, nil
}

func findElementByID(node *html.Node, id string) (*html.Node, error) {
	t := errors.New("regular termination")
	var found *html.Node
	if err := walkHTML(node, func(node *html.Node) error {
		if node.Type == html.ElementNode {
			for _, attr := range node.Attr {
				if attr.Key != "id" {
					continue
				}
				if attr.Val != id {
					break
				}
				found = node
				return t
			}
		}
		return nil
	}); err != nil && err != t {
		return nil, err
	}
	return found, nil
}

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

func Run(url, description string) error {
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
			if filepath.Base(path) == "nav.html" {
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

		b := bytes.NewReader([]byte(content))
		node, err := html.Parse(b)
		if err != nil {
			return err
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
			h1, err := findFirstElementByName(node, "h1")
			if err != nil {
				return err
			}
			if h1 != nil {
				title = fmt.Sprintf("%s - Ebiten", h1.FirstChild.Data)
			}
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
			// When generated on a Windows machine, rel will have \ characters.
			// Use ToSlash to ensure that all path separators are /.
			canonical = url + "/" + filepath.ToSlash(rel)
			canonical = strings.TrimSuffix(canonical, "index.html")
		}

		f := filepath.Join(filepath.Dir(path), "nav.html")
		c, err := ioutil.ReadFile(f)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
		subnav := string(c)

		var meta map[string]interface{}
		n, err := findElementByID(node, "meta")
		if err != nil {
			return err
		}
		if n != nil {
			if err := json.Unmarshal([]byte(n.FirstChild.Data), &meta); err != nil {
				return err
			}
		}

		share := "https://ebiten.org/images/share.png"
		if meta != nil {
			if s, ok := meta["Share"]; ok {
				share = url + s.(string)
			}
		}

		if err := tmpl.Execute(w, map[string]interface{}{
			"Title":     title,
			"Desc":      description,
			"Content":   content,
			"Share":     share,
			"Canonical": canonical,
			"NavExists": nav,
			"SubNav":    subnav,
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
