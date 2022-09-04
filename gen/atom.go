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
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/tools/blog/atom"
)

func toAtomTimeStr(ymd string) atom.TimeStr {
	// TODO: Enable to change the timezone.
	return atom.TimeStr(ymd + "T00:00:00+09:00")
}

func writeAtom(url string) error {
	author := &atom.Person{
		Name:  "Hajime Hoshi",
		URI:   "https://hajimehoshi.com",
		Email: "hajimehoshi@gmail.com",
	}

	feed := &atom.Feed{
		Title: "Ebiten Blog",
		ID:    url + "/blog/",
		Link: []atom.Link{
			{
				Rel:  "self",
				Href: url + "/blog/feed.xml",
			},
		},
		Author: author,
	}

	fs, err := ioutil.ReadDir(filepath.Join("contents", "blog"))
	if err != nil {
		return err
	}

	var pages []*page
	for _, f := range fs {
		if filepath.Ext(f.Name()) != ".html" {
			continue
		}
		if f.Name() == "index.html" {
			continue
		}
		if f.Name() == "nav.html" {
			continue
		}
		path := filepath.Join("contents", "blog", f.Name())
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		p, err := newPage(b, path)
		if err != nil {
			return err
		}
		pages = append(pages, p)
	}
	sort.Slice(pages, func(a, b int) bool {
		ac, _ := pages[a].created()
		bc, _ := pages[b].created()
		return ac > bc
	})

	u, err := pages[0].created()
	if err != nil {
		return err
	}
	feed.Updated = toAtomTimeStr(u)

	for _, p := range pages {
		t, err := p.title()
		if err != nil {
			return err
		}
		created, err := p.created()
		if err != nil {
			return err
		}

		feed.Entry = append(feed.Entry, &atom.Entry{
			Title: t,
			ID:    url + "/blog/" + p.name(),
			Link: []atom.Link{
				{
					Rel:  "alternate",
					Href: url + "/blog/" + p.name(),
				},
			},
			Published: toAtomTimeStr(created),
			Updated:   toAtomTimeStr(created),
			Author:    author,
			Content: &atom.Text{
				Type: "html",
				Body: p.content,
			},
		})
	}

	f, err := os.OpenFile(filepath.Join(outDir, "blog", "feed.xml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	e := xml.NewEncoder(f)
	defer e.Flush()
	e.Indent("", "  ")
	if err := e.Encode(feed); err != nil {
		return err
	}

	return nil
}
