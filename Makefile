PHONY: run templ sass cover

run: templ sass
	go run ./cmd/web

templ:
	templ generate

sass:
	sass ui/scss/styles.scss ui/static/styles.css

cover:
	go test ./... -covermode=count -coverprofile=/tmp/profile.out
	go tool cover -html=/tmp/profile.out
