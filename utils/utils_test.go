package utils_test

import (
	"betme_test/proto"
	"betme_test/utils"
	"fmt"
	"testing"
)

type Testcase struct {
	Dates          *proto.Dates
	TestDate       *proto.Date
	ExpectedResult bool
}

func TestContainsDate(t *testing.T) {
	testcases := []Testcase{
		{Dates: &proto.Dates{Dates: []*proto.Date{{Year: 2022, Month: 2, Day: 7}, {Year: 2023, Month: 1, Day: 7}, {Year: 2023, Month: 2, Day: 6}}}, TestDate: &proto.Date{Year: 2023, Month: 2, Day: 7}, ExpectedResult: false},
		{Dates: &proto.Dates{Dates: []*proto.Date{{Year: 2023}, {Month: 2}, {Day: 7}}}, TestDate: &proto.Date{Year: 2023, Month: 2, Day: 7}, ExpectedResult: true},
	}
	for _, testcase := range testcases {
		actualResult := utils.ContainsDate(testcase.Dates, testcase.TestDate)
		if testcase.ExpectedResult != actualResult {
			fmt.Printf("Testcase failed: %v, Obtained result: %v\n", testcase, actualResult)
			t.Fail()
		}
	}
}
