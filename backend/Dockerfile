FROM golang:1.27.1 

WORKDIR /app/

COPY . .

RUN go build -o app ./cmd

CMD ["./app"] 

