run:
	build

install:
	@go get ./...
	@go mod tidy
	@go mod download

build: install
	env GOOS=darwin GOARCH=arm64 go build -o bin/mac/s3p .
	env GOOS=darwin GOARCH=amd64 go build -o bin/mac/s3p-intel .
	env GOOS=linux GOARCH=amd64 go build -o bin/linux/s3p .

test: install build
	@./bin/mac/s3p --profile "profiles/test_profile_aws.yaml"