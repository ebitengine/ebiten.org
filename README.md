# ebiten.org

This repository manages the old website ebiten.org. All the pages redirect to the new sites ebitengine.org.

## Generating

Edit HTML files under `contents` and run:

```sh
go run gen.go
```

## Test on your local machine

```
go run server.go ./_site
```

Validate your changes by opening http://127.0.0.1:8000.
