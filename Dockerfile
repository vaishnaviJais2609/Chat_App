# ---- build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/server .

# ---- runtime stage ----
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /out/server /app/server
COPY static/ /app/static/

ENV GIN_MODE=release
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/server"]