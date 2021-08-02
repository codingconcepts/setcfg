################
# Test targets #
################

test:
	@go test ./... -v -cover

cover:
	go test --coverprofile=coverage.out
	go tool cover --html=coverage.out

###################
# Release targets #
###################

build:
	GOOS=darwin go build -o setcfg_darwin setcfg.go
	GOOS=darwin go build -o setcfg_linux setcfg.go
	GOOS=darwin go build -o setcfg_windows setcfg.go