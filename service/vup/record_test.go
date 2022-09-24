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

func aTestGetTopGuestRecords(t *testing.T) {
	be, err := GetTopGuestRecords(198297, 5)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, be)
}

func aTestGetTopGuestRecordsNF(t *testing.T) {
	be, err := GetTopGuestRecords(123456789, 5)
	if err != nil {
		t.Fatal(err)
	}
	if be == nil {
		t.Log("not found")
		return
	}
	jsonPrettyPrint(t, be)
}

func aTestGetGlobalRecords(t *testing.T) {
	be, err := GetGlobalRecords("", "", 1, 5, false)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, be)
}
