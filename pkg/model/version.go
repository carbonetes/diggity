package model

// Version - Build Information
type Version struct {
	// Version
	AppName   string `json:"appName"`
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	// Git
	GitCommit string `json:"gitCommit"`
	GitDesc   string `json:"gitDesc"`
	// Golang
	GoVersion string `json:"goVersion"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}
