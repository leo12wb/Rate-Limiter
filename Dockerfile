FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server .

# use a distroless base image with glibc
FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=builder /app/server .
# run as non-privileged user
USER nonroot

ENTRYPOINT ["./server"]