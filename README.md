# Ecommerce API

API REST de e-commerce escrita em Go com `gin`, autenticação por JWT e armazenamento em memória para fins de estudo e evolução incremental.

## Visão geral

O projeto está organizado por domínio em `internal/`:

- `user`: cadastro, login e leitura do usuário autenticado.
- `product`: criação e listagem de produtos.
- `cart`: gerenciamento de carrinho do usuário autenticado.
- `order`: checkout e marcação de pedido como pago.
- `payment`: integração com Stripe para criar `PaymentIntent` e receber webhook.
- `auth`: geração de token JWT e middlewares de autenticação/autorização.

Hoje a aplicação usa repositórios em memória, então os dados são perdidos ao reiniciar a API.

## Stack

- Go `1.26.1`
- Gin
- JWT (`github.com/golang-jwt/jwt/v5`)
- Stripe (`github.com/stripe/stripe-go/v85`)
- UUID (`github.com/google/uuid`)
- Bcrypt (`golang.org/x/crypto/bcrypt`)

## Como rodar

1. Configure a variável `STRIPE_SECRET_KEY` no `.env` ou no shell.
2. Inicie a API:

```bash
go run ./cmd/api
```

A aplicação sobe em `http://localhost:8080`.

## Como testar

Execute toda a suíte:

```bash
go test ./...
```

Cobertura adicionada neste repositório:

- Testes unitários para `internal/user/service.go`
- Testes unitários para `internal/order/service.go`
- Testes de integração HTTP em `cmd/api/main_integration_test.go`

## Endpoints

### Públicos

- `POST /register`: cria um usuário.
- `POST /login`: autentica e retorna um token JWT.
- `GET /products`: lista produtos.
- `POST /webhook`: endpoint de webhook do Stripe.

### Autenticados

Envie `Authorization: Bearer <token>`.

- `GET /me`: retorna `user_id` e `role`.
- `POST /cart/items`: adiciona item ao carrinho.
- `DELETE /cart/items/:product_id`: remove item do carrinho.
- `PUT /cart/items/:product_id`: atualiza a quantidade de um item.
- `GET /cart`: retorna o carrinho atual.
- `POST /checkout`: cria pedido pendente e Payment Intent.
- `POST /orders/:order_id/pay`: marca um pedido como pago.

### Administrador

Requer token de usuário com `role=admin`.

- `POST /products`: cria um produto.

## Fluxo resumido

1. O cliente se registra em `POST /register`.
2. Faz login em `POST /login` e recebe um JWT.
3. Com o token, gerencia o carrinho.
4. Em `POST /checkout`, a API valida estoque, cria o pedido e inicia o pagamento.
5. Após confirmação do pagamento, o pedido pode ser marcado como pago e o estoque é reduzido.

## Limitações atuais

- Persistência apenas em memória.
- `secret` JWT fixo no código.
- Segredo do webhook Stripe ainda está hardcoded.
- Sem migrations, banco de dados ou observabilidade.
- Alguns fluxos ainda dependem de melhorias de validação e tratamento de erro.
