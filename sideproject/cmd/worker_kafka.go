package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var workerKafka = &cobra.Command{
	Use:   "worker_kafka",
	Short: "worker_kafka",
	Long:  "worker_kafka",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kafka consumer....end")
	},
}
