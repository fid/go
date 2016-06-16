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
	testVendorSecret  = "bad_secret"

	testVendor  = "FID"
	testType    = "TE"
	testSubType = "ES"
)

// TestGenerationOfNewID tests that the full ID generation output is as expected
func TestGenerationOfNewID(t *testing.T) {
	inputTime := time.Now()
	t.Log("Testing full ID generation")
	id, err := Generate(IndicatorEntity, testVendor, testType, testSubType, "", "")
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

	keyTime, err := GetTimeFromID(id)
	if err != nil {
		t.Errorf("Error with embedded timekey " + err.Error())
	}

	if keyTime.Day() != inputTime.Day() {
		t.Errorf("Timekey in ID is an unexpected value")
	}
}

// TestGenerationOfNewID tests that the full ID generation output is as expected
func TestGenerationOfNewIDWithVendorSecret(t *testing.T) {
	inputTime := time.Now()
	t.Log("Testing full ID generation with vendor secret")
	id, err := Generate(IndicatorEntity, testVendor, testType, testSubType, "", testVendorSecret)
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
	if result, _ := Verify(testValidID1, ""); result == false {
		t.Errorf("Valid ID failed verification")
	}

	if result, _ := Verify(testValidID2, ""); result == false {
		t.Errorf("Valid ID failed verification")
	}

	if result, _ := Verify(testInvalidID, ""); result == true {
		t.Errorf("Invalid ID did verify")
	}
}

//
func TestValidationFromExternalSource(t *testing.T) {
	fids := []string{
		"IPIH7MHX=-EABCCDEF-MISCR-OI35356", "IPIH7MHX=-EABCCDEF-MISCR-998CT0G", "IPIH7MHX=-EABCCDEF-MISCR-36L8W35", "IPIH7MHX=-EABCCDEF-MISCR-W8H41AF", "IPIH7MHX=-EABCCDEF-MISCR-2XV3I8C",
		"IPIH7MHX=-EABCCDEF-MISCR-5M6YEOB", "IPIH7MHY=-EABCCDEF-MISCR-5HP9E6E", "IPIH7MHY=-EABCCDEF-MISCR-XYKO01L", "IPIH7MHY=-EABCCDEF-MISCR-31XSKHU", "IPIH7MHY=-EABCCDEF-MISCR-LG1BO30",
		"IPIH7MHY=-EABCCDEF-MISCR-7UTZ5AU", "IPIH7MHY=-EABCCDEF-MISCR-5LSMHK4", "IPIH7MHY=-EABCCDEF-MISCR-65B511L", "IPIH7MHY=-EABCCDEF-MISCR-D73BF0O", "IPIH7MHY=-EABCCDEF-MISCR-UG87F5K",
		"IPIH7MHY=-EABCCDEF-MISCR-L3X1597", "IPIH7MHY=-EABCCDEF-MISCR-8CVN055", "IPIH7MHY=-EABCCDEF-MISCR-959015G", "IPIH7MHY=-EABCCDEF-MISCR-NWCHDXX", "IPIH7MHY=-EABCCDEF-MISCR-N4BOT21",
		"IPIH7MHY=-EABCCDEF-MISCR-WOG4HMV", "IPIH7MHZ=-EABCCDEF-MISCR-X321Y7Q", "IPIH7MHZ=-EABCCDEF-MISCR-9T9RB54", "IPIH7MHZ=-EABCCDEF-MISCR-13YA20M", "IPIH7MHZ=-EABCCDEF-MISCR-YZYD4U8",
		"IPIH7MHZ=-EABCCDEF-MISCR-O6W596U", "IPIH7MHZ=-EABCCDEF-MISCR-IO9J865", "IPIH7MHZ=-EABCCDEF-MISCR-E4E3E9Q", "IPIH7MHZ=-EABCCDEF-MISCR-6F12P44", "IPIH7MHZ=-EABCCDEF-MISCR-LCIVXS7",
		"IPIH7MHZ=-EABCCDEF-MISCR-H38US57", "IPIH7MHZ=-EABCCDEF-MISCR-78OTQ5D", "IPIH7MHZ=-EABCCDEF-MISCR-L9O1Z54", "IPIH7MHZ=-EABCCDEF-MISCR-60K27D3", "IPIH7MHZ=-EABCCDEF-MISCR-W03Q083",
		"IPIH7MHZ=-EABCCDEF-MISCR-U551677", "IPIH7MHZ=-EABCCDEF-MISCR-2WJT077", "IPIH7MI0=-EABCCDEF-MISCR-Z9J32OZ", "IPIH7MI0=-EABCCDEF-MISCR-0HGH03S", "IPIH7MI0=-EABCCDEF-MISCR-5R5514M",
		"IPIH7MI0=-EABCCDEF-MISCR-54DP0OA", "IPIH7MI0=-EABCCDEF-MISCR-3L08BXQ", "IPIH7MI0=-EABCCDEF-MISCR-E9AIJHS", "IPIH7MI0=-EABCCDEF-MISCR-YTH5E31", "IPIH7MI0=-EABCCDEF-MISCR-H7K339V",
		"IPIH7MI0=-EABCCDEF-MISCR-6Y32244", "IPIH7MI0=-EABCCDEF-MISCR-V904R12", "IPIH7MI0=-EABCCDEF-MISCR-OVX87IB", "IPIH7MI0=-EABCCDEF-MISCR-0Z6964F", "IPIH7MI0=-EABCCDEF-MISCR-4GNAB68",
		"IPIH7MI0=-EABCCDEF-MISCR-9PN7C1U", "IPIH7MI0=-EABCCDEF-MISCR-OD58P04", "IPIH7MI1=-EABCCDEF-MISCR-N7E56UD", "IPIH7MI1=-EABCCDEF-MISCR-QNEFD02", "IPIH7MI1=-EABCCDEF-MISCR-V1V6X3Q",
		"IPIH7MI1=-EABCCDEF-MISCR-9Y9F776", "IPIH7MI1=-EABCCDEF-MISCR-T984954", "IPIH7MI1=-EABCCDEF-MISCR-WG9H07J", "IPIH7MI1=-EABCCDEF-MISCR-NMP09LB", "IPIH7MI1=-EABCCDEF-MISCR-3182KAP",
		"IPIH7MI1=-EABCCDEF-MISCR-T1TA0AO", "IPIH7MI1=-EABCCDEF-MISCR-6EDO648", "IPIH7MI1=-EABCCDEF-MISCR-OHLP7MS", "IPIH7MI1=-EABCCDEF-MISCR-L6W1I9Y", "IPIH7MI1=-EABCCDEF-MISCR-0Y65485",
		"IPIH7MI1=-EABCCDEF-MISCR-Z0AP5DW", "IPIH7MI1=-EABCCDEF-MISCR-6I3SLYL", "IPIH7MI1=-EABCCDEF-MISCR-49FS6VS", "IPIH7MI2=-EABCCDEF-MISCR-FG2SKZV", "IPIH7MI2=-EABCCDEF-MISCR-3H3260J",
		"IPIH7MI2=-EABCCDEF-MISCR-2X2D6GH", "IPIH7MI2=-EABCCDEF-MISCR-667475Y", "IPIH7MI2=-EABCCDEF-MISCR-XE7335J", "IPIH7MI2=-EABCCDEF-MISCR-9Z6IZ9M", "IPIH7MI2=-EABCCDEF-MISCR-TR54U64",
		"IPIH7MI2=-EABCCDEF-MISCR-170O077", "IPIH7MI2=-EABCCDEF-MISCR-02T033P", "IPIH7MI2=-EABCCDEF-MISCR-V669VFQ", "IPIH7MI2=-EABCCDEF-MISCR-FHZ8818", "IPIH7MI2=-EABCCDEF-MISCR-90ADF6T",
		"IPIH7MI2=-EABCCDEF-MISCR-536R677", "IPIH7MI2=-EABCCDEF-MISCR-9440F75", "IPIH7MI2=-EABCCDEF-MISCR-6W7BFET", "IPIH7MI2=-EABCCDEF-MISCR-ZI5COK2", "IPIH7MI3=-EABCCDEF-MISCR-H7KV320",
		"IPIH7MI3=-EABCCDEF-MISCR-NW3076W", "IPIH7MI3=-EABCCDEF-MISCR-S8D0JFB", "IPIH7MI3=-EABCCDEF-MISCR-9B79531", "IPIH7MI3=-EABCCDEF-MISCR-33SAI39", "IPIH7MI3=-EABCCDEF-MISCR-IE9900X",
		"IPIH7MI3=-EABCCDEF-MISCR-94T912V", "IPIH7MI3=-EABCCDEF-MISCR-B2H2GSC"}

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

// Benchmark single fid generation
func BenchmarkSingleFidGeneration(b *testing.B) {
	for n := 0; n < 20; n++ {
		Generate(IndicatorEntity, testVendor, testType, testSubType, "", "")
	}
}

// Benchmark to generate 1 million FIDs (12.7 seconds on mid 2015 Macbook pro)
func BenchmarkOneMillionFidGeneration(b *testing.B) {
	for i := 0; i < 1000000; i++ {
		Generate(IndicatorEntity, testVendor, testType, testSubType, "", "")
	}
}
