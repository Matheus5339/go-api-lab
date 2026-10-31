# go-api-lab

Projeto de estudo do Matheus para ter Go e nuvem como fato no currículo.
Repositório: `github.com/Matheus5339/go-api-lab` (público).

## Decisões
- Só biblioteca padrão (`net/http`, Go 1.22+ com padrões de rota). Sem framework.
- Porta pela variável `PORT` (padrão 8080).
- Estado compartilhado em `server`; contador com `atomic.Int64`.
- Testes com `httptest` em `server_test.go`; CI roda gofmt, vet, test -race,
  docker build e publica a imagem em `ghcr.io/matheus5339/go-api-lab` (pública).
- Fim de linha LF forçado por `.gitattributes` (o CI checa gofmt).
- **Regra do Matheus: nada que custe dinheiro.** Por isso:
  - GCP abandonado: o teste gratuito com Pix exige pré-pagamento (não feito; nada cobrado).
  - Nuvem usada: **Azure for Students** (sem cartão/Pix; crédito de US$ 100 / 12 meses),
    assinatura "Azure for Students".
  - Imagem no GHCR (grátis) em vez do Azure Container Registry (cobra por dia).
  - Container App com min 0 / max 1 réplica, 0,25 vCPU, sem Log Analytics.
  - Nunca cadastrar cartão ou Pix no Azure.

## Deploy atual (31/10/2026)
- Grupo `rg-go-api-lab`, ambiente `env-go-api-lab`, app `go-api-lab`, região brazilsouth.
- URL: https://go-api-lab.calmwave-44962e56.brazilsouth.azurecontainerapps.io
- Atualizar: `az containerapp update ... --image ghcr.io/matheus5339/go-api-lab:<sha>`.
- Apagar tudo (zera qualquer consumo): `az group delete --name rg-go-api-lab`.

## Pendências
- O código foi escrito pelo Claude a pedido dele (31/10/2026); ele vai estudar
  cada arquivo para conseguir defender em entrevista.
- No CV: Go e Azure saem de "em estudo"; GCP continua "em estudo".
