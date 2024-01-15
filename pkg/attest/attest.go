package attest

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/types"
)

func Run(image string, opts types.AttestationOptions) {

	// Verify if cosign is installed
	err := checkCosign()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			log.Error("Unable to run cosign. Make sure it is installed first.")
		} else {
			log.Error(err)
		}
		return
	}

	log.Print("Running sbom attestation...")
	found, err := helper.IsFileExists(opts.Predicate)
	if err != nil {
		log.Error(err)
		return
	}

	if !found {
		log.Error("Predicate file not found.")
		return
	}
	log.Print("Attestation Complete!")
	err = attest(image, opts)
	if err != nil {
		log.Error(err)
		return
	}

	log.Print("Running verification...")
	err = verify(image, opts)
	if err != nil {
		log.Errorf("Error occurred when verifying attestation. Please make sure that the paths or fields specified are correct. \n%s", err.Error())
		return
	}
	log.Print("Verification Done!")
}

// Attest SBOM
func attest(image string, opts types.AttestationOptions) error {
	args := fmt.Sprintf("attest --yes --key %+v --type %+v --predicate %+v %+v",
		opts.Key, opts.AttestType, opts.Predicate, image)
	attest := strings.Split(args, " ")

	out := new(strings.Builder)
	cmd := exec.Command("cosign", attest...)
	cmd.Stdin = strings.NewReader(opts.Password)
	cmd.Stdout = out

	err := cmd.Run()
	if err != nil {
		return err
	}

	log.Print(out.String())

	return nil
}

func verify(image string, opts types.AttestationOptions) error {
	args := fmt.Sprintf("verify-attestation --key %+v --type %+v %+v", opts.Pub, opts.AttestType, image)
	verify := strings.Split(args, " ")
	if len(opts.OutputFile) != 0 {
		file := fmt.Sprintf("--output-file %+v", helper.AddFileExtension(opts.OutputFile, "json"))
		verify = append(verify, strings.Split(file, " ")...)
	}
	out := new(strings.Builder)
	cmd := exec.Command("cosign", verify...)
	cmd.Stdout = out

	err := cmd.Run()
	if err != nil {
		return err
	}

	log.Print(out.String())

	return nil
}
