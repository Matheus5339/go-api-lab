# go-api-lab

API HTTP pequena em Go, feita para aprender e praticar back-end em Go:
roteamento, JSON, testes, contêiner e CI. Usa **só a biblioteca padrão**
(`net/http`, `encoding/json`), sem framework.

## Endpoints

| Método e rota      | Resposta                                                   |
|--------------------|------------------------------------------------------------|
| `GET /health`      | `{"status":"ok","uptime_seconds":12.3}`                     |
| `GET /metrics`     | `{"requests_total":42}` (conta todas as requisições, inclusive esta) |
| `GET /items/{id}`  | `{"id":1,"name":"Teclado mecânico","price_cents":29990}`    |

Erros em `/items/{id}`:

- `400` se o id não for um inteiro positivo (`/items/abc`, `/items/0`)
- `404` se o item não existir (`/items/999`)

Qualquer outro método nas rotas (ex.: `POST /health`) recebe `405`.

## Como rodar

Requer Go 1.22+ (usa os padrões de rota com método e `{id}`).

```bash
go run .
curl http://localhost:8080/health
```

A porta vem da variável `PORT` (padrão `8080`).

### Com Docker

```bash
docker build -t go-api-lab .
docker run --rm -p 8080:8080 go-api-lab
```

O Dockerfile é multi-stage: compila numa imagem com Go e copia só o
binário estático para uma imagem `distroless` sem shell, rodando como
usuário não-root.

## Testes

```bash
go test -v -cover ./...
```

Os testes usam `net/http/httptest` para chamar os handlers sem abrir porta
de rede. Cobrem os três endpoints, os erros 400/404 e o 405.

O CI (GitHub Actions) roda `gofmt`, `go vet`, `go test -race` e o build da
imagem Docker a cada push na `main`.

## Decisões

- **Biblioteca padrão em vez de framework** (Gin, Echo): desde o Go 1.22 o
  `http.ServeMux` já faz roteamento por método e parâmetro de caminho.
- **Contador com `atomic.Int64`**: cada requisição roda numa goroutine própria;
  um `int` com `++` teria condição de corrida.
- **Preço em centavos (`int`)**: `float` tem erro de arredondamento com dinheiro.
- **Graceful shutdown**: ao receber SIGTERM (o que o Cloud Run faz antes de
  parar o contêiner), o servidor termina as requisições em andamento.

## O que ficou de fora

- Banco de dados: os itens são fixos, em memória.
- Criar/editar/apagar itens (só leitura).
- Autenticação e limite de requisições (rate limiting).
- Métricas no formato Prometheus (o `/metrics` aqui é um JSON simples).
- Deploy em nuvem: planejado para o Google Cloud Run, ainda não feito.
