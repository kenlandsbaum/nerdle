cli:
	go run ./cmd/nerdle-cli/...
api:
	go run ./cmd/nerdle-api/...
test:
	go test ./... -coverprofile cover.out
cov:
	go tool cover -func cover.out
race:
	go test ./... -race

covhtml:
	go tool cover -html=cover.out