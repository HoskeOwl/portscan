/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	ctxstorage "github.com/HoskeOwl/portscan/internal/ctx_storage"
	"github.com/HoskeOwl/portscan/internal/helper"
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

func checkPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "portscan",
	Short: "Simple port scaner",
	Long: `Fast port scanner with parallel execution.
Do not send data, just check connection.
Support json output format.
`,
	Run: func(cmd *cobra.Command, args []string) {
		portRanges, err := cmd.Flags().GetString("port")
		checkPanic(err)
		if portRanges == "" {
			panic(fmt.Errorf("wrong port range format: '%v'", portRanges))
		}
		parsedPortRanges, err := helper.ParsePortRanges(portRanges)
		checkPanic(err)

		host, err := cmd.Flags().GetString("dst")
		checkPanic(err)

		timeoutMs, err := cmd.Flags().GetInt("timeout")
		checkPanic(err)

		connections, err := cmd.Flags().GetInt("connections")
		checkPanic(err)

		retries, err := cmd.Flags().GetInt("retries")
		checkPanic(err)

		verbose, err := cmd.Flags().GetBool("verbose")
		checkPanic(err)

		sort, err := cmd.Flags().GetBool("sort")
		checkPanic(err)

		json, err := cmd.Flags().GetBool("json")
		checkPanic(err)

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
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("port", "p", "", fmt.Sprintf("Port or range. Can be several ranges/ports. Example: '2%v80%v100%v8080'", helper.RangesSeparator, helper.PortRangeSeparator, helper.RangesSeparator))
	rootCmd.Flags().StringP("dst", "d", defaultHost, "Destination address")
	rootCmd.Flags().IntP("timeout", "t", defaultTimeoutMs, "Timeout in milliseconds for each connection")
	rootCmd.Flags().IntP("connections", "c", defaultConnectionsCount, "Connection pool")
	rootCmd.Flags().IntP("retries", "r", defaultRetries, "How many times check unavailable port")
	rootCmd.Flags().BoolP("verbose", "v", false, "Print failed ports")
	rootCmd.Flags().BoolP("sort", "s", false, "Sort output ports (print when checks all)")
	rootCmd.Flags().BoolP("json", "j", false, "Json output (ignore -v and -s)")
}
