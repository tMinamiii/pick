ifdef update
	u=-u
endif

export GO111MODULE=on

.PHONY: deps
deps:
	go get ${u}	 -d
	go mod tidy

.PHONY: build
build: deps
	[ -d "build" ] || mkdir -p build
	go build -o pick main.go
