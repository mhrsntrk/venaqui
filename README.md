# venaqui

**venaqui** is a command-line tool with a Terminal User Interface (TUI) that leverages Real-Debrid premium links and aria2 for high-speed downloads. The tool provides real-time download progress, statistics, and a polished user experience.

> **About the name**: "venaqui" comes from Spanish, meaning "come here" (ven aquí). It's an invitation for your downloads to come to you quickly and efficiently.

## Features

- 🚀 **High-Speed Downloads**: Uses aria2 with optimized settings (16 connections, 16 splits)
- 🔓 **Real-Debrid Integration**: Automatically unrestricts premium hoster links
- 📊 **Real-Time Progress**: Beautiful TUI with live download statistics and speed history graph
- 🎨 **Polished UI**: Modern terminal interface built with Bubble Tea featuring enhanced visuals
- 📈 **Advanced Statistics**: ETA, elapsed time, connections, upload speed, and more
- 🔍 **Post-Download Actions**: Open files directly or reveal in Finder/Explorer
- 🔧 **Cross-Platform**: Works on Windows, macOS, and Linux
- ⚙️ **Easy Configuration**: Simple YAML configuration file

## Prerequisites

### 1. Install aria2

**macOS:**
```bash
brew install aria2
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt install aria2
```

**Linux (Fedora/RHEL):**
```bash
sudo yum install aria2
```

**Windows:**
Download from [aria2 official website](https://aria2.github.io/) and add to PATH.

### 2. Get Real-Debrid API Token

1. Sign up for a [Real-Debrid account](http://real-debrid.com/?id=17216837) (referral link)
2. Visit [My Devices page](https://real-debrid.com/devices)
3. Copy your API Private Token

## Installation

### Homebrew (macOS) - Recommended

The easiest way to install venaqui on macOS is via Homebrew:

```bash
# Install via Homebrew tap
brew tap mhrsntrk/venaqui
brew install venaqui
```

Or install directly from the formula:

```bash
brew install mhrsntrk/venaqui/venaqui
```

**Note:** Homebrew will automatically install `aria2` as a dependency.

### Build from Source

```bash
# Clone the repository
git clone https://github.com/mhrsntrk/venaqui.git
cd venaqui

# Build
go build -o venaqui cmd/venaqui/main.go

# Install globally (optional)
go install cmd/venaqui/main.go
```

### Using Go Install

```bash
go install github.com/mhrsntrk/venaqui/cmd/venaqui@latest
```

## Configuration

Create a configuration file at `~/.venaqui/config.yaml`:

```yaml
realdebrid:
  api_token: "YOUR_RD_API_TOKEN"

aria2:
  rpc_url: "http://localhost:6800/jsonrpc"
  secret: ""  # Optional RPC secret

download:
  default_dir: ""  # Leave empty for OS default Downloads folder
```

### Configuration Options

- **realdebrid.api_token** (required): Your Real-Debrid API token
- **aria2.rpc_url** (optional): aria2 RPC endpoint (default: `http://localhost:6800/jsonrpc`)
- **aria2.secret** (optional): aria2 RPC secret if configured
- **download.default_dir** (optional): Default download directory (default: `~/Downloads`)

## Usage

### Basic Usage

```bash
# Download to default location (~/Downloads)
venaqui "https://mega.nz/file/example"

# Download to specific location
venaqui "https://mega.nz/file/example" "/path/to/download"

# With full path
venaqui "https://1fichier.com/example" "$HOME/Downloads/Movies"
```

### Supported Hosters

venaqui works with all hosters supported by Real-Debrid, including:
- MEGA
- 1fichier
- Rapidgator
- Uploaded
- And [many more](https://real-debrid.com/hosters)

## How It Works

1. **Link Unrestriction**: venaqui sends your hoster link to Real-Debrid API
2. **Direct Link**: Real-Debrid returns an unrestricted direct download link
3. **aria2 Download**: The direct link is passed to aria2 for high-speed downloading
4. **TUI Display**: Real-time progress is shown in a beautiful terminal interface

## TUI Controls

### During Download
- **q** or **Ctrl+C** or **Esc**: Quit the application

### After Download Completes
- **o**: Open the downloaded file with its default application
- **d** or **s**: Show the file in Finder/Explorer (reveals and highlights the file)
- **q**: Quit the application

The TUI stays open after download completion, allowing you to interact with the downloaded file.

## Project Structure

```
venaqui/
├── cmd/venaqui/main.go              # Entry point
├── internal/
│   ├── realdebrid/                  # Real-Debrid API integration
│   ├── aria2/                       # aria2 RPC client
│   ├── tui/                         # Bubble Tea TUI
│   ├── config/                      # Configuration management
│   └── utils/                       # Utilities
├── pkg/models/                      # Shared models
├── go.mod
└── README.md
```

## Troubleshooting

### aria2 Not Found

Ensure aria2 is installed and accessible in your PATH:
```bash
which aria2c  # Should show path to aria2c
```

### Real-Debrid API Errors

- **401 Unauthorized**: Check that your API token is correct in `~/.venaqui/config.yaml`
- **429 Rate Limited**: Too many requests, wait a moment and try again
- **Link Not Supported**: The hoster may not be supported by Real-Debrid

### aria2 Connection Errors

- Ensure aria2 RPC is enabled: `aria2c --enable-rpc`
- Check if port 6800 is available
- Verify RPC URL in configuration matches your aria2 setup

### Download Directory Issues

- Ensure the download directory exists and is writable
- Use absolute paths for download directories
- Check disk space availability

## Development

### Dependencies

- Go 1.22 or later
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Terminal styling
- [arigo](https://github.com/siku2/arigo) - aria2 RPC client
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [Viper](https://github.com/spf13/viper) - Configuration management
- [go-humanize](https://github.com/dustin/go-humanize) - Human-readable sizes

### Running Tests

```bash
go test ./...
```

### Building

```bash
# Build for current platform
go build -o venaqui cmd/venaqui/main.go

# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o venaqui-linux cmd/venaqui/main.go
GOOS=darwin GOARCH=amd64 go build -o venaqui-macos cmd/venaqui/main.go
GOOS=windows GOARCH=amd64 go build -o venaqui.exe cmd/venaqui/main.go
```

## License

This project is open source. Please check the license file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

- [Real-Debrid](http://real-debrid.com/?id=17216837) for the premium link service ([sign up with referral](http://real-debrid.com/?id=17216837))
- [aria2](https://aria2.github.io/) for the powerful download engine
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the TUI framework

## Future Enhancements

- Configuration wizard on first run
- Multiple simultaneous downloads
- Download queue management
- Pause/resume functionality
- Download history
- Bandwidth limiting
- Torrent magnet link support
- Automatic file verification
- Notifications on completion
- Dark/light theme switching
