# go-api-lab

Projeto de estudo do Matheus para ter Go (e, se ele decidir, GCP) como fato
no currículo. Repositório alvo: `github.com/Matheus5339/go-api-lab` (público).

## Decisões
- Só biblioteca padrão (`net/http`, Go 1.22+ com padrões de rota). Sem framework.
- Porta pela variável `PORT` (exigência do Cloud Run).
- Estado compartilhado em `server`; contador com `atomic.Int64`.
- Testes com `httptest` em `server_test.go`; CI roda gofmt, vet, test -race e docker build.
- Fim de linha LF forçado por `.gitattributes` (o CI checa gofmt).

## Pendências
- Push para o GitHub: o Matheus faz (sem `gh` instalado; sem tokens no chat).
- Deploy no Cloud Run: exige conta GCP com faturamento (cartão). Decisão e
  cadastro são dele. Até lá, no README e no CV, GCP fica "em estudo".
- O código foi escrito pelo Claude a pedido dele (31/10/2026); ele vai estudar
  cada arquivo para conseguir defender em entrevista.
