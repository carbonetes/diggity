package eventbus

import (
	sbom "github.com/carbonetes/diggity/internal"
	log "github.com/carbonetes/diggity/internal/logger"
	"github.com/carbonetes/diggity/internal/ui"
	"github.com/carbonetes/diggity/pkg/model"

	"github.com/vmware/transport-go/bus"
	tm "github.com/vmware/transport-go/model"
)

// SetAnalysisRequestHandler for event bus
func SetAnalysisRequestHandler(channelName string) {
	log := log.GetLogger()
	tr := bus.GetBus()
	ui.Disable()
	var event string = "event"
	requestHandler, _ := tr.ListenRequestStream(channelName)
	requestHandler.Handle(
		func(msg *tm.Message) {
			arguments := msg.Payload.(model.Arguments)
			arguments.Output = (*model.Output)(&event)
			sbom.Start(&arguments)
			result := sbom.GetResults()
			if err := tr.SendResponseMessage(channelName, result, msg.DestinationId); err != nil {
				log.Fatalf("Error sending response message: %v", err)
			}
		},
		func(err error) {
			log.Fatalf("Error handling request: %v", err)
		},
	)
}

// GetArguments for event bus
func GetArguments() model.Arguments {
	return model.Arguments{
		DisableFileListing:  new(bool),
		SecretContentRegex:  new(string),
		DisableSecretSearch: new(bool),
		Dir:                 new(string),
		Tar:                 new(string),
		Quiet:               new(bool),
		OutputFile:          new(string),
		ExcludedFilenames:   &[]string{},
		EnabledParsers:      &[]string{},
		RegistryURI:         new(string),
		RegistryUsername:    new(string),
		RegistryPassword:    new(string),
		RegistryToken:       new(string),
	}
}
