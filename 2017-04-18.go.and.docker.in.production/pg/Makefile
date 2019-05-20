# Makefile
APP_DIR         = pg
BUILD_IMG       = pg:build
RELEASE_IMG     = pg:release
SRCS            = pg.go pg_test.go
EXTRA_ASSETS    = /etc/ssl /usr/local/go/lib/time/zoneinfo.zip
PG_IMG          = postgres:9.6
PG_C            = pg_c

build           : .assets.tar
.assets.tar     : Dockerfile $(SRCS)
		docker build --build-arg APP_DIR=$(APP_DIR) -t $(BUILD_IMG) .
		docker run --rm $(BUILD_IMG) tar cf - /app $(EXTRA_ASSETS) > $@

release         : .release
.release        : .assets.tar Dockerfile.release
		docker build -t $(RELEASE_IMG) -f Dockerfile.release .
		@touch $@

start           : .pg_start .release
		docker run --rm -it --link "$(PG_C):pg" $(RELEASE_IMG) -pg "postgres://postgres@pg/postgres?sslmode=disable"

pg_start        : .pg_start
.pg_start       :
		@docker rm -f -v $(PG_C) >& /dev/null || true
		docker run -d --name $(PG_C) $(PG_IMG)
		@while ! docker exec $(PG_C) psql -U postgres -d postgres -c "SELECT;" >& /dev/null; do sleep 1; done
		@echo "Testing Postgres DB up and running."
		@touch $@

test            : .pg_start .assets.tar
		docker run --rm -it --link "$(PG_C):pg" $(BUILD_IMG) go test -v -cover -pg "postgres://postgres@pg/postgres?sslmode=disable"

clean           :
		@docker rm -f -v $(PG_C) >& /dev/null || true
		@rm -f .assets.tar .release .pg_start

.PHONY          : build release test pg_start start
