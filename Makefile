all:
	go build -o parser.exe ./cmd/cli
	./parser -f="file.txt" -t=2s

#generate mocks:
generate:
	mockgen -source ./internal/parser.go -destination ./internal/mock/parser.go

#run all tests:
test:
	go test -v -race ./...

#test coverage
test_coverage:
	go test -short -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	del coverage.out
