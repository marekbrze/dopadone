# Turso Getting Started Guide

This guide walks you through setting up a Turso account, creating databases, and obtaining credentials for use with Dopadone.

## Overview

Turso is a SQLite-compatible database platform that provides cloud-hosted databases with features like embedded replicas, branching, and automatic backups. Dopadone supports three connection modes to Turso:

| Mode | Description | Best For |
|------|-------------|----------|
| **Remote** | Direct connection to Turso cloud | Always-online environments |
| **Replica** | Local replica with cloud sync | Offline-capable with cloud backup |

## Quick Start (5 Minutes)

If you're in a hurry, follow these minimal steps:

```bash
# 1. Install Turso CLI
brew install tursodatabase/tap/turso      # macOS
# curl -sSfL https://get.tur.so/install.sh | bash  # Linux

# 2. Sign up / Login
turso auth signup

# 3. Create a database
turso db create dopadone

# 4. Get your credentials
turso db show dopadone --url              # Database URL
turso db tokens create dopadone           # Auth token
```

---

## Step 1: Account Signup

### Creating a Turso Account

1. Navigate to [turso.tech](https://turso.tech)
2. Click **"Start for free"** or **"Sign up"**
3. Choose your signup method:
   - GitHub account
   - Google account
   - Email address

### Signup via CLI

Alternatively, use the CLI to create an account:

```bash
turso auth signup
```

This opens your browser to complete signup. For headless environments (CI/CD, WSL):

```bash
turso auth signup --headless
```

### Free Tier Limits

The free Starter plan includes:

| Resource | Limit |
|----------|-------|
| Databases | 1 group, up to 3 databases |
| Storage | 1 GB total |
| Row reads | 1 billion/month |
| Row writes | 25 million/month |

### Signup Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Turso Account Signup                     │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │   turso.tech or turso auth    │
              │           signup              │
              └───────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │    Choose signup method:      │
              │  • GitHub  • Google  • Email  │
              └───────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │      Verify email if          │
              │      using email signup       │
              └───────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │     Account created!          │
              │   Ready to create databases   │
              └───────────────────────────────┘
```

---

## Step 2: CLI Installation

The Turso CLI is required to manage databases from the command line.

### macOS

**Option 1: Homebrew (Recommended)**

```bash
brew install tursodatabase/tap/turso
```

**Option 2: Shell Script**

```bash
curl -sSfL https://get.tur.so/install.sh | bash
```

### Linux

```bash
curl -sSfL https://get.tur.so/install.sh | bash
```

### Windows

Turso requires [WSL (Windows Subsystem for Linux)](https://learn.microsoft.com/en-us/windows/wsl/install):

```powershell
# In PowerShell
wsl

# Then in WSL
curl -sSfL https://get.tur.so/install.sh | bash
```

### Verify Installation

```bash
# Open a new shell, then run:
turso --version
```

### Installation Paths Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Turso CLI Installation                    │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
   ┌─────────┐          ┌─────────┐          ┌─────────┐
   │ macOS   │          │ Linux   │          │ Windows │
   └─────────┘          └─────────┘          └─────────┘
        │                     │                     │
        ▼                     ▼                     ▼
 ┌──────────────┐     ┌──────────────┐     ┌──────────────┐
 │   Homebrew   │     │    Shell     │     │  WSL + Shell │
 │   or curl    │     │    script    │     │    script    │
 └──────────────┘     └──────────────┘     └──────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              ▼
              ┌───────────────────────────────┐
              │       Verify with:            │
              │       turso --version         │
              └───────────────────────────────┘
```

### Update CLI

To update to the latest version:

```bash
# Homebrew
brew upgrade tursodatabase/tap/turso

# Shell script (re-run installer)
curl -sSfL https://get.tur.so/install.sh | bash
```

---

## Step 3: Database Creation

### Via CLI

Create a new database:

```bash
# Basic creation (uses default group)
turso db create dopadone

# Create in a specific group
turso db create dopadone --group my-group

# Create and wait until ready
turso db create dopadone --wait
```

### Via Web UI

1. Log in to [Turso Dashboard](https://app.turso.tech)
2. Navigate to **Databases**
3. Click **Create Database**
4. Enter database name (e.g., `dopadone`)
5. Select group (or create new group)
6. Click **Create**

### Import Existing Data

Create from an existing SQLite file:

```bash
# From SQLite file (max 2GB)
turso db create dopadone --from-file ./existing.db

# From SQL dump
turso db create dopadone --from-dump ./dump.sql

# From CSV file
turso db create dopadone --from-csv ./data.csv --csv-table-name my_table

# Copy from another Turso database
turso db create dopadone-backup --from-db dopadone
```

### Database Creation Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                    Database Creation                         │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┴─────────────────────┐
        ▼                                           ▼
┌───────────────────┐                    ┌───────────────────┐
│   CLI Method      │                    │   Web UI Method   │
│                   │                    │                   │
│ turso db create   │                    │ Dashboard >       │
│    <name>         │                    │ Create Database   │
└───────────────────┘                    └───────────────────┘
        │                                           │
        └─────────────────────┬─────────────────────┘
                              ▼
              ┌───────────────────────────────┐
              │    Database created in        │
              │    selected group/region      │
              └───────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │   Verify with:                │
              │   turso db show <name>        │
              └───────────────────────────────┘
```

### List Databases

```bash
# List all databases
turso db list

# Show database details
turso db show dopadone
```

---

## Step 4: Authentication Tokens

Authentication tokens allow your applications to connect to Turso databases securely.

### Database Tokens (Recommended)

Create a token for a specific database:

```bash
# Create a full-access token
turso db tokens create dopadone

# Create a read-only token
turso db tokens create dopadone --read-only

# Create a token that expires in 7 days
turso db tokens create dopadone --expiration 7d

# Create a token that never expires (use with caution)
turso db tokens create dopadone --expiration never
```

### Token Scopes

| Scope | Access Level | Use Case |
|-------|--------------|----------|
| **Full access** | Read + Write | Production apps, Dopadone |
| **Read-only** | Read only | Analytics, reporting dashboards |

### Token Expiration Options

```bash
--expiration 7d      # 7 days
--expiration 30d     # 30 days
--expiration 1d12h   # 1 day 12 hours
--expiration never   # No expiration (not recommended)
```

### Token Security Best Practices

1. **Use read-only tokens** when possible (analytics, reporting)
2. **Set expiration dates** for non-production tokens
3. **Never commit tokens** to source control
4. **Store tokens in environment variables** or secure secret management
5. **Rotate tokens periodically** for production databases

### Invalidate Tokens

```bash
# Invalidate all tokens for a database
turso db tokens invalidate dopadone
```

### Token Flow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                   Token Generation                           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │   turso db tokens create      │
              │        <database>             │
              └───────────────────────────────┘
                              │
              ┌───────────────┴───────────────┐
              ▼                               ▼
     ┌────────────────┐              ┌────────────────┐
     │   Full Access  │              │   Read-Only    │
     │   (default)    │              │  --read-only   │
     └────────────────┘              └────────────────┘
              │                               │
              └───────────────┬───────────────┘
                              ▼
              ┌───────────────────────────────┐
              │   Set expiration (optional)   │
              │   --expiration 7d|30d|never   │
              └───────────────────────────────┘
                              │
                              ▼
              ┌───────────────────────────────┐
              │   Store securely:             │
              │   Environment variable or     │
              │   secret manager              │
              └───────────────────────────────┘
```

---

## Step 5: Database URL

The database URL identifies your Turso database for connections.

### Get URL via CLI

```bash
# Show database URL
turso db show dopadone --url

# Output example:
# libsql://dopadone-organization.turso.io
```

### Get URL via Web UI

1. Log in to [Turso Dashboard](https://app.turso.tech)
2. Navigate to **Databases**
3. Click on your database name
4. Copy the **URL** from the database details

### URL Format

```
libsql://<database-name>-<organization>.turso.io

Example:
libsql://dopadone-myorg.turso.io
```

### URL Structure Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                     Turso Database URL                       │
└─────────────────────────────────────────────────────────────┘

  libsql://dopadone-myorganization.turso.io
  │       │        │               │
  │       │        │               └── Turso domain
  │       │        └── Organization slug
  │       └── Database name
  └── Protocol (libsql)

Components:
  • Protocol: Always "libsql" (Turso's protocol)
  • Database: Your database name
  • Organization: Your Turso organization slug
  • Domain: turso.io (cloud endpoint)
```

---

## Configuring Dopadone with Turso

Once you have your credentials, configure Dopadone:

### Option 1: Environment Variables

```bash
# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export TURSO_DATABASE_URL="libsql://dopadone-myorg.turso.io"
export TURSO_AUTH_TOKEN="your-auth-token-here"

# Use remote mode
export DOPA_DB_MODE=remote

# Or use replica mode (recommended for offline support)
export DOPA_DB_PATH="./dopadone-replica.db"
export DOPA_DB_MODE=replica
```

### Option 2: YAML Configuration File

Create `dopadone.yaml`:

```yaml
database:
  # Remote mode configuration
  mode: remote
  turso:
    url: libsql://dopadone-myorg.turso.io
    token: ${TURSO_AUTH_TOKEN}  # Reference env variable
```

For replica mode:

```yaml
database:
  # Replica mode configuration
  path: ./dopadone-replica.db
  mode: replica
  sync_interval: 60s
  turso:
    url: libsql://dopadone-myorg.turso.io
    token: ${TURSO_AUTH_TOKEN}
```

### Option 3: CLI Flags

```bash
# Remote mode
dopa --turso-url "libsql://dopadone-myorg.turso.io" \
     --turso-auth-token "your-token" \
     --db-mode remote \
     tasks list

# Replica mode
dopa --db ./dopadone-replica.db \
     --turso-url "libsql://dopadone-myorg.turso.io" \
     --turso-auth-token "your-token" \
     --db-mode replica \
     tasks list
```

### Verify Connection

```bash
# Test connection
dopa areas list
```

---

## Troubleshooting

### CLI Installation Fails

**Symptom**: `command not found: turso` after installation

**Solution**:
1. Open a new shell/terminal window
2. Verify PATH includes the installation directory
3. Re-run the installer

### Authentication Errors

**Symptom**: `authentication failed` or `unauthorized`

**Solution**:
1. Verify token is valid: `turso auth token`
2. Re-login: `turso auth login`
3. Check token hasn't expired
4. Ensure token has correct permissions

### Database Connection Issues

**Symptom**: `connection refused` or `timeout`

**Solution**:
1. Verify database URL is correct
2. Check internet connectivity
3. Verify database exists: `turso db list`
4. Check Turso status page for outages

### Token Not Working

**Symptom**: `invalid token` error

**Solution**:
1. Verify token format (should be a long JWT)
2. Ensure token was created for correct database
3. Check if token was invalidated: `turso db tokens invalidate` removes all tokens
4. Create a new token

### Headless Authentication

For CI/CD or WSL environments:

```bash
# Use headless flag for browser-less auth
turso auth login --headless
```

---

## Learn More

### Official Turso Documentation

| Topic | Link |
|-------|------|
| Quickstart | [docs.turso.tech/quickstart](https://docs.turso.tech/quickstart) |
| CLI Reference | [docs.turso.tech/cli/introduction](https://docs.turso.tech/cli/introduction) |
| Embedded Replicas | [docs.turso.tech/features/embedded-replicas](https://docs.turso.tech/features/embedded-replicas/introduction) |
| Usage & Billing | [docs.turso.tech/help/usage-and-billing](https://docs.turso.tech/help/usage-and-billing) |
| Go SDK | [docs.turso.tech/sdk/go/quickstart](https://docs.turso.tech/sdk/go/quickstart) |

### Dopadone Documentation

| Document | Description |
|----------|-------------|
| [Database Modes](DATABASE_MODES.md) | Detailed mode explanations and configuration |
| [Turso Migrations](TURSO_MIGRATIONS.md) | Migration guide for libSQL/Turso integration |
| [Architecture Overview](architecture/01-overview.md) | System architecture |

### Community & Support

- [Turso Discord](https://tur.so/discord) - Community support
- [Turso GitHub](https://github.com/tursodatabase/turso) - Issues and discussions
- [Turso Twitter](https://twitter.com/tursodatabase) - Updates and announcements

---

## Summary

To use Turso with Dopadone:

1. **Install CLI**: `brew install tursodatabase/tap/turso`
2. **Create account**: `turso auth signup`
3. **Create database**: `turso db create dopadone`
4. **Get URL**: `turso db show dopadone --url`
5. **Create token**: `turso db tokens create dopadone`
6. **Configure Dopadone**: Set `TURSO_DATABASE_URL` and `TURSO_AUTH_TOKEN` environment variables

For production use, consider:
- Using **replica mode** for offline support
- Setting **token expiration** for security
- Storing tokens in **environment variables** (never in code)
