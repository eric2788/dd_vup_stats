package vup

import "testing"

func ATestGetTopDDRecords(t *testing.T) {
	be, err := GetTopDDRecords(198297, 5)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, be)
}

func ATestGetTopSelfRecords(t *testing.T) {
	be, err := GetTopSelfRecords(198297, 5)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, be)
}
