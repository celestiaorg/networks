package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: checkuser <user> <pr>")
		os.Exit(1)
	}

	user := os.Args[1]
	pr := os.Args[2]

	if err := checkUser(user, pr); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func checkUser(user string, pr string) error {
	// Open the CSV file
	file, err := os.OpenFile("./.github/workflows/usersprs.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Read the CSV file
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	// Check if the user exists in the CSV file
	var userExists bool
	for _, record := range records {
		if record[0] == user {
			userExists = true
			if record[1] == pr {
				fmt.Println("Pass: User exists with the same PR number")
				return nil
			} else {
				fmt.Println("Fail: User exists with different PR number")
				return fmt.Errorf("User %s already exists with PR number %s", user, record[1])
			}
		}
	}

	// If the user does not exist, add them to the CSV file
	if !userExists {
		w := csv.NewWriter(file)
		defer w.Flush()

		if err := w.Write([]string{user, pr}); err != nil {
			return err
		}

		fmt.Printf("User %s added with PR number %s\n", user, pr)
	}

	return nil
}
