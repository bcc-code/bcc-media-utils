build: dist/watcher-linux-amd64 dist/watcher-linux-arm64 dist/watcher-darwin-amd64 dist/watcher-windows-amd64.exe

dist/watcher-linux-amd64: *.go
	GOOS=linux GOARCH=amd64 go build -o dist/watcher-linux-amd64 .

dist/watcher-linux-arm64: *.go
	GOOS=linux GOARCH=arm64 go build -o dist/watcher-linux-arm64 .

dist/watcher-darwin-amd64: *.go
	GOOS=darwin GOARCH=amd64 go build -o dist/watcher-darwin-amd64 .

dist/watcher-windows-amd64.exe: *.go
	GOOS=windows GOARCH=amd64 go build -o dist/watcher-windows-amd64.exe .

