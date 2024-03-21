package guess

import (
	"fmt"
	"testing"
)

func Test_GetColorCode(t *testing.T) {
	testChar := CharacterType('k')

	result := testChar.GetColorCode(colorRed)
	if result != fmt.Sprintf("%s%c%s", colorRed, testChar, colorReset) {
		t.Fatalf("unexpected result from GetColorCode, got %s\n", result)
	}
}
