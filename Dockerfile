FROM --platform=linux/amd64 golang:1.21.1
WORKDIR /app
COPY . .
RUN go build -o main .
RUN ls -l
CMD ["./main"]