# 🍽️ SaaS Core — Backend (Go + Gin)

Backend core de um sistema de **marketplace de restaurantes + gestão de cardápio**, inspirado em plataformas como **Menudino / AnotaAi**.

O projeto foi pensado desde o início para:

- separar bem domínio, aplicação e infraestrutura
- suportar **usuários comuns** e **donos de loja**
- permitir **cardápios altamente configuráveis**
- escalar futuramente (ex: CassandraDB, pedidos, pagamentos)

---

## 🎯 Objetivo do Projeto

Criar um backend sólido para:

- cadastro e autenticação de usuários
- criação e gestão de lojas
- criação de cardápios por loja
- estrutura pronta para categorias, itens, variações e adicionais
- tomada de decisão de fluxo já no login (UX simples para o front)

---

## 🏗️ Arquitetura

Arquitetura baseada em **Clean Architecture / Hexagonal**:

```

cmd/api              → Entrypoint da aplicação
internal/domain      → Domínio puro (entities, regras)
internal/usecase     → Casos de uso (regras de aplicação)
internal/infra       → HTTP, middleware, repos in-memory
internal/port        → Interfaces (ports)
pkg                  → Infra técnica (jwt, password, uuid)

```

### Princípios

- Domínio não conhece HTTP nem banco
- Usecases orquestram regras de negócio
- Infra apenas implementa contratos
- Fácil troca de persistência (ex: CassandraDB)

---

## 🔐 Autenticação & Sessões

- Autenticação via **JWT (HS256)**
- JWT possui `jti` (ID único do token)
- Sessões são **stateful**
- Cada login gera uma sessão persistida em repositório

### Middleware de autenticação valida:

- Header `Authorization`
- JWT válido
- Sessão existente
- Expiração do token

### Claims do JWT

- `user_id`
- `role` (uso interno)
- `exp`, `iat`, `iss`, `jti`

---

## 👤 Tipos de Usuário e Fluxos

O sistema trabalha com **dois tipos de usuários**, sem expor roles internas ao frontend.

### Entrada no cadastro (`POST /user`)

```json
{
  "user_type": "customer" | "store"
}
```

### Tipos

- **customer** → usuário comum (faz pedidos)
- **store** → usuário que irá criar e gerenciar loja

---

## 🔁 Fluxo de Login (Opção A)

No login, o backend já decide **para onde o usuário deve ir**.

### `POST /login` → Response

```json
{
  "token": "jwt...",
  "user": {
    "id": "u1",
    "name": "Fulano",
    "kind": "store"
  },
  "stores_count": 0,
  "next_step": "CREATE_STORE"
}
```

### Possíveis `next_step`

- `BROWSE_STORES` → usuário comum
- `CREATE_STORE` → dono sem loja
- `STORE_DASHBOARD` → dono com loja(s)

### Regra de decisão

```
if kind == customer:
  BROWSE_STORES
else if kind == store and stores_count == 0:
  CREATE_STORE
else:
  STORE_DASHBOARD
```

---

## 🧭 Diagrama do Fluxo de Onboarding

```
┌─────────────────────────┐
│  SIGNUP (POST /user)    │
│  user_type              │
└────────────┬────────────┘
             │
             ▼
┌─────────────────────────┐
│   LOGIN (POST /login)   │
│ token + next_step       │
└────────────┬────────────┘
             │
             ▼
┌────────────────────────────────────────┐
│ Backend decide next_step               │
│                                        │
│ customer → BROWSE_STORES               │
│ store + 0 lojas → CREATE_STORE         │
│ store + lojas → STORE_DASHBOARD        │
└────────────┬───────────────────────────┘
             │
     ┌───────┼────────────┐
     ▼       ▼            ▼
┌────────┐ ┌──────────┐ ┌─────────────────┐
│Browse  │ │Create    │ │Store Dashboard  │
│Stores  │ │Store     │ │Menus / Gestão   │
└────────┘ └──────────┘ └─────────────────┘
```

---

## 🏪 Domínio de Loja

### Store

- pertence a um usuário (Owner)
- um usuário pode ter **0 ou mais lojas**
- apenas usuários autenticados podem criar loja

---

## 🍽️ Domínio de Cardápio

O cardápio foi modelado para suportar **configurações complexas** (pizza, combos, adicionais, variações).

### Hierarquia

```
Store
 └── Menu
      └── Category
           └── Item
                ├── VariantGroup
                │     └── VariantOption
                └── AddonGroup
                      └── AddonOption
```

---

## 💰 Regra de Preço

```
Preço final =
Item.BasePrice
+ soma(VariantOption.PriceDelta)
+ soma(AddonOption.Price * quantidade)
```

- Variações usam `PriceDelta`
- Adicionais usam `Price`

---

## 🌐 Rotas da API (Estado Atual)

### Públicas

- `POST /user` → cria usuário
- `POST /login` → login (retorna token + next_step)

---

### Protegidas (JWT)

#### Store

- `POST /store` → cria loja
- `GET /store/id/:id` → busca loja por id

#### Menus

- `POST /store/:storeId/menu` → cria menu para a loja
- `GET /store/:storeId/menus` → lista menus da loja
- `GET /menu/:id` → busca menu por id

#### User

- `GET /user/:id`
- `GET /user/email/:email`
- `GET /user/cpf/:cpf`

---

## 🗄️ Persistência (Atual)

- Repositórios **in-memory**
- Estrutura preparada para CassandraDB
- Índices simulados:
  - `byID`
  - `byOwner`
  - `byStore`
  - `byMenu`
  - `byCategory`
  - `byItem`

---

## 🧪 Testes

- Value Objects:
  - CPF
  - CNPJ
- Password hash (bcrypt)

---

## 🛣️ Próximos Passos Planejados

- Categorias do menu
- Itens do menu
- Variações de itens
- Adicionais
- `GET MenuFull` (JSON completo para o app)
- Pedido / Checkout
- Estoque
- Pagamentos
- Migração para CassandraDB
- Observabilidade (logs e métricas)

---

## 🧑‍💻 Status do Projeto

🚧 Em desenvolvimento ativo
🧠 Arquitetura definida
🔐 Autenticação sólida
🍽️ Base de cardápio pronta para evoluir

---

## Como rodar a aplicação

```bash
docker compose up
```

---

## Como rodar os testes

```bash
docker compose exec api go test -v ./...
```
