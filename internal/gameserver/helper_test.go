package gameserver

import (
	"slices"
	"strings"
	"testing"
)

func TestToContainerName_AlphaNumeric_Unchanged(t *testing.T) {
	sut := "AsF_1kl2o"
	result := ToContainerName(sut)

	if result != strings.ToLower(sut) {
		t.Fatalf("No change expected for test string '%s', got: %v", sut, result)
	}
}

func TestToContainerName_String_Changed(t *testing.T) {
	sut := "AsF-1kl2o/:l@0%&91"
	result := ToContainerName(sut)

	if strings.ContainsFunc(result, func(r rune) bool {
		return slices.Contains([]rune{'@', '%', ':'}, r)
	}) {
		t.Fatalf("Change expected for substring '%s', got: %v", "/:l@0%&91", result)
	}
}
