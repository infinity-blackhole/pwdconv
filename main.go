package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var filepath string

var rootCmd = &cobra.Command{
	Use: "pwdconv",
	Run: func(cmd *cobra.Command, args []string) {
		for file := range args {
			f, err := os.Open(file)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			csvReader := csv.NewReader(f)
			data, err := csvReader.ReadAll()
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
