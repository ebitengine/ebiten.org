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
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"cloud.google.com/go/storage"
	"golang.org/x/sync/errgroup"
)

var (
	flagEbitenPath = flag.String("ebitenpath", "", "path to ebiten repository")
)

func examples() ([]string, error) {
	f, err := os.Open(filepath.Join("contents", "examples"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	var es []string
	for _, n := range names {
		if !strings.HasSuffix(n, ".json") {
			continue
		}
		ext := filepath.Ext(n)
		es = append(es, n[:len(n)-len(ext)])
	}

	return es, nil
}

const bucket = "ebiten-dot-org-wasm"

func updateBucket(ctx context.Context) error {
	fmt.Println("Updating bucket...")

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	b := client.Bucket(bucket)
	if _, err := b.Update(ctx, storage.BucketAttrsToUpdate{
		CORS: []storage.CORS{
			{
				Methods: []string{"*"},
				Origins: []string{"*"},
			},
		},
	}); err != nil {
		return err
	}
	return nil
}

func uploadFile(ctx context.Context, name string, r io.Reader) error {
	fmt.Printf("Uploading %s...\n", name)

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	b := client.Bucket(bucket)
	if _, err := b.Update(ctx, storage.BucketAttrsToUpdate{
		CORS: []storage.CORS{
			{
				Methods: []string{"*"},
				Origins: []string{"*"},
			},
		},
	}); err != nil {
		return err
	}

	w := b.Object(name).NewWriter(ctx)
	defer w.Close()

	w.ACL = []storage.ACLRule{
		{
			Entity: storage.AllUsers,
			Role:   storage.RoleReader,
		},
	}

	if _, err := io.Copy(w, r); err != nil {
		return err
	}

	return nil
}

func run() error {
	es, err := examples()
	if err != nil {
		return err
	}

	tmpout, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}

	ctx := context.Background()

	if err := updateBucket(ctx); err != nil {
		return err
	}

	ch := make(chan string, 4)

	var g errgroup.Group
	g.Go(func() error {
		defer close(ch)

		var g errgroup.Group
		for _, e := range es {
			e := e
			g.Go(func() error {
				name := e + ".wasm"
				args := []string{
					"build",
					"-o", filepath.Join(tmpout, name),
					"-tags", "example",
					"./examples/" + e,
				}
				fmt.Println("go", strings.Join(args, " "))
				cmd := exec.Command("go", args...)
				cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
				cmd.Dir = *flagEbitenPath
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					return err
				}

				ch <- e + ".wasm"

				return nil
			})
		}
		return g.Wait()
	})
	g.Go(func() error {
		semaphore := make(chan struct{}, 4)

		var g errgroup.Group
		for name := range ch {
			name := name
			semaphore <- struct{}{}
			g.Go(func() error {
				defer func() {
					<-semaphore
				}()

				f, err := os.Open(filepath.Join(tmpout, name))
				if err != nil {
					return err
				}
				defer f.Close()

				if err := uploadFile(ctx, name, f); err != nil {
					return err
				}

				return nil
			})
		}
		return g.Wait()
	})

	if err := g.Wait(); err != nil {
		return err
	}

	if err := os.RemoveAll(tmpout); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()
	if *flagEbitenPath == "" {
		fmt.Fprintln(os.Stderr, "Specify -ebitenpath")
		os.Exit(1)
	}

	if err := run(); err != nil {
		panic(err)
	}
}
