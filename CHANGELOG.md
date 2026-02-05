# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.2.2] - 2026-02-05

### Added
- **Completion Statistics**: Enhanced completion screen now displays detailed download statistics
  - Total download time
  - Average download speed
  - Peak download speed
  - File size
  - Number of connections used
  - Completion status indicator

### Removed
- **Upload Speed Field**: Removed upload speed display from active download statistics (not applicable for download-only operations)

### Changed
- Improved completion screen layout with organized statistics box
- Better visual presentation of download performance metrics

## [1.2.1] - 2026-02-05

### Changed
- **Updated Color Palette**: Changed primary color to `#EF4444` (red) with harmonized color scheme
  - Secondary accent: `#F87171` (lighter red/pink)
  - Success: `#10B981` (emerald green)
  - Warning: `#F59E0B` (amber)
  - Error: `#DC2626` (darker red)
  - Improved text and border colors for better contrast
- **Branding Update**: Changed "Venaqui" to "venaqui" (lowercase 'v') throughout the application
- **Documentation**: Added explanation about the name origin (Spanish: "ven aquí" meaning "come here")

## [1.2.0] - 2026-02-05

### Added
- **Enhanced TUI Interface**: Completely redesigned terminal user interface with modern styling
  - Beautiful color-coded status indicators with icons (● Active, ✓ Complete, ✗ Error)
  - Boxed sections with rounded borders for better organization
  - Two-column statistics layout for improved readability
- **Advanced Statistics Display**:
  - Real-time download speed with highlighting
  - Upload speed tracking
  - Active connections count
  - Elapsed time display
  - ETA (Estimated Time to Arrival) calculation
  - Remaining bytes display
- **Speed History Graph**: Visual ASCII graph showing download speed trends over time (last 50 samples)
- **Enhanced Progress Visualization**: Improved progress bar with gradient effect
- **Post-Download Actions**:
  - Press `o` to open downloaded file with default application
  - Press `d` or `s` to reveal file in Finder/Explorer (highlights the file)
  - TUI no longer auto-quits on completion, allowing user to choose action
- **Cross-Platform File Opening**: Support for opening files and directories on macOS, Linux, and Windows

### Changed
- TUI now stays open after download completion, allowing users to interact with the downloaded file
- Improved error handling and display
- Better layout and visual hierarchy in the TUI

### Technical Improvements
- Enhanced aria2 status tracking with additional fields (upload speed, connections, pieces)
- Added `GetFilePath()` and `GetFileDirectory()` methods to DownloadStatus
- Added `GetETA()` method for time estimation
- Improved cross-platform file system operations

## [1.1.8] - Previous Release

Previous releases and changes...
