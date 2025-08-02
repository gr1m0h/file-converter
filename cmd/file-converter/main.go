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
		format string
	)

	flag.StringVar(&input, "i", "", "Input file path")
	flag.StringVar(&output, "o", "", "Output file path")
	flag.StringVar(&format, "f", "", "Output format(jpeg, jpg, png, pdf) - auto-detected if not spe")
	flag.Parse()

	if input == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -i input_file [-o output_file] [-f format]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nSupported conversions:\n")
		fmt.Fprintf(os.Stderr, "  - PPTX to PDF\n")
		os.Exit(1)
	}

	if _, err := os.Stat(input); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Input file does not exist: %v\n", err)
		os.Exit(1)
	}

	inputExt := strings.ToLower(filepath.Ext(input))

	switch inputExt {
	case ".pptx":
		if output == "" {
			output = strings.TrimSuffix(input, filepath.Ext(input)) + ".pdf"
		}
		if err := convertPPTXtoPDF(input, output); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unsupported input file format: %s\n", inputExt)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", input, output)
}

func convertPPTXtoPDF(inputPath, outputPath string) error {
	if inputPath == "" {
		return fmt.Errorf("input path cannot be empty")
	}

	if outputPath == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Check if input file exists
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input file does not exist: %w", err)
	}

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
	cmd := exec.Command("which", "libreoffice")
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

	// Check if the generated file exists
	if _, err := os.Stat(generatedFile); err != nil {
		return fmt.Errorf("conversion failed: output file was not created")
	}

	if generatedFile != outputPath {
		if err := os.Rename(generatedFile, outputPath); err != nil {
			return fmt.Errorf("failed to rename output file: %w", err)
		}
	}

	return nil
}
