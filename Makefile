all: journal

journal: $(shell find . -name "*.go") bindata.go
	go build -ldflags="-s -w" -o ./journal

bindata.go: static/bundle.js static/index.html static/global.css static/bundle.css
	go-bindata -o bindata.go static/...

static/bundle.js: $(shell find client)
	./node_modules/.bin/rollup -c
