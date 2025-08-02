# file-converter

A minimal CLI tool for file format conversion written in Go.

## Features

- PPTX to PDF conversion

## Pre-requisites

### Linux/macOS

```bash
# Ubuntu/Debian
sudo apt-get install libreoffice

# Fedora/RHEL/CentOS
sudo dnf install libreoffice

# macOS
brew install --cask libreoffice
```

## Installation

```bash
# Cline the repository
git clone https://github.com/gr1m0h/file-converter.git
cd file-converter

# Build the binary
make build

# Install the binary
make install
```

## Usage

```bash
# Convert a PPTX file to PDF
fconv -i input.pptx

# Specify an output file
fconv -i input.pptx -o output.pdf

# Using make run
make run ARGS="-i input.pptx -o output.pdf"
```

## License

- [MIT License](LICENSE)

