package gofid

/**
* Gofid: Fortifi Open ID
**/

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TypeIndicator system indicators
type TypeIndicator string

// Description of ID input to 'Describe' function
type Description struct {
	Indicator    TypeIndicator
	VendorKey    string
	App          string
	Type         string
	Location     string
	TimeKey      string
	Time         time.Time
	RandomString string
}

const (
	vendorLength          = 3
	appElementLength      = 2
	typeElementLength     = 2
	priLocationLength     = 5
	idLength              = 32
	timeKeyBase           = 36
	randLen               = 7
	delimitChar           = "-"
	idElements            = 4
	timeKeyLength         = 9
	maxTimestampValBase10 = 101559956668415 // Base 10 representation of Max value of FID Base 32 timestamp
	unknownLocationValue  = "MISCR"
	letterBytes           = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	idRegex               = "[A-Z0-9=]{8}-[A-Z0-9=]{9}-[A-Z0-9=]{5}-[A-Z0-9=]{7}\\z"

	/**
	 * System type indicators
	 **/

	// IndicatorEntity system indicator for Standard Entities
	IndicatorEntity TypeIndicator = "E"

	// IndicatorLog system indicator for Log Entries
	IndicatorLog TypeIndicator = "L"

	// IndicatorCache system indicator for Cached Data
	IndicatorCache TypeIndicator = "C"

	// IndicatorMemory system indicator for In Memory Item
	IndicatorMemory TypeIndicator = "M"

	// IndicatorMeta system indicator for Meta Data
	IndicatorMeta TypeIndicator = "A"

	// IndicatorConfiguration system indicator for Configuration Item
	IndicatorConfiguration TypeIndicator = "F"

	// IndicatorTimeSeries system indicator for Time Series Data
	IndicatorTimeSeries TypeIndicator = "T"

	// IndicatorRelationship system indicator for Relationship
	IndicatorRelationship TypeIndicator = "R"

	// IndicatorNote system indicator for Note / Comment
	IndicatorNote TypeIndicator = "N"

	// IndicatorFile system indicator for File (likely stored in a bucket)
	IndicatorFile TypeIndicator = "D"

	// IndicatorChecksum is used to verify input is a valid indicator
	IndicatorChecksum = IndicatorEntity + IndicatorLog + IndicatorCache + IndicatorMemory + IndicatorMeta + IndicatorConfiguration +
		IndicatorTimeSeries + IndicatorTimeSeries + IndicatorRelationship + IndicatorNote + IndicatorFile
)

// Generate returns a new ID in Fortifi Open ID format
func Generate(systemIndicator TypeIndicator, vendor, app, nType, priLocation, vendorSecret string) (string, error) {
	timeKey, err := getBase36TimeKey(time.Now())
	if err != nil {
		return "", err
	}

	// System indicator (Replace unknown with 'entity')
	if !isValidIndicator(string(systemIndicator)) || len(systemIndicator) == 0 {
		systemIndicator = IndicatorEntity
	}

	if len(vendor) != vendorLength {
		return "", fmt.Errorf("Vendor must be of length '%d'", vendorLength)
	}

	if len(app) != typeElementLength {
		return "", fmt.Errorf("App must be of length '%d'", appElementLength)
	}

	if len(nType) != typeElementLength {
		return "", fmt.Errorf("Type must be of length '%d'", typeElementLength)
	}

	if len(priLocation) != priLocationLength {
		priLocation = unknownLocationValue
	}

	randomString := getRandString(randLen)
	result := ""
	preResult := strings.ToUpper(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s", systemIndicator, vendor, app, nType, delimitChar, timeKey, delimitChar, priLocation, delimitChar, randomString))

	if len(vendorSecret) > 0 {
		preResult = preResult[:len(preResult)-1]
		h := md5.New()
		h.Write([]byte(vendorSecret + preResult))
		hexEncoding := hex.EncodeToString(h.Sum(nil))
		result = preResult + strings.ToUpper(string(hexEncoding[0]))
	} else {
		result = preResult
	}

	return result, nil
}

// Verify is true if string is a valid Fortifi Open ID
func Verify(id, vendorSecret string) (bool, error) {
	if len(id) != idLength {
		return false, errors.New("ID is of invalid length")
	}

	re := regexp.MustCompile(idRegex)
	match := re.FindStringSubmatch(id)
	if len(match) != 1 {
		return false, errors.New("ID format is invalid")
	}

	components := strings.Split(id, delimitChar)
	if len(components) != idElements {
		return false, errors.New("Unexpected element count in ID")
	}

	if vendorSecret != "" {
		checkChar := string(id[len(id)-1:])
		idMinusCS := string(id[0:(len(id) - 1)])
		h := md5.New()
		h.Write([]byte(vendorSecret + idMinusCS))
		hexEncoding := hex.EncodeToString(h.Sum(nil))

		if strings.ToUpper(string(hexEncoding[0])) != checkChar {
			return false, errors.New("Checksum does not match vendor secret")
		}
	}

	return true, nil
}

// Describe returns decoded description object for the ID
func Describe(id string) (Description, error) {
	_, err := Verify(id, "")
	if err != nil {
		return Description{}, err
	}

	components := strings.Split(id, delimitChar)

	indicatorCom := components[0]
	sysIndicator := TypeIndicator(indicatorCom[0:1])
	vendorKey := indicatorCom[1:4]
	app := indicatorCom[4:6]
	nType := indicatorCom[6:8]
	location := components[2]
	timeKey := components[1]
	randStr := components[3]

	time, err := getTimeFromID(id)
	if err != nil {
		return Description{}, err
	}

	result := Description{
		Indicator:    sysIndicator,
		VendorKey:    vendorKey,
		App:          app,
		Type:         nType,
		Location:     location,
		TimeKey:      timeKey,
		Time:         time,
		RandomString: randStr,
	}

	return result, nil
}

// getTimeFromID Returns the time from the timekey embedded in ID
func getTimeFromID(id string) (time.Time, error) {
	validate, err := Verify(id, "")
	if validate != true {
		return time.Time{}, err
	}

	components := strings.Split(id, delimitChar)
	miliseconds := components[1]
	msInt, err := strconv.ParseInt(miliseconds, 36, 64)
	if err != nil {
		return time.Time{}, err
	}

	msSinceEpoch := maxTimestampValBase10 - msInt
	revMs := (msSinceEpoch * int64(time.Millisecond))
	return time.Unix(0, revMs), nil
}

// isValidIndicator checks that proposed indicator is valid as per spec
func isValidIndicator(proposed string) bool {
	proposed = strings.ToUpper(proposed)
	return strings.Contains(string(IndicatorChecksum), proposed)
}

// getBase36TimeKey returns a millisecond timestamp in base 36
func getBase36TimeKey(time time.Time) (string, error) {
	nanoTime := time.UnixNano()
	miliTime := nanoTime / 1000000
	revMiliTime := maxTimestampValBase10 - miliTime
	timeKey := strings.ToUpper(strconv.FormatInt(revMiliTime, timeKeyBase))
	paddingLen := timeKeyLength - len(timeKey)

	if paddingLen > 0 {
		timeKey = strings.Repeat("0", paddingLen) + timeKey
	}

	if paddingLen < 0 {
		return "", errors.New("Invalid time key")
	}

	return timeKey, nil
}

// generates a pseudorandom string
func getRandString(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	randBytes := make([]byte, n)
	for i := range randBytes {
		randBytes[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(randBytes)
}
