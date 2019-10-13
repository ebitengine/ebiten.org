# ebiten.org

Edit HTML files under `contents` and run:

```sh
go generate .
```

## Upload WASM files

```
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json GO111MODULE=on go run uploadwasm.go -ebitenpath=../path/to/ebiten
```
