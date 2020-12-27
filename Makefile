build: 
	go build

all: generate build

generate: static/keyevent.js
	go-bindata -o static/static.go -pkg static static/
