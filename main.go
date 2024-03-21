package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/openshift/backplane-tools/cmd/install"
	"github.com/openshift/backplane-tools/cmd/list"
	"github.com/openshift/backplane-tools/cmd/remove"
	"github.com/openshift/backplane-tools/cmd/upgrade"
	"github.com/spf13/cobra"
)

var cmd = cobra.Command{
	Use:   "backplane-tools",
	Short: "An OpenShift tool manager",
	Long:  "This applications manages the tools needed to interact with OpenShift clusters",
	RunE:  help,
}

var (
	cpuProfilingEnvVar = "CPU_PROFILE"
	memProfilingEnvVar = "MEM_PROFILE"
)

func help(cmd *cobra.Command, _ []string) error {
	return cmd.Help()
}

func init() {
	// Subcommands
	cmd.AddCommand(install.Cmd())
	cmd.AddCommand(list.Cmd())
	cmd.AddCommand(remove.Cmd())
	cmd.AddCommand(upgrade.Cmd())
}

func main() {
	cpuProfilePath, found := os.LookupEnv(cpuProfilingEnvVar)
	if found && cpuProfilePath != "" {
		log.Printf("Creating cpu profile at '%s'", cpuProfilePath)
		profile, err := os.Create(cpuProfilePath)
		if err != nil {
			log.Fatalf("failed to create CPU profile '%s': %v", cpuProfilePath, err)
		}

		defer func() {
			closeErr := profile.Close()
			if closeErr != nil {
				fmt.Fprintf(os.Stderr, "failed to close file '%s': %v", cpuProfilePath, err)
			}
		}()

		err = pprof.StartCPUProfile(profile)
		if err != nil {
			log.Fatalf("failed to initiate CPU profiling: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	err := cmd.Execute()
	if err != nil {
		log.Fatalf("Error executing command: %v", err)
	}

	memProfilePath, found := os.LookupEnv(memProfilingEnvVar)
	if found && memProfilePath != "" {
		log.Printf("Creating memory profile at '%s'", memProfilePath)
		profile, err := os.Create(memProfilePath)
		if err != nil {
			log.Fatalf("failed to create memory profile '%s': %v", memProfilePath, err)
		}

		defer func() {
			closeErr := profile.Close()
			if closeErr != nil {
				fmt.Fprintf(os.Stderr, "failed to close file '%s': %v", memProfilePath, err)
			}
		}()

		runtime.GC()
		err = pprof.WriteHeapProfile(profile)
		if err != nil {
			log.Fatalf("failed to initiate memory profiling: %v", err)
		}
	}
}
