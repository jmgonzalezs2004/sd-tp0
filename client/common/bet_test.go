package common

import (
	"testing"
)

func TestBetFields(t *testing.T) {
	bet := Bet{
		FirstName: "Santiago Lionel",
		LastName:  "Lorca",
		Document:  "30904465",
		Birthdate: "1999-03-17",
		Number:    7574,
	}

	if bet.FirstName != "Santiago Lionel" {
		t.Errorf("FirstName = %v; want Santiago Lionel", bet.FirstName)
	}
	if bet.Document != "30904465" {
		t.Errorf("Document = %v; want 30904465", bet.Document)
	}
	if bet.Number != 7574 {
		t.Errorf("Number = %v; want 7574", bet.Number)
	}
}
