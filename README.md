# filebin

A fast, dependency-free CLI for [Filebin](https://filebin.net) — upload, download, list, and manage files directly from your terminal.

[![Go version](https://img.shields.io/badge/Go-%3E%3D1.26-00ADD8?logo=go&logoColor=white)](go.mod)
[![Release](https://img.shields.io/github/v/release/djsnipa1/filebin-cli?logo=github&label=release)](https://github.com/djsnipa1/filebin-cli/releases)
[![License](https://img.shields.io/github/license/djsnipa1/filebin-cli)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey?logo=linux&logoColor=white)](README.md)
[![Last commit](https://img.shields.io/github/last-commit/djsnipa1/filebin-cli)](https://github.com/djsnipa1/filebin-cli)
[![Repo size](https://img.shields.io/github/repo-size/djsnipa1/filebin-cli)](https://github.com/djsnipa1/filebin-cli)

`filebin` makes [Filebin](https://filebin.net/api) easy to use from the command line. No registration required — everyone with a bin URL can read and write to it. Great for quick sharing between machines, distributing files to a team, or scripting uploads in CI.

## Features

- ⬆️ **Upload** with live progress bar and automatic SHA-256 integrity verification
- ⬇️ **Download** with live progress bar, handling Filebin's cookie verification automatically
- 📋 **List** bin contents as a readable table or raw JSON (for scripting)
- 🗑️ **Delete** a single file or an entire bin
- 🔒 **Lock** a bin to make it read-only
- 📦 **Archive** an entire bin as a `.tar` or `.zip`
- 🔳 **QR code** for a bin, ready to scan on mobile
- 🔑 **Checksums** in `sha256sum` format for easy integrity checking

## Installation

### Quick install (curl)

Install the latest release binary (Linux amd64) to `~/.local/bin`:

```bash
curl -fsSL https://raw.githubusercontent.com/djsnipa1/filebin-cli/main/install.sh | sh
```

> [!NOTE]
> The script installs to `$HOME/.local/bin` by default. Override with `FILEBIN_BIN_DIR`:
> `FILEBIN_BIN_DIR=/usr/local/bin curl -fsSL ... | sh`
> If that directory isn't on your `PATH`, add it:
> `echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc`

### Build from source

Requires [Go](https://go.dev/dl/) 1.26+.

```bash
go build -o filebin .
```

> [!TIP]
> The binary is fully self-contained — no runtime dependencies. Copy it anywhere and run.

### Using Nix

```bash
nix-shell -p go --run "go build -o filebin ."
```

## Usage

```bash
filebin upload <bin> <file>       # upload a file
filebin download <bin> <file>     # download a file
filebin list <bin>                # list bin contents
filebin delete <bin> [file]       # delete a file, or whole bin
filebin lock <bin>                # make a bin read-only
filebin archive <bin> --format zip
filebin qr <bin>                  # save a QR code PNG
filebin checksums <bin>           # SHA-256 checksums
```

### Upload

```bash
filebin upload myshare ./photo.jpg
# Uploading [████████████████████████] 100% 482 KB
# Uploaded photo.jpg to bin myshare
# URL: https://filebin.net/myshare/photo.jpg
```

### Download

```bash
filebin download myshare photo.jpg -o ./photo.jpg
```

### List plain or as JSON

```bash
filebin list myshare
FILE        SIZE   TYPE                      CREATED
photo.jpg   482 kB image/jpeg                1 minute ago

filebin list myshare --json   # machine-readable output
```

### Verify integrity

```bash
filebin checksums myshare | sha256sum -c
```

### Pointing at another server

Self-hosted Filebin or a mirror? Use `--base-url`:

```bash
filebin --base-url https://files.example.com upload myshare ./photo.jpg
```

### Archive a whole bin

```bash
filebin archive myshare --format tar -o myshare.tar
filebin archive myshare --format zip   # defaults to myshare.zip
```

## Command reference

| Command | Description |
|---------|-------------|
| `upload <bin> <file>` | Upload a file to a bin |
| `download <bin> <file>` | Download a file from a bin |
| `list <bin>` | Show bin metadata and files (table or `--json`) |
| `delete <bin> [file]` | Delete a file, or the whole bin if no file given |
| `lock <bin>` | Lock a bin (read-only) |
| `archive <bin> -f <tar\|zip>` | Download bin as an archive |
| `qr <bin>` | Save a QR code PNG for the bin |
| `checksums <bin>` | Print SHA-256 checksums (sha256sum format) |

**Global flags:**

- `--base-url <url>` — Filebin server URL (default `https://filebin.net`)
- `-h, --help` — show help for any command
- `-o, --output <path>` — custom output path for downloads/archives/QR

## Development

This is a Go project using [Cobra](https://github.com/spf13/cobra) for the CLI and [progressbar](https://github.com/schollz/progressbar) for progress display.

```bash
nix-shell -p go    # or: install Go 1.26+
go build -o filebin .
```

## Releases

A [GitHub Actions workflow](.github/workflows/release.yml) builds the binary for Linux, macOS, and Windows (amd64 + arm64) and attaches them to a release.

### Trigger via version tag (recommended)

```bash
git tag v1.0.0
git push origin v1.0.0
```

Pushing a `v*` tag automatically builds all platform binaries and creates a release with auto-generated notes.

### Trigger manually

From the **Actions** tab: select **Release** → **Run workflow**. An optional `tag` input can be set; if omitted, a timestamp tag (e.g. `v20260830-192500`) is used.

