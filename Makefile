GIT_COMMIT=$(shell git describe --always)

all: build
default: build

build:
	go build
	mv chatgpt-telegram whisper-telegram
clean:
	rm chatgpt-telegram
