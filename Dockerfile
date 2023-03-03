FROM golang:1.19-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /k8s-chaos-monkey

CMD ["/k8s-chaos-monkey"]