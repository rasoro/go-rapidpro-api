test:
	go test -covermode=count ./...
test-cover:
	go test -coverprofile cover.out -timeout 120s ./... && go tool cover -html=cover.out