COVER_PROFILE=cover.out
COVER_HTML=cover.html

.PHONY: $(COVER_PROFILE) $(COVER_HTML)

all: open

build: clean
	mkdir -p -v ./bin/wpscheduler.app/Contents/Resources
	mkdir -p -v ./bin/wpscheduler.app/Contents/MacOS
	cp ./appInfo/*.plist ./bin/wpscheduler.app/Contents/Info.plist
	cp ./appInfo/*.icns ./bin/wpscheduler.app/Contents/Resources/icon.icns
	cp ./configs/config.yaml ./bin/wpscheduler.app/Contents/MacOS/config.yaml
	go build -o ./bin/wpscheduler.app/Contents/MacOS/wpscheduler cmd/main.go

open: build
	open ./bin

clean:
	rm -rf ./bin

start:
	go run cmd/main.go

test:coverage

coverage: $(COVER_HTML)

$(COVER_HTML): $(COVER_PROFILE)
	go tool cover -html=$(COVER_PROFILE) -o $(COVER_HTML)

$(COVER_PROFILE):
	go test -v -failfast -race -coverprofile=$(COVER_PROFILE) ./...

vet:
	go vet $(shell glide nv)

lint:
	go list ./... | grep -v vendor | grep -v /assets/ |xargs -L1 golint -set_exit_status

.PHONY: build
.PHONY: clean