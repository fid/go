package gofid

import (
	"regexp"
	"strings"
	"testing"
	"time"
)

const (
	specTimeKeyLength = 9
	testVendorSecret  = "bad_secret"
	testVendor        = "FID"
	testApp           = "TE"
	testType          = "ES"
)

// TestGenerationOfNewID tests that the full ID generation output is as expected
func TestGenerationOfNewID(t *testing.T) {
	inputTime := time.Now()
	t.Log("Testing full ID generation")
	id, err := Generate(IndicatorEntity, testVendor, testApp, testType, "", "")
	if err != nil || id == "" || len(id) != idLength {
		t.Errorf("Error generating ID")
	}

	re := regexp.MustCompile(idRegex)
	match := re.FindStringSubmatch(id)
	if len(match) != 1 {
		t.Errorf("Generated ID did not match expected format: " + id)
	}

	if result, err := Verify(id, ""); result == false {
		t.Errorf("Generated ID was invalid: " + err.Error())
	}

	description, err := Describe(id)
	if err != nil {
		t.Errorf("Error with description " + err.Error())
	}

	keyTime := description.Time

	if keyTime.Day() != inputTime.Day() {
		t.Errorf("Timekey in ID is an unexpected value %d : %d", keyTime.Unix(), inputTime.Unix())
	}
}

// TestGenerationOfNewID tests that the full ID generation output is as expected
func TestGenerationOfNewIDWithVendorSecret(t *testing.T) {
	inputTime := time.Now()
	t.Log("Testing full ID generation with vendor secret")
	id, err := Generate(IndicatorEntity, testVendor, testApp, testType, "", testVendorSecret)
	if err != nil || id == "" || len(id) != idLength {
		t.Errorf("Error generating ID")
	}

	re := regexp.MustCompile(idRegex)
	match := re.FindStringSubmatch(id)
	if len(match) != 1 {
		t.Errorf("Generated ID did not match expected format: " + id)
	}

	if result, err := Verify(id, testVendorSecret); result == false {
		t.Errorf("Generated ID was invalid: " + err.Error())
	}

	description, err := Describe(id)
	if err != nil {
		t.Errorf("Error with description " + err.Error())
	}

	keyTime := description.Time

	if keyTime.Day() != inputTime.Day() {
		t.Errorf("Timekey in ID is an unexpected value")
	}

	if result, _ := Verify(id, testVendorSecret+"qweqwe"); result == true {
		t.Errorf("ID passed verification with invalid vendor secret")
	}
}

// TestBase36KeyGeneration tests generation of the base36 time key
func TestBase36KeyGeneration(t *testing.T) {
	t.Log("Testing base 36 key generation")
	key, err := getBase36TimeKey(time.Now())

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

const testValidID1 = "EFORTKPS-55QRHT4ET-USC1B-XRAQPPD"
const testValidID2 = "EFIOUSAC-IP41IR5M=-MISCR-9JQLO0F"
const testInvalidID = "EFOTKPS-55QRHT4E-USC1B-39H6POWT"

// TestIDValidation tests IDs validate correctly
func TestIDValidation(t *testing.T) {
	t.Log("Testing ID Validation")
	if result, err := Verify(testValidID1, ""); result == false {
		t.Errorf("Valid ID failed verification" + err.Error())
	}

	if result, err := Verify(testValidID2, ""); result == false {
		t.Errorf("Valid ID failed verification" + err.Error())
	}

	if result, err := Verify(testInvalidID, ""); result == true {
		t.Errorf("Invalid ID did verify" + err.Error())
	}
}

// TestValidationFromExternalSource tests externally generated IDs
func TestValidationFromExternalSource(t *testing.T) {
	fids := []string{
		"EABCCDEF-IPIH7MHX=-MISCR-OI35356"}

	for _, fid := range fids {
		result, err := Verify(fid, "")
		if result == false {
			t.Errorf("Input fid failed verification %s", fid)
		}

		if err != nil {
			t.Errorf("Error when validating fid")
		}
	}
}

// TestGetDescription generates a description and compares it against input data
func TestGetDescription(t *testing.T) {
	timeKey := "55QRHT4ET"
	systemIndicator := IndicatorLog
	vendorKey := "FOR"
	app := "TE"
	ntype := "ST"
	location := "USC1B"
	rand := "XRAQPPD"
	inputID := string(systemIndicator) + vendorKey + app + ntype + delimitChar + timeKey + delimitChar + location + delimitChar + rand

	description, err := Describe(inputID)
	if err != nil {
		t.Errorf("Error when describing FID")
	}

	if description.TimeKey != timeKey {
		t.Errorf("Description time key %s does not match input %s", description.TimeKey, timeKey)
	}

	if description.Indicator != systemIndicator {
		t.Errorf("Description indicator %s does not match input %s", description.Indicator, systemIndicator)
	}

	if description.VendorKey != vendorKey {
		t.Errorf("Description vendor key %s does not match input %s", description.VendorKey, vendorKey)
	}

	if description.App != app {
		t.Errorf("Description app %s does not match input %s", description.Type, app)
	}

	if description.Type != ntype {
		t.Errorf("Description type %s does not match input %s", description.Type, ntype)
	}

	if description.Location != location {
		t.Errorf("Description location %s does not match input %s", description.Location, location)
	}

	if description.RandomString != rand {
		t.Errorf("Description random string %s does not match input %s", description.TimeKey, timeKey)
	}
}

// Benchmark single fid generation
func BenchmarkSingleFidGeneration(b *testing.B) {
	for n := 0; n < 20; n++ {
		Generate(IndicatorEntity, testVendor, testApp, testType, "", "")
	}
}

// Benchmark to generate 1 million FIDs (12.7 seconds on mid 2015 Macbook pro)
func BenchmarkOneMillionFidGeneration(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		Generate(IndicatorEntity, testVendor, testApp, testType, "", "")
	}
}
