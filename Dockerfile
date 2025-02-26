FROM golang:1.24-bullseye AS builder

WORKDIR /app/

COPY . /app/

RUN apt-get update && apt-get install -y curl && \
    curl -fsSL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
    -o /usr/local/bin/tailwindcss && chmod +x /usr/local/bin/tailwindcss

RUN go install github.com/a-h/templ/cmd/templ@latest

COPY . .

RUN make build-webapp

CMD ["./build/c1"]

FROM golang:alpine3.21

WORKDIR /app/

COPY --from=builder /app/build /app

EXPOSE 7777

CMD ["/app/c1"]
