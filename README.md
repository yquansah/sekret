# sekret

A CLI tool that simplifies Kubernetes secret management by reading environment variables from dotenv files and automatically creating or updating secrets in your cluster.

## Background

Managing secrets in Kubernetes can be cumbersome. To update a secret, you typically need to:

1. Base64 encode each value manually
2. Edit the secret YAML or use complex `kubectl` commands
3. Apply the changes to your cluster

This process is error-prone and time-consuming, especially when dealing with multiple environment variables or frequent updates.

**sekret** streamlines this workflow by automatically reading from dotenv files (`.env`) and handling the base64 encoding and Kubernetes API interactions for you.

## What it does

sekret reads key-value pairs from a dotenv file and creates or updates a Kubernetes secret with those values. It handles:

- ✅ Automatic base64 encoding of values
- ✅ Creating new secrets if they don't exist
- ✅ Updating existing secrets (merge or replace modes)
- ✅ Proper kubeconfig authentication
- ✅ Detailed feedback on what was modified

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

### Basic Usage

```bash
sekret upsert my-secret --namespace default --env-file .env
```

### Command Syntax

```bash
sekret upsert [secret-name] [flags]
```

### Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--namespace` | `-n` | `default` | Kubernetes namespace |
| `--env-file` | `-f` | `.env` | Path to dotenv file |
| `--replace` | | `false` | Replace all existing values instead of merging |

### Examples

#### Create or update a secret (merge mode)
```bash
sekret upsert my-app-secret --namespace production --env-file .env.prod
```

If the secret already exists, new keys from `.env.prod` will be added and existing keys will be updated. Keys not in the `.env.prod` file will remain unchanged.

#### Replace all secret values
```bash
sekret upsert my-app-secret --env-file .env --replace
```

This will completely replace the secret's data with only the values from `.env`. Any existing keys not in the `.env` file will be removed.

#### Use default settings
```bash
sekret upsert my-secret
```

This reads from `./.env` and creates/updates the secret in the `default` namespace.

## Dotenv File Format

Your `.env` file should contain key-value pairs:

```bash
DATABASE_URL=postgresql://user:password@localhost:5432/mydb
API_KEY=secret-api-key-123
DEBUG=true
REDIS_URL=redis://localhost:6379
```

## Behavior

### Merge Mode (Default)

When you run `sekret upsert` without `--replace`:

- **New secret**: Creates the secret with all key-value pairs from the dotenv file
- **Existing secret**: Adds new keys and updates existing keys from the dotenv file, preserves other existing keys

### Replace Mode

When you run `sekret upsert --replace`:

- **New secret**: Creates the secret with all key-value pairs from the dotenv file
- **Existing secret**: Completely replaces the secret's data with only the values from the dotenv file

### Key Modification Counting

sekret tracks and reports the number of keys that were actually modified:

- **Merge mode**: Counts only keys that were added or had their values changed
- **Replace mode**: Counts all keys from the dotenv file

Example output:
```bash
Successfully upserted 3 key(s) for secret 'my-secret' in namespace 'default'
Successfully replaced 5 key(s) for secret 'my-secret' in namespace 'default'
```

## License

[Add your preferred license here]