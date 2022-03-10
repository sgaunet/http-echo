FROM golang:1.17.5-alpine AS builder
LABEL stage=builder

RUN apk add --no-cache git upx

WORKDIR /go/src/http-echo
ENV GOPATH /go

COPY . /go/src/http-echo
RUN echo $GOPATH
RUN go get 
RUN CGO_ENABLED=0 GOOS=linux go build .
RUN upx http-echo



FROM scratch AS final
WORKDIR /
COPY --from=builder /go/src/http-echo/http-echo .
COPY etc /etc
USER MyUser
EXPOSE 8080
CMD [ "/http-echo" ]
