package scanner

import (
	"github.com/carbonetes/diggity/internal/scanner/distro"
	"github.com/carbonetes/diggity/internal/scanner/os/apk"
	"github.com/carbonetes/diggity/internal/scanner/os/dpkg"
	"github.com/carbonetes/diggity/internal/scanner/secret"
	"github.com/carbonetes/diggity/pkg/stream"
)

func Init() {
	stream.Attach(apk.Type, apk.Scan)
	stream.Attach(dpkg.Type, dpkg.Scan)
	stream.Attach(distro.Type, distro.Scan)
	stream.Attach(secret.Type, secret.Scan)
}
