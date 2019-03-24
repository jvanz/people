package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/user"
	"regexp"
	"strings"
)

// People is the type used to store all the data from a entry in the database
type People struct {
	ID       uuid.UUID
	Nickname string
	Name     string
	Email    string
}

// mutt is the command used to mutt to search and add new entries in the database
var mutt = &cobra.Command{
	Use:   "mutt",
	Short: "Command used to mutt to get the address",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// list is the subcommand within the mutt command to search for address in the
// database
var list = &cobra.Command{
	Use:   "list",
	Short: "Command used to mutt to get the address",
	Run: func(cmd *cobra.Command, args []string) {
		var address_book = loadPeople()
		for _, people := range address_book {
			if strings.Contains(strings.ToLower(people.Nickname), args[0]) ||
				strings.Contains(strings.ToLower(people.Name), args[0]) ||
				strings.Contains(people.Email, args[0]) {
				formatOutput(&people)
			}
		}
	},
}

// add is the subcommand within the mutt command to add new address in the
// database
var add = &cobra.Command{
	Use:   "add",
	Short: "Command used to mutt to add address",
	Run: func(cmd *cobra.Command, args []string) {
		scanner := bufio.NewScanner(os.Stdin)
		var re = regexp.MustCompile(`From: (\"?(?P<name>[\w ]*)\"? )?<(?P<email>.*)>`)
		for scanner.Scan() {
			line := scanner.Text()
			if re.MatchString(line) {
				email := re.ReplaceAllString(line, fmt.Sprintf("${%s}", re.SubexpNames()[3]))
				id, err := uuid.NewRandom()
				if err != nil {
					panic(err)
				}
				new_people := People{
					ID:       id,
					Nickname: strings.Split(email, "@")[0],
					Name:     re.ReplaceAllString(line, fmt.Sprintf("${%s}", re.SubexpNames()[2])),
					Email:    email,
				}
				var address_book = loadPeople()
				add := true
				for _, people := range address_book {
					if isTheSamePerson(&people, &new_people) {
						fmt.Printf("%s found in the database, skipping...\n", new_people.Email)
						add = false
						break
					}
				}
				if add {
					json, err := json.Marshal(new_people)
					if err != nil {
						panic(err)
					}
					writeJsonFile(json)
					fmt.Printf("Added %s\n", new_people.Email)
				}
			}
		}
	},
}

func init() {
	mutt.AddCommand(list)
	mutt.AddCommand(add)
	rootCmd.AddCommand(mutt)
}

// isTheSamePerson checks if the p1 and p2 is the same entry. In other words,
// if the data from both as equivalent
func isTheSamePerson(p1, p2 *People) bool {
	return p1.Email == p2.Email ||
		strings.Contains(strings.ToLower(p1.Nickname), strings.ToLower(p2.Nickname)) ||
		strings.Contains(strings.ToLower(p1.Name), strings.ToLower(p2.Nickname))
}

// writeJsonFile writes the json data of a new entry in the database
func writeJsonFile(json []byte) {
	f, err := os.OpenFile(getDataFilename(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	entry := fmt.Sprintf("%s\n", json)
	w := bufio.NewWriter(f)
	nn, err := w.WriteString(entry)
	if nn < len(json) {
		log.Print(err)
	}
	w.Flush()
}

// loadPeople reads the databse and returns a slice with all entries read
func loadPeople() []People {
	file, err := os.Open(getDataFilename())
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var address_book []People
	for scanner.Scan() {
		line := scanner.Bytes()
		if !json.Valid(line) {
			panic("Invalid JSON")
		}
		var people People
		err := json.Unmarshal(line, &people)
		if err != nil {
			panic("Unmarshal failed")
		}
		address_book = append(address_book, people)
	}
	return address_book
}

// formatOutput prints in the stdout the data from the given person. This printed
// data is formated to allow mutt read it and used the address
func formatOutput(person *People) {
	fmt.Printf("\n%s\t%s\t%s\n", person.Email, person.Name, " ")
}

// getDataFilename returns a string of the database file path
func getDataFilename() string {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%s/.people/data", user.HomeDir)
}
