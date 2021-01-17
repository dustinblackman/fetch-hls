.DEFAULT_GOAL := build

build:
	go build -o fetch-hls -ldflags="-X main.version=$$(cat VERSION)"

install:
	go install -ldflags="-X main.version=$$(cat VERSION)"

lint:
	gomodrun golangci-lint run

lint-fix:
	gomodrun golangci-lint run --fix

release: lint
	go mod tidy
	git add .
	git commit -m "v$$(cat VERSION)"
	git tag -a "v$$(cat VERSION)" -m "v$$(cat VERSION)"
	git push
	git push --tags
	gomodrun goreleaser release --rm-dist

release-snapshot:
	gomodrun goreleaser release --snapshot --skip-publish --rm-dist
