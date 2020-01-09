FROM golang:1-alpine as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o . cmd/kathra-catalogmanager-helm-server/main.go

FROM alpine/helm:3.0.0
RUN apk --no-cache add ca-certificates bash sed grep gawk
WORKDIR /root/
COPY repositories.yaml repositories.yaml
COPY --from=builder /app/main .
EXPOSE 8080
ENTRYPOINT [ "/bin/sh" ]
CMD ["-c", "./main"] 