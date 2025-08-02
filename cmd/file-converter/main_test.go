package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestConvertPPTXtoPDF(t *testing.T) {
	// Skip tests on unsupported platforms
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("Test only runs on macOS and Linux")
	}

	// Create a temporary test file
	tmpDir := t.TempDir()
	testPPTX := filepath.Join(tmpDir, "test.pptx")
	if err := os.WriteFile(testPPTX, []byte("dummy pptx content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		inputPath   string
		outputPath  string
		wantErr     bool
		errContains string
	}{
		{
			name:        "non-existent input file",
			inputPath:   "/tmp/non-existent.pptx",
			outputPath:  "/tmp/output.pdf",
			wantErr:     true,
			errContains: "input file does not exist",
		},
		{
			name:        "empty input path",
			inputPath:   "",
			outputPath:  "/tmp/output.pdf",
			wantErr:     true,
			errContains: "input path cannot be empty",
		},
		{
			name:        "empty output path",
			inputPath:   testPPTX,
			outputPath:  "",
			wantErr:     true,
			errContains: "output path cannot be empty",
		},
		{
			name:        "valid input file with dummy content",
			inputPath:   testPPTX,
			outputPath:  filepath.Join(tmpDir, "output.pdf"),
			wantErr:     false, // LibreOffice might convert even dummy content
			errContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertPPTXtoPDF(tt.inputPath, tt.outputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertPPTXtoPDF() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !containsString(err.Error(), tt.errContains) {
					t.Errorf("convertPPTXtoPDF() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

func TestConvertWithLibreOffice(t *testing.T) {
	// Skip tests on unsupported platforms
	if runtime.GOOS != "darwin" && runtime.GOOS != "linux" {
		t.Skip("Test only runs on macOS and Linux")
	}

	// Create a temporary test file
	tmpDir := t.TempDir()
	testInput := filepath.Join(tmpDir, "test.pptx")
	testOutput := filepath.Join(tmpDir, "test.pdf")

	// Create a dummy PPTX file for testing
	if err := os.WriteFile(testInput, []byte("dummy pptx content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name        string
		inputPath   string
		outputPath  string
		wantErr     bool
		errContains string
	}{
		{
			name:        "valid paths but no real PPTX content",
			inputPath:   testInput,
			outputPath:  testOutput,
			wantErr:     false, // LibreOffice might succeed even with dummy content
			errContains: "",
		},
		{
			name:        "non-existent input",
			inputPath:   "/tmp/non-existent.pptx",
			outputPath:  testOutput,
			wantErr:     true,
			errContains: "conversion failed: output file was not created",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := convertWithLibreOffice(tt.inputPath, tt.outputPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertWithLibreOffice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errContains != "" {
				if err == nil || !containsString(err.Error(), tt.errContains) {
					t.Errorf("convertWithLibreOffice() error = %v, want error containing %q", err, tt.errContains)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
