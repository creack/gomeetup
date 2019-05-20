# Makefile
APP_DIR         = fetch
BUILD_IMG       = fetch:build
RELEASE_IMG     = fetch:release
SRCS            = fetch.go
EXTRA_ASSETS    = /etc/ssl

build           : .assets.tar
.assets.tar     : Dockerfile $(SRCS)
		docker build --build-arg APP_DIR=$(APP_DIR) -t $(BUILD_IMG) .
		docker run --rm $(BUILD_IMG) tar cf - /app $(EXTRA_ASSETS) > $@

release         : .release
.release        : .assets.tar Dockerfile.release
		docker build -t $(RELEASE_IMG) -f Dockerfile.release .
		@touch $@

clean           :
		@rm -f .assets.tar .release

.PHONY          : build release
