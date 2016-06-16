package gofid

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

// Test params
const (
	specTimeKeyLength = 9

	testVendor  = "FID"
	testType    = "TE"
	testSubType = "ES"
)

// TestGenerationOfNewID tests that the full ID generation output is as expected
func TestGenerationOfNewID(t *testing.T) {
	inputTime := time.Now()
	t.Log("Testing full ID generation")
	id, err := Generate(IndicatorEntity, testVendor, testType, testSubType, "")
	if err != nil || id == "" || len(id) != idLength {
		t.Errorf("Error generating ID")
	}

	re := regexp.MustCompile(idRegex)
	match := re.FindStringSubmatch(id)
	if len(match) != 1 {
		t.Errorf("Generated ID did not match expected format: " + id)
	}

	if result, err := Verify(id); result == false {
		t.Errorf("Generated ID was invalid: " + err.Error())
	}

	keyTime, err := GetTimeFromID(id)
	if err != nil {
		t.Errorf("Error with embedded timekey " + err.Error())
	}

	if keyTime.Day() != inputTime.Day() {
		t.Errorf("Timekey in ID is an unexpected value")
	}
}

// TestBase36KeyGeneration tests generation of the base36 time key
func TestBase36KeyGeneration(t *testing.T) {
	t.Log("Testing base 36 key generation")
	key, err := getBase32TimeKey(time.Now())

	if err != nil {
		t.Errorf("Error generating ID time key")
	}

	if len(key) != specTimeKeyLength {
		t.Errorf("Time key generated was not of specification length")
	}
}

// TestRandomStringGeneration tests generation of the random string componenet
func TestRandomStringGeneration(t *testing.T) {
	t.Log("Testing random string generation")
	randLength := 32
	strRan := getRandString(randLength)
	if len(strRan) != randLength {
		t.Errorf("Random string generated with unexpected length")
	}

	if strRan != strings.ToUpper(strRan) {
		t.Errorf("Random string contains unexpected characters")
	}
}

const validIndicator = IndicatorEntity
const invalidIndicator = TypeIndicator("Z")

// TestTypeIndicatorValidation tests that Type Indicator validation correctly identifies valid Indicators
func TestTypeIndicatorValidation(t *testing.T) {
	t.Log("Testing type indicator validation")
	if isValidIndicator(string(validIndicator)) != true {
		t.Errorf("Valid type indicator flagged as invalid")
	}

	if isValidIndicator(string(invalidIndicator)) != false {
		t.Errorf("Invalid type indicator flagged as valid")
	}
}

const testValidID1 = "55QRHT4ET-EFORTKPS-USC1B-XRAQPPD"
const testValidID2 = "IP41IR5M=-EFIOUSAC-MISCR-9JQLO0F"
const testInvalidID = "55QRHT4E-EFOTKPS-USC1B-39H6POWT"

// TestIDValidation tests IDs validate correctly
func TestIDValidation(t *testing.T) {
	t.Log("Testing ID Validation")
	if result, _ := Verify(testValidID1); result == false {
		t.Errorf("Valid ID failed verification")
	}

	if result, _ := Verify(testValidID2); result == false {
		t.Errorf("Valid ID failed verification")
	}

	if result, _ := Verify(testInvalidID); result == true {
		t.Errorf("Invalid ID did verify")
	}
}
