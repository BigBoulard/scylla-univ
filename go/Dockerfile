# Build the manager binary
FROM golang:1.21 as builder

COPY . /build
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /build/mutants_01 ./cmd/mutants_01
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /build/mutants_02 ./cmd/mutants_02
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /build/mutants_03 ./cmd/mutants_03

FROM alpine:3.18.4

COPY --from=builder /build/mutants_01 /usr/local/bin/mutants_01
COPY --from=builder /build/mutants_02 /usr/local/bin/mutants_02
COPY --from=builder /build/mutants_03 /usr/local/bin/mutants_03
RUN chmod +x /usr/local/bin/mutants_01
RUN chmod +x /usr/local/bin/mutants_02
RUN chmod +x /usr/local/bin/mutants_03

# images for mutants_02
RUN mkdir -p /usr/share/icons/mutants
COPY --from=builder /build/cmd/mutants_02/img/*.jpg /usr/share/icons/mutants/

# a simple way to create a process that runs in the foreground.
# /dev/null is a special file in Unix-like operating systems that discards all data written to it.
# This means that the tail -f /dev/null command will simply read from the file and discard all of the data.
CMD tail -f /dev/null
