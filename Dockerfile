FROM golang:1.14-alpine AS build
WORKDIR /go/src/proglog
COPY . .
RUN CGO_ENABLE=0 go build -o /go/bin/proglog ./cmd/proglog
FROM scratch
COPY --from=build /go/bin/proglog /bin/proglog
ENTRYPOINT ["/bin/proglog"]