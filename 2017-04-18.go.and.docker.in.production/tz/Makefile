# Makefile
APP_DIR         = tz
BUILD_IMG       = tz:build
RELEASE_IMG     = tz:release
SRCS            = tz.go
EXTRA_ASSETS    = /etc/ssl /usr/local/go/lib/time/zoneinfo.zip

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
