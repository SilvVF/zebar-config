package main

import (
	"testing"
)

func TestMain(t *testing.T) {

	_, _, err := DailyNote(ZZZConfig)
	if err != nil {
		t.Error(err)
	}
	_, _, err = DailyNote(StarRailConfig)
	if err != nil {
		t.Error(err)
	}
	_, _, err = DailyNote(GenshinConfig)
	if err != nil {
		t.Error(err)
	}
}
