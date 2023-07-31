build: dist/mediainfoserver-linux-amd64 dist/mediainfoserver-linux-arm64 dist/mediainfoserver-darwin-amd64 dist/mediainfoserver-windows-amd64.exe

dist/mediainfoserver-linux-amd64: *.go
	GOOS=linux GOARCH=amd64 go build -o dist/mediainfoserver-linux-amd64 ./cmd

dist/mediainfoserver-linux-arm64: *.go
	GOOS=linux GOARCH=arm64 go build -o dist/mediainfoserver-linux-arm64 ./cmd

dist/mediainfoserver-darwin-amd64: *.go
	GOOS=darwin GOARCH=amd64 go build -o dist/mediainfoserver-darwin-amd64 ./cmd

dist/mediainfoserver-windows-amd64.exe: *.go
	GOOS=windows GOARCH=amd64 go build -o dist/mediainfoserver-windows-amd64.exe ./cmd

