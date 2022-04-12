package vup

import (
	"testing"
	"vup_dd_stats/service/blive"
)

func ATestGetVupStats(t *testing.T) {
	vup, err := GetStats(392505232, 5)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, vup)
}

func ATestGetVupStatsCommand(t *testing.T) {
	vup, err := GetStatsCommand(690608693, 5, blive.InteractWord)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, vup)
}

func aTestGetGlobalStats(t *testing.T) {
	vup, err := GetGlobalStats(5)
	if err != nil {
		t.Fatal(err)
	}
	jsonPrettyPrint(t, vup)
}
