package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

// Define function to read CSV file and return an array of objects
func readCSV(filePath string) (map[string]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	accounts := make(map[string]string)
	for _, record := range records {
		accounts[record[0]] = record[1]
	}
	return accounts, nil
}

// Define function to read JSON file and return its contents
func readJSON(filePath string) (map[string]interface{}, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data map[string]interface{}
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Define function to verify address and balance in JSON file against CSV list
func verifyJSON(filePath string) error {
	// Read 80% csv file
	eightyPercentMap, err := readCSV("./.github/workflows/80_percent_accounts.csv")
	if err != nil {
		return err
	}

	// Read 10 utia csv file
	tenUtiaMap, err := readCSV("./.github/workflows/10_utia_accounts.csv")
	if err != nil {
		return err
	}

	// Read gentx
	jsonData, err := readJSON(filePath)
	if err != nil {
		return err
	}

	// Parse the json
	body := jsonData["body"].(map[string]interface{})
	messages := body["messages"].([]interface{})
	if len(messages) != 1 {
		log.Fatalf("Invalid number of messages: expected 1 got %v", len(messages))
	}

	// Grab the address and balance info
	delegatorAddress := messages[0].(map[string]interface{})["delegator_address"].(string)
	value := messages[0].(map[string]interface{})["value"].(map[string]interface{})
	denom := value["denom"]
	if denom != "utia" {
		return fmt.Errorf("Invalid denom: expected utia got %v", denom)
	}
	amount := value["amount"]

	expectedAmount80, ok80 := eightyPercentMap[delegatorAddress]
	_, ok10 := tenUtiaMap[delegatorAddress]
	switch {
	case ok80:
		// convert string to uint and verify that amount is at least 80% of expected
		expectedAmountInt, err := strconv.ParseUint(expectedAmount80, 10, 64)
		if err != nil {
			return err
		}
		amountInt, err := strconv.ParseUint(amount.(string), 10, 64)
		if err != nil {
			return err
		}
		if amountInt < expectedAmountInt*8/10 {
			return fmt.Errorf("Invalid amount: expected at least %v got %v", expectedAmountInt*8/10, amountInt)
		}
	case ok10:
		// convert amount to uint and verify it is between 1 and 9 tia
		amountInt, err := strconv.ParseUint(amount.(string), 10, 64)
		if err != nil {
			return err
		}
		// 1 tia = 1000000 utia
		if amountInt < 1000000 || amountInt > 9000000 {
			return fmt.Errorf("Invalid amount: expected between 1000000 and 9000000 utia, got %v", amountInt)
		}
	default:
		return errors.New("Address not found in approved list")
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide the path to the JSON file")
	}
	jsonFilePath := os.Args[1]
	err := verifyJSON(jsonFilePath)
	if err != nil {
		log.Fatal(err)
	}
}
