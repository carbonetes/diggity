package command

import (
	"errors"
	"os"
	"time"

	"github.com/carbonetes/diggity/cmd/diggity/build"
	"github.com/carbonetes/diggity/cmd/diggity/config"
	"github.com/carbonetes/diggity/cmd/diggity/ui"
	"github.com/carbonetes/diggity/internal/helper"
	"github.com/carbonetes/diggity/internal/log"
	"github.com/carbonetes/diggity/pkg/cdx"
	"github.com/carbonetes/diggity/pkg/cdx/dependency"
	"github.com/carbonetes/diggity/pkg/reader"
	"github.com/carbonetes/diggity/pkg/types"
	"github.com/golistic/urn"
	"github.com/spf13/cobra"
)

const (
	shortDesc = "BOM Diggity Scanner"
	longDesc  = `BOM Diggity is an open-source tool developed to streamline the critical process of generating a comprehensive Software Bill of Materials (SBOM) for Container Images and File Systems across various supported ecosystems.`
)

var scan = &cobra.Command{
	Use:   "diggity",
	Args:  cobra.MaximumNArgs(1),
	Short: shortDesc,
	Long:  longDesc,
	Run:   scanCmd,
}

func scanCmd(c *cobra.Command, args []string) {
	versionArg, _ := c.Flags().GetBool("version")
	if versionArg {
		log.Print(build.FromBuild().Version)
		os.Exit(0)
	}

	tarball, _ := c.Flags().GetString("tar")
	filesystem, _ := c.Flags().GetString("dir")

	params := types.DefaultParameters()
	if err := setParams(&params, c, filesystem, tarball, args); err != nil {
		log.Fatal(err)
	}

	quiet, err := c.Flags().GetBool("quiet")
	if err != nil {
		log.Debug(err.Error())
	}

	outputFormat, err := c.Flags().GetString("output")
	if err != nil {
		log.Debug(err.Error())
	}

	file, err := c.Flags().GetString("file")
	if err != nil {
		log.Debug(err.Error())
	}

	if len(file) > 0 {
		params.SaveToFile = file
	}

	if !types.IsValidOutputFormat(outputFormat) {
		log.Debug("Invalid output format parameter")
	}

	params.Quiet = quiet
	params.SaveToFile = file
	params.OutputFormat = types.OutputFormat(outputFormat)

	if !params.Quiet {
		ui.Run()
	} else {
		if params.OutputFormat == types.Table {
			params.OutputFormat = types.JSON
		}
	}

	if len(params.Input) == 0 {
		log.Fatal("No input provided")
	}

	start := time.Now()
	if addr, err := Scan(params); err != nil {
		ui.Error(err)
		os.Exit(1)
	} else {
		ui.DisplayResult(params, time.Since(start).Seconds(), addr)
	}
}

func setParams(params *types.Parameters, c *cobra.Command, filesystem, tarball string, args []string) error {
	if filesystem != "" {
		if found, _ := helper.IsDirExists(filesystem); !found {
			return errors.New("directory not found: " + filesystem)
		}
		params.ScanType = types.Filesystem
		params.Input = filesystem
	}

	if tarball != "" {
		if found, _ := helper.IsFileExists(tarball); !found {
			return errors.New("tarball not found: " + tarball)
		}
		params.Input = tarball
		params.ScanType = types.Tarball
	}

	if filesystem == "" && tarball == "" {
		if len(args) > 0 {
			params.Input = helper.FormatImage(args[0])
			params.ScanType = types.Image
		} else {
			_ = c.Help()
			os.Exit(0)
		}
	}

	return nil
}

func Scan(parameters types.Parameters) (*urn.URN, error) {

	// Generate unique address for the scan
	addr, err := types.NewAddress()
	if err != nil {
		log.Fatal(err)
	}

	cdx.New(addr)
	dependency.NewDependencyNodes(addr)

	switch parameters.ScanType {
	case 1:
		if err := handleImageScan(parameters, addr); err != nil {
			return nil, err
		}
	case 2:
		if err := handleTarballScan(parameters, addr); err != nil {
			return nil, err
		}
	case 3:
		if err := handleFilesystemScan(parameters, addr); err != nil {

			return nil, err
		}
	default:
		return nil, errors.New("invalid scan type")
	}

	return addr, nil
}

func handleImageScan(parameters types.Parameters, addr *urn.URN) error {
	target := parameters.Input
	image, ref, err := reader.GetImage(target, &config.Config.Registry)
	if err != nil {
		return err
	}

	if image == nil {
		return errors.New("image not found")
	}
	if ref != nil {
		cdx.SetMetadataComponent(addr, cdx.SetImageMetadata(*image, *ref, target))
	}

	err = reader.ReadFiles(image, addr)
	if err != nil {
		return err
	}

	return nil
}

func handleTarballScan(parameters types.Parameters, addr *urn.URN) error {
	image, err := reader.ReadTarball(parameters.Input)
	if err != nil {
		return err
	}
	err = reader.ReadFiles(image, addr)
	if err != nil {
		return err
	}

	return nil
}

func handleFilesystemScan(parameters types.Parameters, addr *urn.URN) error {
	err := reader.FilesystemScanHandler(parameters.Input, addr)
	if err != nil {
		return err
	}

	return nil
}
