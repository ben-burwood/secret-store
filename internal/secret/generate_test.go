package secret

import (
	"strings"
	"testing"
)

func TestGenerateSecret_Length(t *testing.T) {
	secret := GenerateSecret(16, false, false)
	if len(secret) != 16 {
		t.Errorf("Expected length 16, got %d", len(secret))
	}
}

func TestGenerateSecret_OnlyLetters(t *testing.T) {
	secret := GenerateSecret(32, false, false)
	for _, c := range secret {
		if !strings.ContainsRune(Letters, c) {
			t.Errorf("Found non-letter character: %c", c)
		}
	}
}

func TestGenerateSecret_IncludeNumbers(t *testing.T) {
	secret := GenerateSecret(64, true, false)
	hasNumber := false
	for _, c := range secret {
		if strings.ContainsRune(Numbers, c) {
			hasNumber = true
			break
		}
	}
	if !hasNumber {
		t.Error("Expected at least one number in secret")
	}
}

func TestGenerateSecret_IncludeSymbols(t *testing.T) {
	secret := GenerateSecret(64, false, true)
	hasSymbol := false
	for _, c := range secret {
		if strings.ContainsRune(Symbols, c) {
			hasSymbol = true
			break
		}
	}
	if !hasSymbol {
		t.Error("Expected at least one symbol in secret")
	}
}

func TestGenerateSecret_IncludeNumbersAndSymbols(t *testing.T) {
	secret := GenerateSecret(128, true, true)
	hasNumber := false
	hasSymbol := false
	for _, c := range secret {
		if strings.ContainsRune(Numbers, c) {
			hasNumber = true
		}
		if strings.ContainsRune(Symbols, c) {
			hasSymbol = true
		}
	}
	if !hasNumber {
		t.Error("Expected at least one number in secret")
	}
	if !hasSymbol {
		t.Error("Expected at least one symbol in secret")
	}
}
