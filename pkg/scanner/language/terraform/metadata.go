package terraform

import (
	"github.com/carbonetes/diggity/internal/log"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type TerraformLockFile struct {
	Providers []Provider `hcl:"provider,block"`
}

type Provider struct {
	URL         string   `hcl:",label" json:"url"`
	Constraints string   `hcl:"constraints" json:"constraints"`
	Version     string   `hcl:"version" json:"version"`
	Hashes      []string `hcl:"hashes" json:"hashes"`
}

func readLockfile(content []byte) (*TerraformLockFile, error) {
	parser := hclparse.NewParser()
	file, diag := parser.ParseHCL(content, "lockfile.hcl")
	if diag.HasErrors() {
		log.Debugf("Error parsing HCL: %s", diag.Error())
		return nil, diag
	}

	var metadata TerraformLockFile
	decodeCtx := &hcl.EvalContext{}
	diag = gohcl.DecodeBody(file.Body, decodeCtx, &metadata)
	if diag.HasErrors() {
		log.Debugf("Error decoding HCL body: %s", diag.Error())
		return nil, diag
	}

	return &metadata, nil
}
