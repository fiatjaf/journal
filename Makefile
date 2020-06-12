all: journal

dev: $(shell find . -name "*.go")
	go-bindata -debug -o bindata.go static/...
	go build -ldflags="-s -w" -o ./journal

journal: $(shell find . -name "*.go") bindata.go
	go build -ldflags="-s -w" -o ./journal

server/bindata.go: static/bundle.js static/index.html static/global.css static/bundle.css
	go-bindata -o server/bindata.go static/...

static/bundle.js: $(shell find client)
	./node_modules/.bin/rollup -c
