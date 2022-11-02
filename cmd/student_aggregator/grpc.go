/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package studentaggregator

import (
	"fmt"

	"github.com/spf13/cobra"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		portNum, _ := cmd.Flags().GetInt("port")

		fmt.Printf("grpc server called with port number %d\n", portNum)
	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)

	grpcCmd.PersistentFlags().Int("port", 9090, "grpc server port to be served")
}
