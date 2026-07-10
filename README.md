# ebiten.org

This repository manages the old website ebiten.org. All the pages redirect to the new site ebitengine.org.

The site is a static set of redirect pages under `dist`, which is deployed to GitHub Pages as-is. Edit the HTML files under `dist` directly.

## Test on your local machine

Serve the `dist` directory with any static file server, e.g.:

```sh
python3 -m http.server 8000 --directory dist
```

Validate your changes by opening http://127.0.0.1:8000.
