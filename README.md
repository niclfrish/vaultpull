# vaultpull

> CLI tool to sync HashiCorp Vault secrets to local `.env` files with namespace support

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/vaultpull.git
cd vaultpull
go build -o vaultpull .
```

---

## Usage

Set your Vault address and token, then pull secrets into a `.env` file:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.yourtoken"

# Pull secrets from a path into a .env file
vaultpull pull --path secret/myapp/prod --out .env

# Use a namespace (Vault Enterprise)
vaultpull pull --path secret/myapp/prod --namespace team-a --out .env.production

# Preview without writing
vaultpull pull --path secret/myapp/prod --dry-run
```

The resulting `.env` file will contain key-value pairs sourced directly from the specified Vault path:

```env
DATABASE_URL=postgres://user:pass@host/db
API_KEY=abc123
DEBUG=false
```

---

## Configuration

| Flag | Env Var | Description |
|------|---------|-------------|
| `--path` | — | Vault secret path |
| `--out` | `VAULTPULL_OUT` | Output `.env` file path |
| `--namespace` | `VAULT_NAMESPACE` | Vault namespace |
| `--dry-run` | — | Print secrets without writing |

---

## License

[MIT](LICENSE)