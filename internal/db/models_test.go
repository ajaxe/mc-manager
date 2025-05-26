package db

import (
	"encoding/base64"
	"testing"
	"time"
)

func TestEncodePaginationCursor(t *testing.T) {
	id := "12345"
	dt := time.Now().Format(time.RFC3339)

	expected := base64.URLEncoding.EncodeToString([]byte(id + "|" + dt))

	result := EncodePaginationCursor(&PaginationCursor{
		ID:         id,
		LaunchDate: dt,
	})

	if result != expected {
		t.Fatalf("error encoding cursor. expected:%v, got:%v", expected, result)
	}
}

func TestDecodePaginationCursor(t *testing.T) {
	id := "12345"
	dt := time.Now().Format(time.RFC3339)

	coded := base64.URLEncoding.EncodeToString([]byte(id + "|" + dt))

	result, e := DecodePaginationCursor(coded)

	if e != nil {
		t.Fatalf("failed to decode cursor: %v", e)
	}

	if result.ID != id {
		t.Fatalf("invlaid id, expected: %v, got:%v", id, result.ID)
	}
	if result.LaunchDate != dt {
		t.Fatalf("invlaid id, expected: %v, got:%v", dt, result.LaunchDate)
	}
}
