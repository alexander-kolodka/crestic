package tests

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/alexander-kolodka/crestic/tests/harness"
)

func TestMain(m *testing.M) {
	assertRestic()
	assertCrestic()

	os.Exit(m.Run())
}

func assertRestic() {
	_, err := exec.LookPath("restic")
	if err == nil {
		return
	}

	_, _ = fmt.Fprintln(
		os.Stderr,
		"restic is required to run integration tests.\n"+
			"Please install it and ensure it is available in PATH.",
	)
	os.Exit(1)
}

func assertCrestic() {
	_, err := os.Stat(harness.CresticBin())
	if err == nil {
		return
	}

	if os.IsNotExist(err) {
		_, _ = fmt.Fprintln(
			os.Stderr,
			"crestic binary is required for integration tests.\n"+
				"Build it first, e.g. `go build -o bin/crestic .`",
		)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "failed to access bin/crestic: %v\n", err)
	os.Exit(1)
}
