package idgen

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/fortifi/fident/logging"
)

const specTimeKeyLength = 9

// Test params
const testVendor = "FID"
const testType = "TE"
const testSubType = "ES"

// TestGenerationOfNewID tests that the full ID generation output is as expected
func TestGenerationOfNewID(t *testing.T) {
	inputTime := time.Now()
	logging.Write(logging.LevelTest, "Testing full ID generation")
	id, err := New(IndicatorEntity, testVendor, testType, testSubType, "")
	if err != nil || id == "" || len(id) != idLength {
		logging.Write(logging.LevelTestFailure, "Error generating ID")
		t.Fail()
	}

	re := regexp.MustCompile(idRegex)
	match := re.FindStringSubmatch(id)
	if len(match) != 1 {
		logging.Write(logging.LevelTestFailure, "Generated ID did not match expected format: "+id)
		t.Fail()
	}

	if result, err := Validate(id); result == false {
		logging.Write(logging.LevelTestFailure, "Generated ID was invalid: "+err.Error())
		t.Fail()
	}

	keyTime, err := GetTimeFromID(id)
	if err != nil {
		logging.Write(logging.LevelTestFailure, "Error with embedded timekey "+err.Error())
		t.Fail()
	}

	if keyTime.Day() != inputTime.Day() {
		logging.Write(logging.LevelTestFailure, "Timekey in ID is an unexpected value")
		t.Fail()
	}
}

// TestBase36KeyGeneration tests generation of the base36 time key
func TestBase36KeyGeneration(t *testing.T) {
	logging.Write(logging.LevelTest, "Testing base 36 key generation")
	key, err := getBase32TimeKey(time.Now())

	if err != nil {
		logging.Write(logging.LevelTestFailure, "Error generating ID time key")
		t.Fail()
	}

	if len(key) != specTimeKeyLength {
		logging.Write(logging.LevelTestFailure, "Time key generated was not of specification length")
		t.Fail()
	}
}

// TestRandomStringGeneration tests generation of the random string componenet
func TestRandomStringGeneration(t *testing.T) {
	logging.Write(logging.LevelTest, "Testing random string generation")
	randLength := 32
	strRan := getRandString(randLength)
	if len(strRan) != randLength {
		logging.Write(logging.LevelTestFailure, "Random string generated with unexpected length")
		t.Fail()
	}

	if strRan != strings.ToUpper(strRan) {
		logging.Write(logging.LevelTestFailure, "Random string contains unexpected characters")
		t.Fail()
	}
}

const validIndicator = IndicatorEntity
const invalidIndicator = TypeIndicator("Z")

// TestTypeIndicatorValidation tests that Type Indicator validation correctly identifies valid Indicators
func TestTypeIndicatorValidation(t *testing.T) {
	logging.Write(logging.LevelTest, "Testing type indicator validation")
	if isValidIndicator(string(validIndicator)) != true {
		logging.Write(logging.LevelTestFailure, "Valid type indicator flagged as invalid")
		t.Fail()
	}

	if isValidIndicator(string(invalidIndicator)) != false {
		logging.Write(logging.LevelTestFailure, "Invalid type indicator flagged as valid")
		t.Fail()
	}
}

const testValidID1 = "55QRHT4ET-EFORTKPS-USC1B-XRAQPPD"
const testValidID2 = "IP41IR5M=-EFIOUSAC-MISCR-9JQLO0F"
const testInvalidID = "55QRHT4E-EFOTKPS-USC1B-39H6POWT"

// TestIDValidation tests IDs validate correctly
func TestIDValidation(t *testing.T) {
	logging.Write(logging.LevelTest, "Testing ID Validation")
	if result, _ := Validate(testValidID1); result == false {
		logging.Write(logging.LevelTestFailure, "Valid ID failed to validate")
		t.Fail()
	}

	if result, _ := Validate(testValidID2); result == false {
		logging.Write(logging.LevelTestFailure, "Valid ID failed to validate")
		t.Fail()
	}

	if result, _ := Validate(testInvalidID); result == true {
		logging.Write(logging.LevelTestFailure, "Invalid ID did validate")
		t.Fail()
	}
}
