run:
	go run main.go

test:
	go test ./... -coverprofile=coverage.out -timeout=2m -parallel=4

coverage: test
	go tool cover -html=coverage.out
