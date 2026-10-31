# Etapa 1 (build): imagem grande, com o compilador Go.
FROM golang:1.27-alpine AS build
WORKDIR /src

# Copiar o go.mod antes do código aproveita o cache de camadas do Docker:
# se só o código mudar, as dependências não são baixadas de novo.
COPY go.mod ./
RUN go mod download

COPY *.go ./
# CGO_ENABLED=0 gera um binário estático (não depende de bibliotecas C),
# que roda numa imagem sem sistema operacional completo.
# -ldflags="-s -w" remove símbolos de depuração e diminui o binário.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api .

# Etapa 2 (runtime): só o binário, sem compilador nem shell.
# A distroless "nonroot" já roda como usuário sem privilégios de root.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
