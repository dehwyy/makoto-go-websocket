FROM golang

WORKDIR /src

COPY . .

RUN cd cmd && go build -o main .

CMD ["./cmd/main"]
