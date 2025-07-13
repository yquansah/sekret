# sekret

A CLI tool that simplifies Kubernetes secret management. It can create or update secrets from dotenv files, list existing secret contents, and delete specific keys from secrets.

## Background

Managing secrets in Kubernetes can be cumbersome. To update a secret, you typically need to:

1. Base64 encode each value manually
2. Edit the secret YAML or use complex `kubectl` commands
3. Apply the changes to your cluster

This process is error-prone and time-consuming, especially when dealing with multiple environment variables or frequent updates.

**sekret** streamlines this workflow by automatically reading from dotenv files (`.env`) and handling the base64 encoding and Kubernetes API interactions for you.

## What it does

sekret provides comprehensive Kubernetes secret management capabilities:

- ✅ **Create/Update**: Read from dotenv files and create or update secrets with automatic base64 encoding
- ✅ **List**: Display secret contents in JSON format with base64 decoded values
- ✅ **Delete Keys**: Remove specific keys from secrets with interactive or non-interactive modes
- ✅ **Detailed Feedback**: Reports on what was modified in each operation

## Installation

### Option 1: Install directly with go install

```bash
go install github.com/yquansah/sekret@latest
```

This will install sekret to your `~/go/bin` directory. Make sure `~/go/bin` is in your PATH.

### Option 2: Build from source

1. Clone this repository
2. Build the binary:
   ```bash
   go build -o sekret
   ```
3. Move to your PATH (optional):
   ```bash
   mv sekret /usr/local/bin/
   ```

## Usage

sekret provides three main commands: `upsert`, `list`, and `delete-keys`.

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--kubeconfig` | `~/.kube/config` | Path to kubeconfig file |

### Command: upsert

Create or update secrets from dotenv files.

```bash
sekret upsert [secret-name] [flags]
```

#### Flags
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--namespace` | `-n` | `default` | Kubernetes namespace |
| `--env-file` | `-f` | `.env` | Path to dotenv file |
| `--replace` | | `false` | Replace all existing values instead of merging |

### Command: list

Display secret contents in JSON format with decoded values.

```bash
sekret list [secret-name] [flags]
```

#### Flags
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--namespace` | `-n` | `default` | Kubernetes namespace |

### Command: delete-keys

Delete specific keys from a secret (interactive or non-interactive).

```bash
sekret delete-keys [secret-name] [flags]
```

#### Flags
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--namespace` | `-n` | `default` | Kubernetes namespace |
| `--keys` | | | Comma-separated list of keys to delete (non-interactive mode) |

### Examples

#### Create or update a secret (merge mode)
```bash
sekret upsert my-app-secret --namespace production --env-file .env.prod
```

#### Replace all secret values
```bash
sekret upsert my-app-secret --env-file .env --replace
```

#### List secret contents
```bash
sekret list my-secret --namespace production
```

Output:
```json
{
  "data": [
    {"key": "DATABASE_URL", "value": "postgresql://user:password@localhost:5432/mydb"},
    {"key": "API_KEY", "value": "secret-api-key-123"}
  ]
}
```

#### Delete keys interactively
```bash
sekret delete-keys my-secret --namespace production
```

This will show a menu to select which key to delete and confirm the action.

#### Delete specific keys
```bash
sekret delete-keys my-secret --namespace production --keys API_KEY,OLD_TOKEN
```

#### Use custom kubeconfig
```bash
sekret list my-secret --kubeconfig /path/to/custom/kubeconfig
```

## Dotenv File Format

Your `.env` file should contain key-value pairs:

```bash
DATABASE_URL=postgresql://user:password@localhost:5432/mydb
API_KEY=secret-api-key-123
DEBUG=true
REDIS_URL=redis://localhost:6379
```

## Behavior

### Upsert Command

#### Merge Mode (Default)
- **New secret**: Creates the secret with all key-value pairs from the dotenv file
- **Existing secret**: Adds new keys and updates existing keys, preserves other existing keys

#### Replace Mode
- **New secret**: Creates the secret with all key-value pairs from the dotenv file  
- **Existing secret**: Completely replaces the secret's data with only values from the dotenv file

### List Command

- Retrieves all key-value pairs from the specified secret
- Automatically decodes base64 values for display
- Returns data in JSON format with `{"key": "...", "value": "..."}` structure

### Delete Keys Command

#### Interactive Mode (default)
- Lists all available keys in the secret
- Prompts user to select which key to delete
- Asks for confirmation before deletion

#### Non-Interactive Mode (`--keys` flag)
- Deletes specified comma-separated keys directly
- Reports the number of keys successfully deleted

### Output Messages

sekret provides clear feedback on operations:

```bash
# Upsert examples
Successfully upserted 3 key(s) for secret 'my-secret' in namespace 'default'
Successfully replaced 5 key(s) for secret 'my-secret' in namespace 'default'

# Delete examples  
Successfully deleted 2 key(s) from secret 'my-secret' in namespace 'default'
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.