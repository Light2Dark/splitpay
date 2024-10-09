setup-mac:
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-macos-arm64
	chmod +x tailwindcss-macos-arm64
	mv tailwindcss-macos-arm64 static/css/tailwindcss

	curl -sL https://unpkg.com/htmx.org/dist/htmx.min.js > ./static/htmx.min.js
	go install github.com/bokwoon95/wgo@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go mod download

setup:
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
	chmod +x tailwindcss-linux-x64
	mv tailwindcss-linux-x64 static/css/tailwindcss

	curl -sL https://unpkg.com/htmx.org/dist/htmx.min.js > ./static/htmx.min.js

	go install github.com/bokwoon95/wgo@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go mod download

docker-setup:
	curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.11/tailwindcss-linux-arm64
	chmod +x tailwindcss-linux-arm64
	mv tailwindcss-linux-arm64 static/css/tailwindcss

	go install github.com/bokwoon95/wgo@latest
	go install github.com/a-h/templ/cmd/templ@latest

dev:
	make -j 3 run_tailwind run_templ run_go port?=8080

run_tailwind:
	./static/css/tailwindcss -i static/css/input.css -o static/css/output.css --watch

run_go:
	wgo run ./cmd/api -port=$(port)

run_templ:
	templ generate --watch --proxy="http://localhost:$(port)"

build:
	go mod tidy && \
   	templ generate && \
	./static/css/tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify && \
	CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o bin/go-starter ./cmd/api && \
	docker build -t splitpay .

update_packages:
	go get -d -u ./...
	go mod tidy