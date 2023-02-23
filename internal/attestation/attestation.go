package attestation

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	sbom "github.com/carbonetes/diggity/internal"
	"github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/model"
	"github.com/carbonetes/diggity/internal/ui"
	"github.com/google/uuid"
)

var (
	log        = logger.GetLogger()
	cosign     = "cosign"
	sbomPrefix = "diggity-sbom-"

	// Arguments
	Arguments model.Arguments
)

// Run SBOM Attestation
func Attest(image string, attestationOptions *model.AttestationOptions) {
	var predicate string

	// Verify if cosign is installed
	checkCosign()

	// Generate SBOM as needed
	if *attestationOptions.Predicate == "" {
		predicate = generateBom(image, attestationOptions.BomArgs, *attestationOptions.OutputType)
	} else {
		predicate = *attestationOptions.Predicate
	}

	// Attest specified BOM file
	attestSpinner := ui.InitSpinner("Attesting SBOM...")
	go ui.RunSpinner(attestSpinner)
	attestBom(image, predicate, attestationOptions)
	ui.DoneSpinner(attestSpinner)

	// Get Attestation
	verifySpinner := ui.InitSpinner("Verifying Attestation...")
	go ui.RunSpinner(verifySpinner)
	getAttestation(image, attestationOptions)
	ui.DoneSpinner(verifySpinner)
}

// Check if cosign is installed on machine
func checkCosign() {
	cmd := exec.Command("cosign")
	err := cmd.Run()

	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			log.Print("[warning]: Unable to run cosign. Make sure it is installed first.")
			os.Exit(1)
		}
		log.Fatal(err.Error())
	}
}

// Attest SBOM
func attestBom(image string, predicate string, attestationOptions *model.AttestationOptions) {
	args := fmt.Sprintf("attest --key %+v --type %+v --predicate %+v %+v",
		*attestationOptions.Key, *attestationOptions.AttestType, predicate, image)
	attest := strings.Split(args, " ")

	cmd := exec.Command(cosign, attest...)
	cmd.Stdin = strings.NewReader(*attestationOptions.Password)
	cmd.Stdout = os.Stdout
	err := cmd.Run()

	if err != nil {
		log.Fatal("[error]: Error occurred when running SBOM attestation. Please make sure that the paths or fields specified are correct.")
		os.Exit(1)
	}
}

// Get Attestation by Verifying
func getAttestation(image string, attestationOptions *model.AttestationOptions) {
	args := fmt.Sprintf("verify-attestation --key %+v --type %+v %+v", *attestationOptions.Pub, *attestationOptions.AttestType, image)
	verify := strings.Split(args, " ")

	if *attestationOptions.OutputFile != "" {
		fileArg := fmt.Sprintf("--output-file %+v", *attestationOptions.OutputFile)
		verify = append(verify, strings.Split(fileArg, " ")...)
	}

	cmd := exec.Command(cosign, verify...)
	if *attestationOptions.OutputFile == "" {
		cmd.Stdout = os.Stdout
	}
	err := cmd.Run()

	if err != nil {
		log.Fatal("[error]: Error occurred when verifying attestation. Please make sure that the paths or fields specified are correct.")
		os.Exit(1)
	}
}

// Generate SBOM
func generateBom(image string, arguments *model.Arguments, outputType string) string {
	// Generate Temp Bom Filename
	var bomFileName string
	switch outputType {
	case "json", "cyclonedx-json", "spdx-json", "cyclonedxjson", "spdxjson":
		bomFileName = sbomPrefix + uuid.NewString() + ".json"
	case "cyclonedx", "cyclonedx-xml", "cyclonedxxml", "cdx":
		bomFileName = sbomPrefix + uuid.NewString() + ".cdx"
	case "spdx-tag-value", "spdxtagvalue", "spdxtv", "spdx":
		bomFileName = sbomPrefix + uuid.NewString() + ".spdx"
	default:
		bomFileName = sbomPrefix + uuid.NewString() + ".json"
	}

	// Init Args
	bomPath := filepath.Join(".", bomFileName)
	Arguments = *arguments
	Arguments.Image = &image
	Arguments.OutputFile = &bomPath
	Arguments.Output = (*model.Output)(&outputType)

	// Start SBOM
	sbom.Start(&Arguments)

	return bomPath
}
