FROM node:22-alpine as tailwind
WORKDIR /app
COPY ui ./ui
RUN npm install -D @tailwindcss/cli && \
    npx @tailwindcss/cli -i ./ui/static/src/input.css -o ./ui/static/src/output.css && \
    npm cache clean --force


FROM golang:1.26-alpine AS go-builder
WORKDIR /app


ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off


COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o web ./cmd/web


FROM alpine
WORKDIR /app

COPY ./ui ./ui

COPY --from=tailwind /app/ui/static/src/output.css ./ui/static/src/output.css
COPY --from=go-builder /app/web ./web

EXPOSE 8080
CMD ["./web"]