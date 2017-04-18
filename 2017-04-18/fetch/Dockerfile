# Dockerfile
FROM            golang:1.8

# Set workdir.
ARG             APP_DIR
ENV             APP_PATH $GOPATH/src/$APP_DIR
WORKDIR         $APP_PATH

# Add our project.
ADD             . $APP_PATH

# Fetch dependencies.
RUN             go get -d ./...

# Compile static binary.
# NOTE: We could use ARG for the binary path too.
RUN             CGO_ENABLED=0 go build -ldflags -d -o /app
