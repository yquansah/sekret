package cmd

import (
	"fmt"
	"runtime"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version information can be set at build time using ldflags:
// go build -ldflags "-X github.com/yquansah/sekret/cmd.version=v1.0.0 -X github.com/yquansah/sekret/cmd.commit=abc123 -X github.com/yquansah/sekret/cmd.date=2024-01-01T00:00:00Z"

var (
	version   = "dev"
	commit    = "unknown"
	date      = "unknown"
	goVersion = runtime.Version()
)

func GetVersion() string {
	initVersion()
	return version
}

func initVersion() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if version == "dev" {
			version = info.Main.Version
			if version == "(devel)" {
				version = "dev"
			}
		}

		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				if commit == "unknown" {
					commit = setting.Value
				}
			case "vcs.time":
				if date == "unknown" {
					date = setting.Value
				}
			}
		}
	}
}

func init() {
	initVersion()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information for sekret including build details.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sekret version %s\n", version)
		fmt.Printf("  commit: %s\n", commit)
		fmt.Printf("  built: %s\n", date)
		fmt.Printf("  go version: %s\n", goVersion)
	},
}