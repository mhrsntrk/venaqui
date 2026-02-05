package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/mhrsntrk/venaqui/internal/aria2"
	"github.com/mhrsntrk/venaqui/internal/config"
	"github.com/mhrsntrk/venaqui/internal/realdebrid"
	"github.com/mhrsntrk/venaqui/internal/tui"
	"github.com/mhrsntrk/venaqui/internal/utils"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "venaqui [link] [location]",
	Short: "Download files via Real-Debrid and aria2",
	Long: `Venaqui is a command-line tool with a Terminal User Interface (TUI) that
leverages Real-Debrid premium links and aria2 for high-speed downloads.`,
	Args: cobra.MinimumNArgs(1),
	Run:  run,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("venaqui version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built: %s\n", date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
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
		fmt.Fprintf(os.Stderr, "Please create a config file at ~/.venaqui/config.yaml\n")
		os.Exit(1)
	}

	// Parse arguments
	link := args[0]
	downloadDir := cfg.DefaultDownloadDir
	if len(args) > 1 {
		downloadDir = args[1]
	}

	// Validate URL
	if err := utils.ValidateURL(link); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid URL: %v\n", err)
		os.Exit(1)
	}

	// Normalize and validate download directory
	downloadDir = utils.NormalizePath(downloadDir)
	if err := utils.ValidatePath(downloadDir); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid download directory: %v\n", err)
		os.Exit(1)
	}

	// Ensure download directory exists
	if err := utils.EnsureDirExists(downloadDir); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create download directory: %v\n", err)
		os.Exit(1)
	}

	// Start aria2c if not running
	if err := ensureAria2Running(cfg.Aria2RPCUrl); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start aria2: %v\n", err)
		fmt.Fprintf(os.Stderr, "Please ensure aria2 is installed and accessible\n")
		os.Exit(1)
	}

	// Initialize Real-Debrid client
	rdClient := realdebrid.NewClient(cfg.RealDebridAPIToken)

	// Validate token (optional check)
	if err := rdClient.ValidateToken(); err != nil {
		fmt.Fprintf(os.Stderr, "Real-Debrid API token validation failed: %v\n", err)
		os.Exit(1)
	}

	var unrestrictedLink *realdebrid.UnrestrictedLink
	var filename string

	// Check if this is a torrent or magnet link
	if utils.IsTorrentLink(link) || utils.IsMagnetLink(link) {
		// Handle torrent/magnet link
		fmt.Println("Adding torrent to Real-Debrid...")
		var torrentResp *realdebrid.AddTorrentResponse
		var err error

		if utils.IsMagnetLink(link) {
			torrentResp, err = rdClient.AddMagnet(link)
		} else {
			torrentResp, err = rdClient.AddTorrent(link)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "RD API error: %v\n", err)
			os.Exit(1)
		}

		// Check if files need to be selected
		torrentInfo, err := rdClient.GetTorrentInfo(torrentResp.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "RD API error: %v\n", err)
			os.Exit(1)
		}

		// Select all files if needed
		if torrentInfo.Status == "waiting_files_selection" {
			fileIDs := []int{}
			for _, file := range torrentInfo.Files {
				fileIDs = append(fileIDs, file.ID)
			}
			if len(fileIDs) > 0 {
				fmt.Println("Selecting files...")
				if err := rdClient.SelectFiles(torrentResp.ID, fileIDs); err != nil {
					fmt.Fprintf(os.Stderr, "RD API error: %v\n", err)
					os.Exit(1)
				}
			}
		}

		fmt.Println("Waiting for torrent to be processed...")
		torrentInfo, err = rdClient.WaitForTorrentReady(torrentResp.ID, 5*time.Minute)
		if err != nil {
			fmt.Fprintf(os.Stderr, "RD API error: %v\n", err)
			os.Exit(1)
		}

		if len(torrentInfo.Links) == 0 {
			fmt.Fprintf(os.Stderr, "No download links available from torrent\n")
			os.Exit(1)
		}

		// Use the first link (or we could unrestrict all links)
		// For now, we'll unrestrict the first link
		fmt.Println("Unrestricting download link...")
		unrestrictedLink, err = rdClient.UnrestrictLink(torrentInfo.Links[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "RD API error: %v\n", err)
			os.Exit(1)
		}

		filename = torrentInfo.Filename
		if filename == "" {
			filename = unrestrictedLink.Filename
		}
	} else {
		// Handle regular hoster link
		fmt.Println("Unrestricting link via Real-Debrid...")
		var err error
		unrestrictedLink, err = rdClient.UnrestrictLink(link)
		if err != nil {
			fmt.Fprintf(os.Stderr, "RD API error: %v\n", err)
			os.Exit(1)
		}
		filename = unrestrictedLink.Filename
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
	if filename == "" {
		filename = unrestrictedLink.Filename
		if filename == "" {
			filename = filepath.Base(unrestrictedLink.Link)
		}
	}

	model := tui.InitialModel(aria2Client, gid, filename)
	p := tea.NewProgram(model)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}

// ensureAria2Running checks if aria2 is running and starts it if needed
func ensureAria2Running(rpcURL string) error {
	// Try to connect to aria2 RPC
	client, err := aria2.NewClient(rpcURL, "")
	if err == nil {
		// Connection successful, check if it's responsive
		if err := client.Ping(); err == nil {
			client.Close()
			return nil
		}
		client.Close()
	}

	// aria2 is not running, try to start it
	fmt.Println("Starting aria2 daemon...")

	cmd := exec.Command("aria2c",
		"--enable-rpc",
		"--rpc-listen-all",
		"--daemon=true",
		"--max-connection-per-server=16",
		"--split=16",
		"--min-split-size=1M",
		"--rpc-allow-origin-all",
	)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start aria2: %w", err)
	}

	// Wait a bit for aria2 to start
	time.Sleep(2 * time.Second)

	// Try to connect again
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		client, err := aria2.NewClient(rpcURL, "")
		if err == nil {
			if err := client.Ping(); err == nil {
				client.Close()
				return nil
			}
			client.Close()
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("aria2 failed to start or is not accessible")
}

// checkPort checks if a port is listening
func checkPort(host, port string) bool {
	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	if conn != nil {
		conn.Close()
		return true
	}
	return false
}
