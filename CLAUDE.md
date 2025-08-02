# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build and Run
- `make build` - Build the binary to ./bin/fconv
- `make run ARGS="-i input.pptx -o output.pdf"` - Build and run with arguments
- `make install` - Build and install binary to /usr/local/bin/

### Development
- `make fmt` - Format Go code
- `make lint` - Run golangci-lint (requires golangci-lint to be installed)
- `make test` - Run all tests
- `make deps` - Download and tidy dependencies
- `make clean` - Clean build artifacts and Go cache

### Direct Binary Usage
```bash
fconv -i input.pptx              # Convert to PDF with default output name
fconv -i input.pptx -o output.pdf # Convert with specific output name
```

## Architecture

This is a minimal Go CLI tool for file format conversion, currently supporting PPTX to PDF conversion using LibreOffice as the backend converter.

### Key Components

1. **Entry Point**: `cmd/file-converter/main.go` - CLI interface with flag parsing
2. **Conversion Logic**: The main.go file contains:
   - `convertPPTXtpPDF()` - Routes to appropriate converter based on OS
   - `convertWithLibreOffice()` - Uses LibreOffice headless mode for conversion

### Important Implementation Details

- The tool requires LibreOffice to be installed on the system
- Only macOS and Linux are currently supported (Windows returns an error)
- The converter uses LibreOffice in headless mode with `--convert-to pdf`
- Output files are renamed if LibreOffice generates a different name than requested
- The tool checks for both `libreoffice` and `soffice` commands (line 59 has a typo: "whitch" should be "which")

### Build Configuration
- Binary name: `fconv`
- Build flags: `-ldflags="-s -w"` (strip debug info and symbol table)
- Go version: 1.24.1