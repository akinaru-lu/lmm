.PHONY: app api image

all:
	make run -j

install:
	cd api && dep ensure
	cd image && dep ensure
	npm --prefix app install

run: app api image

app: app/package.json
	npm --prefix app run dev

api: api/main.go
	go run api/main.go

image: image/main.go
	go run image/main.go

