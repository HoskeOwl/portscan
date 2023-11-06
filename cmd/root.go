/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
	"github.com/HoskeOwl/portscan/internal/helper"
	"github.com/HoskeOwl/portscan/internal/version"
	"github.com/HoskeOwl/portscan/internal/worker"
	"github.com/spf13/cobra"
)

const (
	defaultTimeoutMs        int    = 1000
	defaultConnectionsCount int    = 50
	defaultRetries          int    = 2
	defaultHost             string = "localhost"
	defaultVerbose          bool   = false
)

func checkCmdError(cmd *cobra.Command, err error) {
	if err == nil {
		return
	}
	cmd.Println(err)
	cmd.Println(cmd.UsageString())
	os.Exit(1)
}

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	SilenceErrors: true,
	Use:           "portscan",
	Short:         "Simple port scaner",
	Long: `Fast port scanner with parallel execution.
Do not send data, just check connection.
Support json output format.
`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := cmd.Flags().GetBool("version")
		checkCmdError(cmd, err)
		if v {
			fmt.Println(version.BuildVersion())
			os.Exit(0)
		}

		portRanges, err := cmd.Flags().GetString("port")
		checkCmdError(cmd, err)
		if portRanges == "" {
			fmt.Printf("wrong port range format: '%v'\n", portRanges)
			cmd.Println(cmd.UsageString())
			os.Exit(1)
		}
		parsedPortRanges, err := helper.ParsePortRanges(portRanges)
		checkCmdError(cmd, err)

		host, err := cmd.Flags().GetString("dst")
		checkCmdError(cmd, err)

		timeoutMs, err := cmd.Flags().GetInt("timeout")
		checkCmdError(cmd, err)

		connections, err := cmd.Flags().GetInt("connections")
		checkCmdError(cmd, err)

		retries, err := cmd.Flags().GetInt("retries")
		checkCmdError(cmd, err)

		verbose, err := cmd.Flags().GetBool("verbose")
		checkCmdError(cmd, err)

		sort, err := cmd.Flags().GetBool("sort")
		checkCmdError(cmd, err)

		json, err := cmd.Flags().GetBool("json")
		checkCmdError(cmd, err)

		storage := ctxstorage.CtxStorage{
			ConnDuration: time.Duration(timeoutMs * int(time.Millisecond)),
			Retries:      retries,
			Verbose:      verbose,
		}

		var hooker worker.ResultProcessor
		switch {
		case json:
			hooker = worker.MakeJsonResultProcessor()
		case sort:
			hooker = worker.MakeSortedResultProcessor()
		default:
			hooker = worker.MakeRealtimeResultProcessor()
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, ctxstorage.StorageKey, storage)

		if !json {
			fmt.Printf("Scan '%v' for port ranges '%v' ...\n", host, portRanges)
		}
		start := time.Now()

		scanTasks := helper.CreateScanTasks(host, parsedPortRanges)
		if len(scanTasks) < connections {
			connections = len(scanTasks)
		}
		pool := worker.MakePool(ctx, scanTasks, connections).WithSuccessHook(hooker.Success).WithFailHook(hooker.Fail)
		pool.Execute(ctx)

		elapsed := time.Since(start)
		hooker.Print(ctx)
		if !json {
			fmt.Printf("\nDone   (took %.3f seconds)\n", (float64(elapsed.Milliseconds()) / 1000))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var ErrSilent = errors.New("ErrSilent")

func init() {
	RootCmd.Flags().BoolP("version", "v", false, "Print program version and exit")
	RootCmd.Flags().StringP("port", "p", "", fmt.Sprintf("Port or range. Can be several ranges/ports. Example: '2%v80%v100%v8080'", helper.RangesSeparator, helper.PortRangeSeparator, helper.RangesSeparator))
	RootCmd.Flags().StringP("dst", "d", defaultHost, "Destination address")
	RootCmd.Flags().IntP("timeout", "t", defaultTimeoutMs, "Timeout in milliseconds for each connection")
	RootCmd.Flags().IntP("connections", "c", defaultConnectionsCount, "Connection pool")
	RootCmd.Flags().IntP("retries", "r", defaultRetries, "How many times check unavailable port")
	RootCmd.Flags().BoolP("verbose", "b", false, "Print failed ports")
	RootCmd.Flags().BoolP("sort", "s", false, "Sort output ports (print when checks all)")
	RootCmd.Flags().BoolP("json", "j", false, "Json output (ignore -v and -s)")
}
