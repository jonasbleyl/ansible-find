test:
	go test ./... -count=1 -tags=test

coverage:
	go test ./... -count=1 -tags=test -coverprofile=coverage.out

report: coverage
	go tool cover -html=coverage.out

build:
	go build -o ansible-vars cmd/ansible-vars/main.go