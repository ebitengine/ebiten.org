# ebiten.org

## Generating

Edit HTML files under `contents` and run:

```sh
go generate .
```

## Contributions

Contributions are welcome. However, if you try to contribute one or more sections or articles, please ask Hajime Hoshi <hajimehoshi@gmail.com> before writing.

## Uploading WASM files

```
GOOGLE_APPLICATION_CREDENTIALS=/path/to/credentials.json go run ./cmd/uploadwasm/ -ebitenpath=../path/to/ebiten -upload
```

## Test on your local machine

Run the following command:
```
go run server.go ./docs
```

Navigate to [http://localhost:8000/](http://localhost:8000/) to view the changes
