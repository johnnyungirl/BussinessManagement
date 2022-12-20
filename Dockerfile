FROM golang:alpine
LABEL maintainer="DEN"
RUN apk update && apk add --no-cache git && apk add --no-cache bash
RUN mkdir /app
WORKDIR /app
COPY . . 
COPY .env .
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o /build
CMD [ "/build" ]