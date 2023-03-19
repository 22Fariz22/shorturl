.PHONY: cover
cover:
	go test -short -count=1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

.PHONY: gen
gen:
	mockgen -source=internal/usecase/interfaces.go -destination=internal/usecase/mocks/mock_interfaces.go