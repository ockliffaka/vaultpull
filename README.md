# vaultpull

> CLI tool to sync secrets from HashiCorp Vault into local `.env` files safely

---

## Installation

```bash
go install github.com/yourusername/vaultpull@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpull/releases).

---

## Usage

Set your Vault address and token, then run `vaultpull` pointing at a Vault path:

```bash
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

vaultpull --path secret/data/myapp --output .env
```

This will fetch all key/value pairs stored at the given Vault path and write them to `.env`:

```
DB_HOST=db.example.com
DB_PASSWORD=supersecret
API_KEY=abc123
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--path` | *(required)* | Vault secret path to read from |
| `--output` | `.env` | Output file path |
| `--append` | `false` | Append to existing file instead of overwriting |
| `--dry-run` | `false` | Print secrets to stdout without writing |

```bash
# Preview without writing
vaultpull --path secret/data/myapp --dry-run

# Append to an existing .env file
vaultpull --path secret/data/shared --output .env --append
```

---

## Requirements

- Go 1.21+
- A running HashiCorp Vault instance
- A valid `VAULT_TOKEN` or other supported auth method

---

## License

[MIT](LICENSE) © 2024 yourusername