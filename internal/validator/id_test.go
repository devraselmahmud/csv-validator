package validator

import "testing"

func TestIDTrackerRecognizesDuplicateWithSurroundingWhitespace(t *testing.T) {
	tracker := idTracker{}

	if tracker.isDuplicate("42") {
		t.Fatal("first ID reported as duplicate")
	}

	if !tracker.isDuplicate(" 42 ") {
		t.Fatal("ID with surrounding whitespace not reported as duplicate")
	}
}

func TestIDTrackerDoesNotTreatBlankIDsAsDuplicates(t *testing.T) {
	tracker := idTracker{}

	for _, id := range []string{"", " ", "\t", ""} {
		if tracker.isDuplicate(id) {
			t.Fatalf("isDuplicate(%q) = true for blank ID, want false", id)
		}
	}
}

func TestIDTrackerKeepsIDsCaseSensitive(t *testing.T) {
	tracker := idTracker{}

	for _, id := range []string{"ABC", "abc"} {
		if tracker.isDuplicate(id) {
			t.Fatalf("isDuplicate(%q) = true, want false", id)
		}
	}

	if !tracker.isDuplicate("ABC") {
		t.Fatal("repeated uppercase ID not reported as duplicate")
	}
}

func TestIDTrackersDoNotShareState(t *testing.T) {
	first := idTracker{}
	second := idTracker{}

	if first.isDuplicate("42") {
		t.Fatal("first tracker reported new ID as duplicate")
	}

	if second.isDuplicate("42") {
		t.Fatal("second tracker inherited first tracker's IDs")
	}
}
