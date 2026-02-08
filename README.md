# 🌿 Leafy

> [!WARNING]
> This project is currently in development and subject to changes. Use it at your own risk.

**Leafy** is your cozy terminal companion for transferring media from external drives. It turns the chore of manual file transfers into a breezy, interactive experience right in your shell.

### 💡 Motivation

I started vlogging my motorcycle journey recently, and found files transferring from my action camera and external mic is a *very* tedious task. So I created **Leafy** to automate the process. Now I can quickly and easily transfer all my media to my editing PC by simply running in the terminal! 🙌

### 🚀 Quickstart

A binary will be released on GitHub later. Alternatively, you can install it using Go.

```bash
# Option 1: Install with Go
go install github.com/Nightails/Leafy@latest

# Option 2: Download the binary
# 1. Download: Download the latest binary from the Releases page.
# 2. Permissions: Ensure the binary is executable:
chmod +x leafy
# 3. Run: Execute the binary:
./leafy
```

### 📖 Usage

Leafy provides an interactive interface to monitor your USB devices:

```text
- Automatic Scanning: The app starts scanning immediately upon launch.
- Manual Refresh: Press 's' to trigger a new scan.
- Navigation:
    - Use 'j' or 'k' (or arrow keys) to navigate.
    - Press 'enter' to select a specific device.
- Exiting: Press 'ctrl+c' to safely exit the application.
```

### 🤝 Contributing

To set up Leafy for local development, ensure you have the following dependencies installed:

-   **Go** (version 1.25 or later)
-   **lsblk**: The application relies on the `lsblk` command-line utility (usually pre-installed on most Linux distributions).

#### Development Steps

```bash
# 1. Clone the repository:
git clone https://github.com/user/leafy.git
cd leafy

# 2. Install dependencies:
go mod download

# 3. Run locally:
go run main.go

# 4. Build:
go build -o leafy main.go
```

### 📜 License

This project is licensed under the terms of the MIT License. See [LICENSE](LICENSE) for more details.
