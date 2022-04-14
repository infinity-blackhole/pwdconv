package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

type BitwardenRecord struct {
	Folder        string
	Favorite      string
	Type          string
	Name          string
	Notes         string
	Fields        string
	Reprompt      bool
	LoginUri      string
	LoginUsername string
	LoginPassword string
	LoginTOTP     string
}

type GooglePasswordRecord struct {
	Uri      string
	Username string
	Password string
}

var rootCmd = &cobra.Command{
	Use: "pwdconv",
	Run: func(cmd *cobra.Command, args []string) {
		for _, file := range args {
			f1, err := os.Open(file)
			if err != nil {
				log.Fatal(err)
			}
			defer f1.Close()
			csvReader := csv.NewReader(f1)

			f2, err := os.Create(fmt.Sprintf("google-%s", file))
			if err != nil {
				log.Fatal(err)
			}
			defer f2.Close()
			csvWriter := csv.NewWriter(f2)
			if err := csvWriter.Write([]string{"uri", "username", "password"}); err != nil {
				log.Fatal(err)
			}
			defer csvWriter.Flush()

			for i := 1; i <= SkipLeardingRows; i++ {
				_, err := csvReader.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
				}
			}
			for {
				rec, err := csvReader.Read()
				if err != nil {
					if err == io.EOF {
						break
					}
					log.Fatal(err)
				}

				reprompt, err := strconv.ParseBool(rec[6])
				if err != nil {
					log.Fatal(err)
				}

				bitwRec := BitwardenRecord{
					Folder:        rec[0],
					Favorite:      rec[1],
					Type:          rec[2],
					Name:          rec[3],
					Notes:         rec[4],
					Fields:        rec[5],
					Reprompt:      reprompt,
					LoginUri:      rec[7],
					LoginUsername: rec[8],
					LoginTOTP:     rec[9],
				}

				googRec := GooglePasswordRecord{
					Uri:      bitwRec.LoginUri,
					Username: bitwRec.LoginUsername,
					Password: bitwRec.LoginTOTP,
				}

				line := []string{googRec.Uri, googRec.Username, googRec.Password}
				if err := csvWriter.Write(line); err != nil {
					log.Fatal(err)
				}
			}
		}
	},
}

var SkipLeardingRows int

func init() {
	rootCmd.Flags().IntVar(&SkipLeardingRows, "skip-leading-rows", 0, "Number of rows to skip")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
