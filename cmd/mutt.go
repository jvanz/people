package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var mutt = &cobra.Command{
	Use:   "mutt",
	Short: "Command used to mutt to get the address",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.Open("/home/jvanz/.people/address")
		if err != nil {
			panic(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Bytes()
			if !json.Valid(line) {
				panic("Invalid JSON")
			}
			var v map[string]interface{}
			err := json.Unmarshal(line, &v)
			if err != nil {
				panic("Unmarshal failed")
			}
			for _, value := range v {
				if strings.Contains(value.(string), args[0]) {
					formatOutput(v)
					break
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(mutt)
}

func formatOutput(people map[string]interface{}) {
	fmt.Printf("%s\t%s\t%s\n", people["email"], people["name"], people["nickname"])
}
