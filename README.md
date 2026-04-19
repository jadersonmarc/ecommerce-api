# Ecommerce API

solution of https://roadmap.sh/projects/ecommerce-api

API REST de e-commerce escrita em Go com `gin`, autenticação por JWT, PostgreSQL e GORM.

## Visão geral

O projeto está organizado por domínio em `internal/`:

- `database`: conexão com PostgreSQL, criação do banco e bootstrap do schema com GORM.
- `user`: cadastro, login e leitura do usuário autenticado.
- `product`: criação e listagem de produtos.
- `cart`: gerenciamento de carrinho do usuário autenticado.
- `order`: checkout e marcação de pedido como pago.
- `payment`: integração com Stripe para criar `PaymentIntent` e receber webhook.
- `auth`: geração de token JWT e middlewares de autenticação/autorização.

Ao subir a aplicação, o banco configurado é criado automaticamente se ainda não existir, junto com as tabelas necessárias.

## Stack

- Go `1.26.1`
- Gin
- JWT (`github.com/golang-jwt/jwt/v5`)
- PostgreSQL
- GORM (`gorm.io/gorm`)
- Driver PostgreSQL do GORM (`gorm.io/driver/postgres`)
- Stripe (`github.com/stripe/stripe-go/v85`)
- UUID (`github.com/google/uuid`)
- Bcrypt (`golang.org/x/crypto/bcrypt`)

## Como rodar

1. Garanta que o PostgreSQL esteja rodando.
2. Configure as variáveis de ambiente no shell ou `.env`:

```bash
export POSTGRES_HOST=localhost
export POSTGRES_PORT=5432
export POSTGRES_USER=postgres
export POSTGRES_PASSWORD=123mudar
export POSTGRES_DB=ecommerce
export POSTGRES_SSLMODE=disable
export STRIPE_SECRET_KEY=sk_test_xxx
```

3. Inicie a API:

```bash
go run ./cmd/api
```

A aplicação sobe em `http://localhost:8080`.

Se o banco `ecommerce` ainda não existir, a API tentará criá-lo automaticamente usando a conexão administrativa no banco `postgres`. Depois disso, o GORM executa `AutoMigrate` para garantir o schema básico.

Se quiser subir o PostgreSQL via Docker:

```bash
docker compose up -d postgres
```

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

- `secret` JWT fixo no código.
- Segredo do webhook Stripe ainda está hardcoded.
- Bootstrap de schema feito com `AutoMigrate`; ainda não há versionamento formal de migrations.
- Alguns fluxos ainda dependem de melhorias de validação e tratamento de erro.
