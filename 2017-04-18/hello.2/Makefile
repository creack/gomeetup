# Makefile
APP_DIR         = hello
BUILD_IMG       = hello:build
RELEASE_IMG     = hello:release

build           : .asset.tar
.asset.tar      : Dockerfile hello.go
		docker build --build-arg APP_DIR=$(APP_DIR) -t $(BUILD_IMG) .
		docker run --rm $(BUILD_IMG) tar cf - /app > $@

release         : .release
.release        : .asset.tar Dockerfile.release
		docker build -t $(RELEASE_IMG) -f Dockerfile.release .
		@touch $@

clean           :
		@rm -f .asset.tar .release

.PHONY          : build release
