# Dockerfile
FROM            golang:1.8

# Add our project.
ADD             . $APP_PATH

# Set entrypoint.
ENTRYPOINT      ["/app"]

# Compile.
RUN             go build -o /app
