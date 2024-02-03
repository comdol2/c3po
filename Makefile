.XPORT_ALL_VARIABLES:

# Set it to 1 to indicate make is running via the Jenkinsfile.
VERSION=default

# Go tests are performed in a Docker container.
SUFFIX=x86_64.zip
GITHUB=github.com
ORG=comdol2
REPO=snow
GOWS=/home/jenkins/agent
GOENV_flags = -e GO111MODULE=on -e GOPRIVATE=github.com -e GOPATH=$(GOWS) -e GOBIN=$(GOWS)/bin -e CGO_ENABLED=0
GOENV_flags_jenkins = GO111MODULE=on GOPRIVATE=github.com GOPATH=$(GOWS) GOBIN=$(GOWS)/bin CGO_ENABLED=0
OUTPUT = ./bin

# https://blog.golang.org/cover
cmd_test := go test -v -covermode=count -coverprofile=count.out -cover ./api
cmd_cover_func := go tool cover -func=count.out
cmd_go_build := go build -ldflags '-s -w -X $(GITHUB)/$(ORG)/$(REPO)/cmd.version=$(VERSION)' -a -o $(OUTPUT)/snow snow.go
cmd_go_exebuild := go build -ldflags '-s -w -X $(GITHUB)/$(ORG)/$(REPO)/cmd.version=$(VERSION)' -a -o $(OUTPUT)/snow.exe snow.go
cmd_go_build_linux = $(cmd_go_build) && zip -j $(OUTPUT)/snow_Linux_x86_64.zip $(OUTPUT)/snow
cmd_go_build_darwin = $(cmd_go_build) && zip -j $(OUTPUT)/snow_Darwin_x86_64.zip $(OUTPUT)/snow
cmd_go_build_windows = $(cmd_go_exebuild) && zip -j $(OUTPUT)/snow_Windows_amd64.zip $(OUTPUT)/snow.exe
cmd_vendor = go mod vendor
cmd_update = go mod tidy -v
cmd_create_release := $(OUTPUT)/github release create -r $(ORG)/$(REPO) -t $(VERSION) -n $(VERSION) -a $(GOWS)/.github_token
cmd_attach_linux_to_release := $(OUTPUT)/github release upload -r $(ORG)/$(REPO) -t $(VERSION) -f $(OUTPUT)/snow_Linux_x86_64.zip -c application/zip -a $(GOWS)/.github_token
cmd_attach_darwin_to_release := $(OUTPUT)/github release upload -r $(ORG)/$(REPO) -t $(VERSION) -f $(OUTPUT)/snow_Darwin_x86_64.zip -c application/zip -a $(GOWS)/.github_token
cmd_attach_darwin_to_release := $(OUTPUT)/github release upload -r $(ORG)/$(REPO) -t $(VERSION) -f $(OUTPUT)/snow_Windows_amd64.zip -c application/zip -a $(GOWS)/.github_token


.PHONY: all info tests build vendor release


all: info
info:
	@echo "This makefile encompasses the development workflow.\nTypical targets are: \n     tests, ws, clean."
	@echo "\nExamples:"
	@echo "    make tests                 - Run go tests."
	@echo "    make ws                    - Run go on workstation."
	@echo "    make clean                 - Remove vendor and glide.lock"


tests: unit-tests int-tests
unit-tests:
	@echo "  >  Running unit tests..."
	# Jenkins will already spin a k8s pod or docker container before running below commands.
	$(GOENV_flags_jenkins) $(cmd_test)
	$(GOENV_flags_jenkins) $(cmd_cover_func)
int-tests:
	@echo "  >  Running integration tests..."
	# Jenkins will already spin a k8s pod or docker container before running below commands.
	$(GOENV_flags_jenkins) $(cmd_test)
	$(GOENV_flags_jenkins) $(cmd_cover_func)

update:
	@echo "  >  Update Go vendors using go modules..."
	# Jenkins will already spin a k8s pod or docker container before running below commands.
	$(GOENV_flags_jenkins) $(cmd_update)

vendor: update
	@echo "  >  Downloading Go vendors into vendor directory using go modules..."
	# Jenkins will already spin a k8s pod or docker container before running below commands.
	$(GOENV_flags_jenkins) $(cmd_vendor)

build: clean vendor build_windows build_linux build_darwin
build_docker:
    docker build . -t 
build_linux:
	@echo "  >  Creating GitHub binary for Linux x86_64..."
	GOOS=linux $(GOENV_flags_jenkins) $(cmd_go_build_linux)
build_darwin:
	@echo "  >  Creating GitHub binary for Darwin x86_64..."
	GOOS=darwin $(GOENV_flags_jenkins) $(cmd_go_build_darwin)
build_windows:
	@echo "  >  Creating GitHub binary for Windows amd64..."
	GOOS=windows $(GOENV_flags_jenkins) $(cmd_go_build_windows)

# For workstation only
ws: clean
	$(cmd_godep) && $(cmd_test) && $(cmd_cover_func)

clean:
	rm -rf vendor bin
	rm -f count.out _github_permissions
