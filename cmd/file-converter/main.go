package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func main() {
	var (
		input  string
		output string
	)

	flag.StringVar(&input, "i", "", "Input file path")
	flag.StringVar(&output, "o", "", "Output file path")
	flag.Parse()

	if input == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -i input file  [-o output file]\n", os.Args[0])
		os.Exit(1)
	}

	if output == "" {
		output = strings.TrimSuffix(input, filepath.Ext(input)) + ".pdf"
	}

	if err := convertPPTXtpPDF(input, output); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", input, output)
}

func convertPPTXtpPDF(inputPath, outputPath string) error {
	absInput, err := filepath.Abs(inputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve input path: %w", err)
	}

	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	if runtime.GOOS == "darwin" || runtime.GOOS == "linux" {
		return convertWithLibreOffice(absInput, absOutput)
	} else { // For Windows, we can use unioffice library directly
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func convertWithLibreOffice(inputPath, outputPath string) error {
	cmd := exec.Command("whitch", "libreoffice")
	if err := cmd.Run(); err != nil {
		cmd = exec.Command("which", "soffice")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("LibreOffice not found: %w", err)
		}
	}

	outputDir := filepath.Dir(outputPath)

	cmdArgs := []string{
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outputDir,
		inputPath,
	}

	cmd = exec.Command("libreoffice", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		cmd = exec.Command("soffice", cmdArgs...)
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("conversion failed: %w\nOutput: %s", err, string(output))
		}
	}

	generatedFile := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))+".pdf")
	if generatedFile != outputPath {
		if err := os.Rename(generatedFile, outputPath); err != nil {
			return fmt.Errorf("failed to rename output file: %w", err)
		}
	}

	return nil
}
