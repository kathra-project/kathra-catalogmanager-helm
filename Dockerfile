FROM golang:1-alpine as builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine/helm:2.15.2 
RUN apk --no-cache add ca-certificates bash sed grep gawk
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"] 