package utils

import (
	"betme_test/proto"
	"encoding/json"
	"log"
	"strconv"
	"strings"
)

// DataMatchesAnyDate - Returns if the fileData from feedName matches any of the dates from request
func DataMatchesAnyDate(feedName proto.Feed, fileData []byte, dates *proto.Dates) bool {

	// Extract Date from fileData for different feeds
	var thisDate *proto.Date
	if feedName == proto.Feed_FEED_X {
		thisDate = GetDateFromFeedXData(fileData)
	} else if feedName == proto.Feed_FEED_Y {
		thisDate = GetDateFromFeedYData(fileData)
	}

	// return false if failed to extract Date
	if thisDate == nil {
		log.Println("Error extracting date from fileData")
		return false
	}

	return ContainsDate(dates, thisDate)
}

func ContainsDate(dates *proto.Dates, thisDate *proto.Date) bool {
	// Check if extracted date matches any of the dates from input
	for _, needDate := range dates.Dates {
		if needDate.Year != 0 && needDate.Year != thisDate.Year {
			continue
		}
		if needDate.Month != 0 && needDate.Month != thisDate.Month {
			continue
		}
		if needDate.Day != 0 && needDate.Day != thisDate.Day {
			continue
		}

		// Match found!
		// log.Printf("%v matched with %v\n", needDate, thisDate)
		return true
	}

	// None of the dates matched!
	// log.Printf("No match between %v and %v\n", thisDate, dates)
	return false
}

// Get DateFromFeedXData - Parses the fileData from Feed X to extract the date
func GetDateFromFeedXData(fileData []byte) *proto.Date {
	var err error
	myMap := make(map[string]interface{})
	err = json.Unmarshal(fileData, &myMap)
	if err != nil {
		return nil
	}
	// fileData["data"]
	fileData, err = json.Marshal(myMap["data"])
	if err != nil {
		return nil
	}
	myMap = make(map[string]interface{})
	err = json.Unmarshal(fileData, &myMap)
	if err != nil {
		return nil
	}
	// fileData["data"]["time"]
	fileData, err = json.Marshal(myMap["time"])
	if err != nil {
		return nil
	}
	myMap = make(map[string]interface{})
	err = json.Unmarshal(fileData, &myMap)
	if err != nil {
		return nil
	}
	// fileData["data"]["time"]["starting_at"]
	fileData, err = json.Marshal(myMap["starting_at"])
	if err != nil {
		return nil
	}
	myMap = make(map[string]interface{})
	err = json.Unmarshal(fileData, &myMap)
	if err != nil {
		return nil
	}
	// fileData["data"]["time"]["starting_at"]["date"]
	date, err := json.Marshal(myMap["date"])
	if err != nil {
		return nil
	}
	dateStr := strings.Trim(string(date), "\"")
	parts := strings.Split(dateStr, "-")
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}
	return &proto.Date{
		Year:  int32(year),
		Month: int32(month),
		Day:   int32(day),
	}
}

// Get DateFromFeedYData - Parses the fileData from Feed Y to extract the date
func GetDateFromFeedYData(fileData []byte) *proto.Date {
	var err error
	// fileData[0]
	fileDataStr := strings.Trim(strings.Trim(string(fileData), "["), "]")
	myMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(fileDataStr), &myMap)
	if err != nil {
		return nil
	}
	// filedata[0]["fixture"]
	fileData, err = json.Marshal(myMap["fixture"])
	if err != nil {
		return nil
	}
	myMap = make(map[string]interface{})
	err = json.Unmarshal(fileData, &myMap)
	if err != nil {
		return nil
	}
	// fileData[0]["fixture"]["date"]
	date, err := json.Marshal(myMap["date"])
	if err != nil {
		return nil
	}
	dateStrs := strings.Split(strings.Trim(string(date), "\""), "T")
	if len(dateStrs) < 1 {
		return nil
	}
	parts := strings.Split(dateStrs[0], "-")
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil
	}
	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil
	}
	return &proto.Date{
		Year:  int32(year),
		Month: int32(month),
		Day:   int32(day),
	}
}
