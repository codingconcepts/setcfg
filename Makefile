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
	GOOS=linux go build -o setcfg_linux setcfg.go
	GOOS=windows go build -o setcfg_windows setcfg.go

##################
# Docker targets #
##################

docker_build:
	GOOS=linux go build -o setcfg setcfg.go
	docker build -t codingconcepts/setcfg:latest .
	docker push codingconcepts/setcfg:latest

docker_run:
	docker run --rm -it \
		-v c:/dev/github.com/codingconcepts/setcfg/examples:/examples \
		codingconcepts/setcfg:latest \
			-i examples/input.yaml -e examples/dev.yaml > blah.yaml
