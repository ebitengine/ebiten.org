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
	"errors"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

type page struct {
	node    *html.Node
	path    string
	content string
}

func newPage(content []byte, path string) (*page, error) {
	b := bytes.NewReader([]byte(content))
	node, err := html.Parse(b)
	if err != nil {
		return nil, err
	}
	c := &page{
		node:    node,
		path:    path,
		content: strings.ReplaceAll(string(content), "\r\n", "\n"),
	}
	return c, nil
}

func (p *page) title() (string, error) {
	h1, err := findFirstElementByName(p.node, "h1")
	if err != nil {
		return "", err
	}
	return h1.FirstChild.Data, nil
}

func (p *page) redirect() (string, error) {
	r, err := findElementByID(p.node, "meta-redirect")
	if err != nil {
		return "", err
	}
	if r != nil {
		return r.FirstChild.Data, nil
	}
	rel, err := filepath.Rel("contents", p.path)
	if err != nil {
		return "", err
	}
	path := filepath.ToSlash(rel)
	if strings.HasSuffix(path, "/index.html") {
		path = path[:len(path)-len("index.html")]
	}
	return "https://ebitengine.org/en/" + path, nil
}

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
