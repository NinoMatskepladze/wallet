FROM golang:1.17-alpine3.14 

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./main"]