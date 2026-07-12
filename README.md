# ebiten.org

This repository manages the old website ebiten.org. It has moved to
[ebitengine.org](https://ebitengine.org/); every path is redirected there.

The redirects are defined in [`dist/_redirects`](dist/_redirects) and deployed
to [Cloudflare Pages](https://developers.cloudflare.com/pages/), which serves
them as server-side 301 redirects. To change where a path redirects, edit
`dist/_redirects` directly.
