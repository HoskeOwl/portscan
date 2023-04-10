/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/HoskeOwl/portscan/internal/scanworker"
	"github.com/spf13/cobra"
)

const (
	defaultTimeoutMs        int    = 1000
	defaultConnectionsCount int    = 1000
	defaultHost             string = "localhost"
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
	Long: `Fast port scanner for huge amount of ports.
Do not send data, just check if can connect.
`,
	Run: func(cmd *cobra.Command, args []string) {
		portRanges, err := cmd.Flags().GetString("port")
		checkPanic(err)
		if portRanges == "" {
			panic(fmt.Errorf("wrong port range format: '%v'", portRanges))
		}
		parsedPortRanges, err := scanworker.ParsePortRanges(portRanges)
		checkPanic(err)

		host, err := cmd.Flags().GetString("dst")
		checkPanic(err)

		timeoutMs, err := cmd.Flags().GetInt("timeout")
		checkPanic(err)

		connections, err := cmd.Flags().GetInt("connections")
		checkPanic(err)

		ctx := context.Background()
		success, err := scanworker.StartPortScan(ctx, connections, host, parsedPortRanges, timeoutMs)
		if len(success) == 0 {
			fmt.Println("No opened ports")
		} else {
			for _, port := range success {
				fmt.Printf("\t%v: success\n", port)
			}
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
	rootCmd.Flags().StringP("port", "p", "", "Port or range. Can be several ranges/ports. Example: '2,80:100,8080'")
	rootCmd.Flags().StringP("dst", "d", defaultHost, "Destination address")
	rootCmd.Flags().IntP("timeout", "t", defaultTimeoutMs, "Timeout in milliseconds for each connection")
	rootCmd.Flags().IntP("connections", "c", defaultConnectionsCount, "Connection pool")
}
