# First stage: build the executable.
FROM golang:alpine AS builder

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk add --no-cache ca-certificates git sqlite gcc musl-dev
# Install packr for building the binary
RUN go get -u github.com/gobuffalo/packr/v2/packr2

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Import the code from the context.
COPY ./ ./
RUN mkdir config && chown nobody:nobody config && cp config-template.json config/config.json && touch wtfd-db.log

# Build the executable to `/app`. Mark the build as statically linked.
RUN cd ./internal && $(go env GOPATH)/bin/packr2

RUN go build -ldflags="-linkmode external -extldflags -static -s -w" \
    ./cmd/wtfd.go
RUN ls -al . internal

# Final stage: the running container.
FROM scratch AS final
WORKDIR /

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the first stage.
COPY --from=builder /app/wtfd /app/config /

# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8080

# Set env
ENV WTFD_CONFIG_FILE=/config/config.json
ENV WTFD_DB_FILE=/config/wtfd-db.db
ENV WTFD_DBLOG_FILE=/config/wtfd-db.log

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/wtfd"]
