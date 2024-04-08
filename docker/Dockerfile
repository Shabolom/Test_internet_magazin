FROM golang:alpine
LABEL author=Timur
LABEL tipe=demo

WORKDIR /app
COPY cmd/main.go .

RUN go build -o main main.go
CMD [". /main"]
