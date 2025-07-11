package scan

import (
	"context"
	"testing"
	"time"
)

func TestEngine(t *testing.T) {
	engine := NewEngine(nil)

	target := ScanTarget{
		Type: TargetTypeDirectory,
		Path: ".",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := engine.ScanTarget(ctx, target)
	if err != nil {
		t.Fatalf("ScanTarget failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	scanners := engine.GetAvailableScanners()
	if len(scanners) == 0 {
		t.Error("Expected some available scanners, got none")
	}

	t.Logf("Available scanners: %v", scanners)
	t.Logf("Scan completed successfully in %v", result.Duration)
}
