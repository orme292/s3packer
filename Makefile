install:
	@go get ./...
	@go mod tidy
	@go mod download

build: install
	env GOOS=darwin GOARCH=arm64 go build -o bin/mac/s3p cmd/s3p/*.go
	env GOOS=darwin GOARCH=amd64 go build -o bin/mac-intel/s3p cmd/s3p/*.go
	env GOOS=linux GOARCH=amd64 go build -o bin/linux/s3p cmd/s3p/*.go

run: build

test: install build
	@go test ./internal/conf
	@./bin/mac/s3p use -f "docs/test/test_profile_aws.yaml"