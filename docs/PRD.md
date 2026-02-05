# Venaqui - PRD

### Overview

**Venaqui** is a command-line tool with a Terminal User Interface (TUI) that leverages Real-Debrid premium links and aria2 for high-speed downloads. The tool provides real-time download progress, statistics, and a polished user experience using the Bubble Tea framework.[^1][^2][^3][^4]

### Architecture

#### Component Stack

1. **CLI Layer**: Command parsing and argument handling
2. **Real-Debrid Integration**: API communication for link unrestriction
3. **aria2 RPC Client**: Download management via JSON-RPC interface
4. **TUI Layer**: Bubble Tea framework for interactive display
5. **File System Handler**: Cross-platform download directory management

### Core Components

## 1. Project Structure

```
venaqui/
├── cmd/
│   └── venaqui/
│       └── main.go                 # Entry point
├── internal/
│   ├── realdebrid/
│   │   ├── client.go              # RD API client
│   │   ├── auth.go                # Authentication handling
│   │   └── unrestrict.go          # Link unrestriction
│   ├── aria2/
│   │   ├── client.go              # aria2 RPC client
│   │   ├── download.go            # Download management
│   │   └── status.go              # Status tracking
│   ├── tui/
│   │   ├── model.go               # Bubble Tea model
│   │   ├── update.go              # Update logic
│   │   ├── view.go                # Render logic
│   │   └── styles.go              # UI styling
│   ├── config/
│   │   ├── config.go              # Configuration management
│   │   └── paths.go               # Path resolution
│   └── utils/
│       ├── downloads.go           # Download folder detection
│       └── validation.go          # Input validation
├── pkg/
│   └── models/
│       ├── download.go            # Download state model
│       └── types.go               # Shared types
├── go.mod
├── go.sum
└── README.md
```


## 2. Dependencies

```go
// go.mod
module github.com/yourusername/venaqui

go 1.22

require (
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/siku2/arigo v0.2.0           // aria2 RPC client
    github.com/spf13/cobra v1.8.0           // CLI framework
    github.com/spf13/viper v1.18.2          // Configuration
    github.com/dustin/go-humanize v1.0.1    // Human-readable sizes
)
```


## 3. Real-Debrid API Integration

### Authentication \& Configuration

```go
// internal/config/config.go
package config

import (
    "os"
    "path/filepath"
)

type Config struct {
    RealDebridAPIToken string
    Aria2RPCUrl        string
    Aria2Secret        string
    DefaultDownloadDir string
}

func Load() (*Config, error) {
    homeDir, _ := os.UserHomeDir()
    configPath := filepath.Join(homeDir, ".venaqui", "config.yaml")
    
    // Load from config file using viper
    // Return config with RD API token, aria2 settings
}

func GetDefaultDownloadDir() (string, error) {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        return "", err
    }
    
    // Cross-platform downloads folder
    // Windows: C:\Users\<user>\Downloads
    // macOS/Linux: /Users/<user>/Downloads or /home/<user>/Downloads
    return filepath.Join(homeDir, "Downloads"), nil
}
```


### Real-Debrid Client

```go
// internal/realdebrid/client.go
package realdebrid

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

const BaseURL = "https://api.real-debrid.com/rest/1.0"

type Client struct {
    apiToken   string
    httpClient *http.Client
}

func NewClient(apiToken string) *Client {
    return &Client{
        apiToken:   apiToken,
        httpClient: &http.Client{},
    }
}

// UnrestrictLink converts a hoster link to unrestricted download link
func (c *Client) UnrestrictLink(link string) (*UnrestrictedLink, error) {
    endpoint := fmt.Sprintf("%s/unrestrict/link", BaseURL)
    
    payload := map[string]string{
        "link": link,
    }
    
    jsonData, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
    req.Header.Set("Authorization", "Bearer "+c.apiToken)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("RD API error: %d", resp.StatusCode)
    }
    
    var result UnrestrictedLink
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}

type UnrestrictedLink struct {
    ID       string `json:"id"`
    Filename string `json:"filename"`
    MimeType string `json:"mimeType"`
    Filesize int64  `json:"filesize"`
    Link     string `json:"link"`      // Direct download link
    Host     string `json:"host"`
    Chunks   int    `json:"chunks"`
    Download string `json:"download"`
}
```


## 4. aria2 Integration

### aria2 Client Setup

```go
// internal/aria2/client.go
package aria2

import (
    "context"
    "github.com/siku2/arigo"
)

type Client struct {
    rpc    arigo.Client
    ctx    context.Context
}

func NewClient(rpcURL, secret string) (*Client, error) {
    // Default: http://localhost:6800/jsonrpc
    client, err := arigo.Dial(rpcURL, secret)
    if err != nil {
        return nil, err
    }
    
    return &Client{
        rpc: client,
        ctx: context.Background(),
    }, nil
}

func (c *Client) AddDownload(url, downloadDir string) (string, error) {
    options := arigo.Options{
        "dir": downloadDir,
        "max-connection-per-server": "16",
        "split": "16",
        "min-split-size": "1M",
    }
    
    gid, err := c.rpc.AddURI([]string{url}, options)
    if err != nil {
        return "", err
    }
    
    return gid, nil
}

func (c *Client) GetStatus(gid string) (*DownloadStatus, error) {
    status, err := c.rpc.TellStatus(gid, 
        "gid", "status", "totalLength", "completedLength",
        "downloadSpeed", "files", "errorMessage")
    if err != nil {
        return nil, err
    }
    
    return &DownloadStatus{
        GID:             status.Gid,
        Status:          status.Status,
        TotalLength:     status.TotalLength,
        CompletedLength: status.CompletedLength,
        DownloadSpeed:   status.DownloadSpeed,
        Files:           status.Files,
    }, nil
}

func (c *Client) Close() error {
    return c.rpc.Close()
}

type DownloadStatus struct {
    GID             string
    Status          string
    TotalLength     int64
    CompletedLength int64
    DownloadSpeed   int64
    Files           []arigo.FileData
}
```


## 5. Bubble Tea TUI Implementation

### Model Definition

```go
// internal/tui/model.go
package tui

import (
    "time"
    "github.com/charmbracelet/bubbletea"
    "github.com/yourrepo/venaqui/internal/aria2"
)

type Model struct {
    aria2Client   *aria2.Client
    gid           string
    status        *aria2.DownloadStatus
    filename      string
    err           error
    quitting      bool
    lastUpdate    time.Time
}

type tickMsg time.Time
type statusMsg *aria2.DownloadStatus
type errMsg error

func InitialModel(aria2Client *aria2.Client, gid, filename string) Model {
    return Model{
        aria2Client: aria2Client,
        gid:         gid,
        filename:    filename,
        lastUpdate:  time.Now(),
    }
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(
        tickCmd(),
        m.fetchStatus,
    )
}

func tickCmd() tea.Cmd {
    return tea.Tick(time.Second, func(t time.Time) tea.Msg {
        return tickMsg(t)
    })
}

func (m Model) fetchStatus() tea.Msg {
    status, err := m.aria2Client.GetStatus(m.gid)
    if err != nil {
        return errMsg(err)
    }
    return statusMsg(status)
}
```


### Update Logic

```go
// internal/tui/update.go
package tui

import (
    "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            m.quitting = true
            return m, tea.Quit
        }
    
    case tickMsg:
        return m, tea.Batch(
            tickCmd(),
            m.fetchStatus,
        )
    
    case statusMsg:
        m.status = msg
        
        // Check if download is complete
        if m.status.Status == "complete" {
            m.quitting = true
            return m, tea.Quit
        }
        
        // Check for errors
        if m.status.Status == "error" {
            m.err = fmt.Errorf("download error")
            return m, tea.Quit
        }
        
        return m, nil
    
    case errMsg:
        m.err = msg
        return m, tea.Quit
    }
    
    return m, nil
}
```


### View Rendering

```go
// internal/tui/view.go
package tui

import (
    "fmt"
    "github.com/charmbracelet/lipgloss"
    "github.com/dustin/go-humanize"
)

func (m Model) View() string {
    if m.err != nil {
        return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
    }
    
    if m.quitting {
        if m.status != nil && m.status.Status == "complete" {
            return successStyle.Render("✓ Download complete!\n")
        }
        return "Exiting...\n"
    }
    
    if m.status == nil {
        return "Initializing download...\n"
    }
    
    // Calculate progress
    var progress float64
    if m.status.TotalLength > 0 {
        progress = float64(m.status.CompletedLength) / float64(m.status.TotalLength) * 100
    }
    
    // Format sizes
    completed := humanize.Bytes(uint64(m.status.CompletedLength))
    total := humanize.Bytes(uint64(m.status.TotalLength))
    speed := humanize.Bytes(uint64(m.status.DownloadSpeed))
    
    // Build view
    s := titleStyle.Render("Venaqui - Download Manager") + "\n\n"
    s += fmt.Sprintf("File: %s\n", filenameStyle.Render(m.filename))
    s += fmt.Sprintf("Status: %s\n", statusStyle.Render(m.status.Status))
    s += fmt.Sprintf("Progress: %s\n", progressBar(progress))
    s += fmt.Sprintf("Size: %s / %s (%.2f%%)\n", completed, total, progress)
    s += fmt.Sprintf("Speed: %s/s\n", speed)
    s += "\n" + helpStyle.Render("Press 'q' to quit")
    
    return s
}

func progressBar(percent float64) string {
    width := 50
    filled := int(percent / 100 * float64(width))
    
    bar := ""
    for i := 0; i < width; i++ {
        if i < filled {
            bar += "█"
        } else {
            bar += "░"
        }
    }
    
    return progressStyle.Render(bar)
}
```


### Styles

```go
// internal/tui/styles.go
package tui

import (
    "github.com/charmbracelet/lipgloss"
)

var (
    titleStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#00D9FF")).
        MarginBottom(1)
    
    filenameStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#FFEB3B"))
    
    statusStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#4CAF50"))
    
    progressStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#00D9FF"))
    
    successStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#4CAF50"))
    
    errorStyle = lipgloss.NewStyle().
        Bold(true).
        Foreground(lipgloss.Color("#F44336"))
    
    helpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#757575"))
)
```


## 6. Main Application

```go
// cmd/venaqui/main.go
package main

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    
    tea "github.com/charmbracelet/bubbletea"
    "github.com/spf13/cobra"
    "github.com/yourrepo/venaqui/internal/aria2"
    "github.com/yourrepo/venaqui/internal/config"
    "github.com/yourrepo/venaqui/internal/realdebrid"
    "github.com/yourrepo/venaqui/internal/tui"
)

var rootCmd = &cobra.Command{
    Use:   "venaqui [link] [location]",
    Short: "Download files via Real-Debrid and aria2",
    Args:  cobra.MinimumNArgs(1),
    Run:   run,
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}

func run(cmd *cobra.Command, args []string) {
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
        os.Exit(1)
    }
    
    // Parse arguments
    link := args[^0]
    downloadDir := cfg.DefaultDownloadDir
    if len(args) > 1 {
        downloadDir = args[^1]
    }
    
    // Ensure download directory exists
    os.MkdirAll(downloadDir, 0755)
    
    // Start aria2c if not running
    ensureAria2Running()
    
    // Initialize Real-Debrid client
    rdClient := realdebrid.NewClient(cfg.RealDebridAPIToken)
    
    // Unrestrict link
    fmt.Println("Unrestricting link via Real-Debrid...")
    unrestrictedLink, err := rdClient.UnrestrictLink(link)
    if err != nil {
        fmt.Fprintf(os.Stderr, "RD API error: %v\n", err)
        os.Exit(1)
    }
    
    // Initialize aria2 client
    aria2Client, err := aria2.NewClient(cfg.Aria2RPCUrl, cfg.Aria2Secret)
    if err != nil {
        fmt.Fprintf(os.Stderr, "aria2 connection error: %v\n", err)
        os.Exit(1)
    }
    defer aria2Client.Close()
    
    // Add download to aria2
    fmt.Println("Starting download...")
    gid, err := aria2Client.AddDownload(unrestrictedLink.Link, downloadDir)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Download error: %v\n", err)
        os.Exit(1)
    }
    
    // Start TUI
    model := tui.InitialModel(aria2Client, gid, unrestrictedLink.Filename)
    p := tea.NewProgram(model)
    
    if _, err := p.Run(); err != nil {
        fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
        os.Exit(1)
    }
}

func ensureAria2Running() {
    // Check if aria2c is running on RPC port
    // If not, start it with: aria2c --enable-rpc --rpc-listen-all
    cmd := exec.Command("aria2c", 
        "--enable-rpc", 
        "--rpc-listen-all",
        "--daemon=true",
        "--max-connection-per-server=16",
        "--split=16",
        "--min-split-size=1M")
    cmd.Start()
}
```


## 7. Configuration File

```yaml
# ~/.venaqui/config.yaml
realdebrid:
  api_token: "YOUR_RD_API_TOKEN"

aria2:
  rpc_url: "http://localhost:6800/jsonrpc"
  secret: ""  # Optional RPC secret

download:
  default_dir: ""  # Leave empty for OS default Downloads folder
```


## 8. Building and Running

### Prerequisites

1. **Install aria2**:
    - macOS: `brew install aria2`
    - Linux: `sudo apt install aria2` or `sudo yum install aria2`
    - Windows: Download from https://aria2.github.io/
2. **Get Real-Debrid API Token**: Visit https://real-debrid.com/apitoken

### Build Commands

```bash
# Initialize module
go mod init github.com/yourusername/venaqui
go mod tidy

# Build
go build -o venaqui cmd/venaqui/main.go

# Install
go install cmd/venaqui/main.go
```


### Usage Examples

```bash
# Download to default location
venaqui "https://mega.nz/file/example"

# Download to specific location
venaqui "https://mega.nz/file/example" "/path/to/download"

# With full path
venaqui "https://1fichier.com/example" "$HOME/Downloads/Movies"
```


## 9. Key Features to Implement

### Phase 1 (MVP)

- [x] Real-Debrid link unrestriction[^5][^1]
- [x] aria2 RPC integration[^4][^6]
- [x] Basic TUI with progress display[^2][^3]
- [x] Cross-platform download folder detection[^7][^8]


### Phase 2 (Enhancements)

- [ ] Configuration wizard on first run
- [ ] Multiple simultaneous downloads
- [ ] Download queue management
- [ ] Pause/resume functionality
- [ ] Download history
- [ ] Bandwidth limiting


### Phase 3 (Advanced)

- [ ] Torrent magnet link support via RD
- [ ] Automatic file verification
- [ ] Notifications on completion
- [ ] Dark/light theme switching
- [ ] Plugin system for other premium link services


## 10. Testing Strategy

```go
// Example test structure
package aria2_test

import (
    "testing"
    "github.com/yourrepo/venaqui/internal/aria2"
)

func TestAddDownload(t *testing.T) {
    client, _ := aria2.NewClient("http://localhost:6800/jsonrpc", "")
    gid, err := client.AddDownload("http://example.com/file.zip", "/tmp")
    
    if err != nil {
        t.Fatalf("Failed to add download: %v", err)
    }
    
    if gid == "" {
        t.Error("Expected non-empty GID")
    }
}
```


## 11. Error Handling

Key error scenarios to handle:

- Invalid Real-Debrid API token
- Link not supported by Real-Debrid
- aria2 not running or unreachable
- Network interruptions during download
- Insufficient disk space
- Invalid download paths
- API rate limiting


## 12. Performance Considerations

- aria2 uses 16 connections per server by default for maximum speed[^9]
- Real-Debrid API has rate limits - implement exponential backoff
- Use WebSocket connection to aria2 for real-time updates instead of polling[^4]
- Implement connection pooling for multiple downloads


[^1]: https://api.real-debrid.com

[^2]: https://www.youtube.com/watch?v=ERaZi0YvBRs

[^3]: https://themarkokovacevic.com/posts/terminal-ui-with-bubbletea/

[^4]: https://pkg.go.dev/github.com/siku2/arigo

[^5]: https://valentingot.github.io/real-debrid/available_requests/unrestrict.html

[^6]: https://github.com/kahosan/aria2-rpc

[^7]: https://www.reddit.com/r/learnpython/comments/p0jgki/crossplatform_way_to_fetch_default_download/

[^8]: https://stackoverflow.com/questions/77934569/finding-downloads-folder-programmatically-go

[^9]: https://blog.csdn.net/qq_21484461/article/details/138390681

[^10]: https://github.com/ValentinGot/real-debrid

[^11]: https://valentingot.github.io/real-debrid/

[^12]: https://valentingot.github.io/real-debrid/basic_usage/index.html

[^13]: https://www.reddit.com/r/RealDebrid/comments/h9tei2/real_debrid_api_access/

[^14]: https://troypoint.com/real-debrid/

[^15]: https://github.com/siku2/arigo

[^16]: https://pypi.org/project/aria2-rpc-client/

[^17]: https://pkg.go.dev/github.com/saeidrp/aria2-rpc

[^18]: https://pkg.go.dev/github.com/myanimestream/arias/aria2

[^19]: https://github.com/5aitama/RealDebrid4Node

