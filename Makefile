test:
	go test --race --cover $$(go list ./...)

fmt:
	gofmt -l -w ./

tidy:
	go mod tidy