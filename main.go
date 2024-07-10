package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	req "github.com/dombartenope/filterfinder/req"
)

type Filter struct {
	Key      string `json:"key"`
	Field    string `json:"field"`
	Value    string `json:"value"`
	Relation string `json:"relation"`
}

type Operator struct {
	Operator string `json:"operator"`
}

func main() {
	userStruct := req.ViewUser()

	tagData, err := json.Marshal(userStruct.Properties.Tags)
	if err != nil {
		log.Fatalf("Error marshalling the JSON tag data: %v", err)
	}

	languageData, err := json.Marshal(userStruct.Properties.Language)
	if err != nil {
		log.Fatalf("Error marshalling the JSON tag data: %v", err)
	}

	//Input User Tags here
	userTags := string(tagData)
	userLanguage := string(languageData)
	fmt.Printf("\nAttempting to find match given this input: %s\n", userTags)
	fmt.Printf("Attempting to find match given this input: %s\n\n", userLanguage)
	count := 0

	//Read the input of the csv here and create an output file to put matches in
	input, err := os.Open("input.csv")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer input.Close()
	out, err := os.Create("out.csv")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer out.Close()

	reader := csv.NewReader(input)
	writer := csv.NewWriter(out)
	defer writer.Flush()

	//Column Names and write to CSV
	col, err := reader.Read()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	colAndIndex := make(map[string]int)
	for i, v := range col {
		colAndIndex[v] = i
	}

	row, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	for _, v := range row {
		if v[colAndIndex["filters"]] != "" {
			csvFilters := v[colAndIndex["filters"]]
			filterGroups := parseFilters(csvFilters)
			userTagsMap := parseUserTags(userTags)
			userLang := parseUserLanguage(userLanguage)
			// Check if any group of conditions is met
			match := false
			for i, group := range filterGroups {
				groupMatch := checkConditionGroup(group, userTagsMap, userLang)
				if groupMatch {
					fmt.Printf("Group %d: %v\n", i+1, groupMatch)
					match = true
				}
			}
			if match {
				writer.Write(v)
				fmt.Println("Match found! At least one group of conditions is met.")
				fmt.Println(csvFilters)
				count++
			}
		}
	}
	fmt.Printf("%d found", count)
}

func parseFilters(csvFilters string) [][]Filter {
	// Remove the "=>" syntax and replace with ":"
	csvFilters = strings.ReplaceAll(csvFilters, `""`, `"`)
	csvFilters = strings.ReplaceAll(csvFilters, "=>", ":")

	var allFilters []json.RawMessage
	err := json.Unmarshal([]byte(csvFilters), &allFilters)
	if err != nil {
		fmt.Println("Error parsing filters:", err)
		return nil
	}

	var filterGroups [][]Filter
	var currentGroup []Filter

	for _, rawFilter := range allFilters {
		var operator Operator
		err := json.Unmarshal(rawFilter, &operator)
		if err == nil && operator.Operator == "OR" {
			if len(currentGroup) > 0 {
				filterGroups = append(filterGroups, currentGroup)
				currentGroup = []Filter{}
			}
		} else {
			var filter Filter
			err := json.Unmarshal(rawFilter, &filter)
			if err == nil {
				currentGroup = append(currentGroup, filter)
			}
		}
	}

	if len(currentGroup) > 0 {
		filterGroups = append(filterGroups, currentGroup)
	}

	return filterGroups
}

func parseUserLanguage(userLanguage string) string {

	userLanguage = strings.TrimSpace(userLanguage)
	userLanguage = strings.Trim(userLanguage, "{}")
	userLanguage = strings.ReplaceAll(userLanguage, "\"", "")

	return userLanguage
}

func parseUserTags(userTags string) map[string]string {
	// Remove the curly braces
	userTags = strings.TrimSpace(userTags)
	userTags = strings.Trim(userTags, "{}")
	userTags = strings.ReplaceAll(userTags, "\"", "")

	// Split the section into key-value pairs
	pairs := strings.Split(userTags, ",")
	userTagsMap := make(map[string]string)

	for _, pair := range pairs {
		kv := strings.SplitN(pair, ":", 2)
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			userTagsMap[key] = value
		}
	}

	return userTagsMap
}

func checkConditionGroup(group []Filter, userTags map[string]string, userLanguage string) bool {
	for _, filter := range group {
		if !checkCondition(filter, userTags, userLanguage) {
			return false
		}
	}
	return true
}

// Check the validity of the relation and whether this script supports it
func checkCondition(filter Filter, userTags map[string]string, userLanguage string) bool {

	if filter.Field == "language" {
		return compareValues(userLanguage, filter.Value, filter.Relation)
	}

	switch filter.Relation {
	case "exists":
		_, exists := userTags[filter.Key]
		return exists
	case "not_exists":
		_, exists := userTags[filter.Key]
		return !exists
	case "=", "!=", "<", ">":
		userValue, exists := userTags[filter.Key]
		if !exists {
			return false
		}
		return compareValues(userValue, filter.Value, filter.Relation)
	default:
		fmt.Printf("Unsupported relation: %s\n", filter.Relation)
		return false
	}
}

// Compare the values by converting to integers and using basic operands
func compareValues(userValue, filterValue, relation string) bool {
	// Try to convert both values to integers for numeric comparisons
	userInt, userErr := strconv.Atoi(userValue)
	filterInt, filterErr := strconv.Atoi(filterValue)

	if userErr == nil && filterErr == nil {
		// Both values are integers, perform numeric comparison
		switch relation {
		case "=":
			return userInt == filterInt
		case "!=":
			return userInt != filterInt
		case "<":
			return userInt < filterInt
		case ">":
			return userInt > filterInt
		}
	} else {
		// At least one value is not an integer, perform string comparison
		switch relation {
		case "=":
			return userValue == filterValue
		case "!=":
			return userValue != filterValue
		case "<":
			return userValue < filterValue
		case ">":
			return userValue > filterValue
		}
	}
	return false
}
