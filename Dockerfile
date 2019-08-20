FROM golang:1.12.9-alpine as build-env

RUN mkdir /build 

ADD main.go /build/

WORKDIR /build 

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

FROM scratch

COPY --from=build-env /build/main /app/

WORKDIR /app

CMD ["./main"]
