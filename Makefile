.PHONY: cover

cover:
	go test -coverprofile=coverage.tmp.out ./...
	# remove generated code from coverage 
	cat coverage.tmp.out | grep -v "_gen.go" > coverage.out
	go tool cover -func=coverage.out
