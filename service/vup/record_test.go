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

func ATestGetTopGuestRecords(t *testing.T) {
	be, err := GetTopGuestRecords(198297, 5)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, be)
}

func ATestGetGlobalRecords(t *testing.T) {
	be, err := GetGlobalRecords("", 1, 5, false)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, be)
}
