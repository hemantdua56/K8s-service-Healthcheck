FROM golang:1.14.1
ENV GO111MODULE=on 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app
RUN go build -o main ./src/
CMD ["/app/main"]
