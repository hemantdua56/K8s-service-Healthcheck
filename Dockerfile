FROM golang:1.14.1 

ENV GO111MODULE=on 

RUN  mkdir /app 
ADD . /app/ 
WORKDIR /app
RUN go mod init main.go && go get k8s.io/client-go@v0.17.0
RUN go build -o main . 
CMD ["/app/main"]
