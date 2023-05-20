.PHONY: cover
cover:
	go test -short -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: gen
gen:
	mockgen -source=internal/usecase/interfaces.go -destination=internal/usecase/mocks/mock_interfaces.go

.PHONY: total
total:
	go test -coverpkg=./... -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -func=coverage.out
	rm coverage.out

.PHONY: run
run:
	go run cmd/shortener/main.go -d="postgres://postgres:55555@127.0.0.1:5432/dburl" -t="127.0.0.1/8"