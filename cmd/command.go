package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"node-allocation-exporter/exporter"
	"os"
)

// command define command-cli
var command = &cobra.Command{
	Use:   "node_allocation_exporter",
	Short: "Prometheus exporter for node allocation metrics",
	Long:  `Prometheus exporter collecting allocation of the node.`,
	Run: func(cmd *cobra.Command, args []string) {
		e := exporter.NewExporter(viper.GetString("bind-address"))
		e.RunServer()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	flags := command.PersistentFlags()
	flags.StringP("bind-address", "b", "0.0.0.0:9965", "Address to bind to")

	viper.BindPFlags(flags)
}

func Execute() {
	if err := command.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}
