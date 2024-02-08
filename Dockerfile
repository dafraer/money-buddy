FROM golang:1.21

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o MoneyBuddy ./

ENV PORT 8000

EXPOSE 8000

CMD ["./MoneyBuddy"]