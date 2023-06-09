package ui

import "fmt"

func OnCheckingImageFromLocal(image string) {
	if disabled {
		return
	}
	pb.Describe(fmt.Sprintf("Checking image %s from local", image))
	go run()
}

func OnPullingPublicImage(image string) {
	if disabled {
		return
	}
	pb.Describe(fmt.Sprintf("Pulling public image %s", image))
	go run()
}

func OnPullingImageFromRegistry(image string) {
	if disabled {
		return
	}
	pb.Describe(fmt.Sprintf("Pulling image %s from privae registry", image))
	go run()
}

func OnExtractingImage(image string) {
	if disabled {
		return
	}
	pb.Describe(fmt.Sprintf("Extracting image %s", image))
	go run()
}

func OnScanningImage(image string) {
	if disabled {
		return
	}
	pb.Describe(fmt.Sprintf("Scanning image %s", image))
	go run()
}

func OnScanningTar(tar string) {
	if disabled {
		return
	}
	pb.Describe(fmt.Sprintf("Scanning tar file %s", tar))
	go run()
}

func OnScanningDir(dir string) {
	if disabled {
		return
	}
	pb.Describe(fmt.Sprintf("Scanning directory %s", dir))
	go run()
}

func OnSbomAttestation() {
	if disabled {
		return
	}
	pb.Describe("Attesting sbom")
	go run()
}

func OnVerifyingAttestation() {
	if disabled {
		return
	}
	pb.Describe("Verifying attestation")
	go run()
}
