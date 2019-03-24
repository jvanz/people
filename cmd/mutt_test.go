package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/user"
	"testing"
)

func TestGetDataFilePath(t *testing.T) {
	filename := getDataFilename()
	if len(filename) == 0 {
		t.Error("Missing filename")
	}
	user, err := user.Current()
	if err != nil {
		t.Error(err)
	}
	if filename != fmt.Sprintf("%s/.people/data", user.HomeDir) {
		t.Errorf("Invalid data file: %s", filename)
	}
}

func TestFormatOutput(t *testing.T) {
	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	person := People{
		ID:       id,
		Nickname: "test",
		Name:     "Test User",
		Email:    "test@test.com",
	}
	str := formatPeopleOutput(&person)
	expected := fmt.Sprintf("%s\t%s\t%s\n", person.Email, person.Name, " ")
	if str != expected {
		t.Errorf("Invalid format: %s", str)
	}
}

func TestWriteEntry(t *testing.T) {
	// clean database
	err := os.Remove(getDataFilename())
	if err != nil {
		t.Error(err)

	}
	//  create a fake entry
	id, err := uuid.NewRandom()
	if err != nil {
		t.Error(err)
	}
	person := People{
		ID:       id,
		Nickname: "test",
		Name:     "Test User",
		Email:    "test@test.com",
	}
	jsondata, err := json.Marshal(person)
	if err != nil {
		panic(err)
	}
	writeJsonFile(jsondata)

	file, err := os.Open(getDataFilename())
	if err != nil {
		t.Error(err)
	}
	defer file.Close()

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		count = count + 1
		line := scanner.Bytes()
		if !json.Valid(line) {
			t.Error("Invalid JSON")
		}
		var person2 People
		err := json.Unmarshal(line, &person2)
		if err != nil {
			t.Error(err)
		}
		if person.ID != person2.ID {
			t.Errorf("Invalid ID: %s", person2.ID)
		}
		if person.Nickname != person2.Nickname {
			t.Errorf("Invalid Nickname: %s", person2.Nickname)
		}
		if person.Name != person2.Name {
			t.Errorf("Invalid Name: %s", person2.Name)
		}
		if person.Email != person2.Email {
			t.Errorf("Invalid Email: %s", person2.Email)
		}
	}
	if count != 1 {
		t.Errorf("Written one entry but read %d", count)
	}
}

func TestReadEntry(t *testing.T) {
	// clean database
	err := os.Remove(getDataFilename())
	if err != nil {
		t.Error(err)

	}
	//  create a fake entries
	all_people := make(map[uuid.UUID]People)
	for i := 0; i < 2; i++ {
		id, err := uuid.NewRandom()
		if err != nil {
			t.Fatal(err)
		}
		person := People{
			ID:       id,
			Nickname: fmt.Sprintf("test %d", i),
			Name:     fmt.Sprintf("Test User %d", i),
			Email:    fmt.Sprintf("test%d@test.com", i),
		}
		jsondata, err := json.Marshal(person)
		if err != nil {
			t.Fatal(err)
		}
		writeJsonFile(jsondata)
		all_people[person.ID] = person
	}
	// read database
	address_book := loadPeople()
	if len(address_book) != len(all_people) {
		t.Errorf("Address book has %d entries. Expected %d", len(address_book), len(all_people))
	}
	for _, person := range address_book {
		if _, ok := all_people[person.ID]; !ok {
			t.Errorf("Entry %s not expected", person)
		} else {
			delete(all_people, person.ID)
		}
	}
	if len(all_people) > 0 {
		for _, p := range all_people {
			t.Errorf("Missing %s", p)
		}
	}

}
