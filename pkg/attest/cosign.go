package attest

import "os/exec"

// Check if cosign is installed on machine
func checkCosign() error {
	cmd := exec.Command("cosign")
	return cmd.Run()
}
