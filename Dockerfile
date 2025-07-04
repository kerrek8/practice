FROM golang:1.24-alpine AS build
RUN apk add --no-cache alpine-sdk

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go
RUN CGO_ENABLED=1 GOOS=linux go build -o migrator cmd/migrator/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app
COPY --from=build /app/main /app/main
COPY --from=build /app/migrator /app/migrator
COPY --from=build /app/.env /app/.env
RUN mkdir /app/database
RUN touch /app/database/database.db
COPY --from=build /app/migrations /app/migrations
COPY --from=build /app/front /app/front
EXPOSE 8080
RUN ./migrator --migrations-path=./migrations --migrations-table=mig
CMD ["./main"]


