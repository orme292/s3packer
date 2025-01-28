run:
	build

install:
	@go get ./...
	@go mod tidy
	@go mod download

build: install
	env GOOS=darwin GOARCH=arm64 go build -o bin/mac/s3p cmd/app/*.go
	env GOOS=darwin GOARCH=amd64 go build -o bin/mac-intel/s3p cmd/app/*.go
	env GOOS=linux GOARCH=amd64 go build -o bin/linux/s3p cmd/app/*.go

test: install build
	@./bin/mac/s3p --profile "profiles/test_profile_aws.yaml"