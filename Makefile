build:
	go build -o cli
test:
	go test ./... -timeout=2m -parallel=4
