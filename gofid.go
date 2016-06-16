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

const (
	vendorLength         = 3
	typeElementLength    = 2
	priLocationLength    = 5
	idLength             = 32
	timeKeyBase          = 36
	randLen              = 7
	paddingChar          = "="
	delimitChar          = "-"
	idElements           = 4
	timeKeyLength        = 9
	unknownLocationValue = "MISCR"
	letterBytes          = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	idRegex              = "[A-Z0-9=]{9}-[A-Z0-9=]{8}-[A-Z0-9=]{5}-[A-Z0-9=]{7}\\z"

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
func Generate(systemIndicator TypeIndicator, vendor, nType, nSubType, priLocation, vendorSecret string) (string, error) {
	timeKey, err := getBase32TimeKey(time.Now())
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

	if len(nType) != typeElementLength {
		return "", fmt.Errorf("Type must be of length '%d'", typeElementLength)
	}

	if len(nSubType) != typeElementLength {
		nSubType = nType
	}

	if len(priLocation) != priLocationLength {
		priLocation = unknownLocationValue
	}

	randomString := getRandString(randLen)

	result := ""
	preResult := strings.ToUpper(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s", timeKey, delimitChar, systemIndicator, vendor, nType, nSubType, delimitChar, priLocation, delimitChar, randomString))

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

// GetTimeFromID Returns the time from the timekey embedded in ID
func GetTimeFromID(id string) (time.Time, error) {
	validate, err := Verify(id, "")
	if validate != true {
		return time.Time{}, err
	}

	components := strings.Split(id, delimitChar)
	miliseconds := components[0]
	miliseconds = strings.Replace(miliseconds, paddingChar, "", -1)
	msInt, err := strconv.ParseInt(miliseconds, 36, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(0, msInt*int64(time.Millisecond)), nil
}

// isValidIndicator checks that proposed indicator is valid as per spec
func isValidIndicator(proposed string) bool {
	proposed = strings.ToUpper(proposed)
	return strings.Contains(string(IndicatorChecksum), proposed)
}

// getBase32TimeKey returns a millisecond timestamp in base 36
func getBase32TimeKey(time time.Time) (string, error) {
	nanoTime := time.UnixNano()
	miliTime := nanoTime / 1000000
	timeKey := strings.ToUpper(strconv.FormatInt(miliTime, timeKeyBase))
	paddingLen := ((len(timeKey) - timeKeyLength) * -1)

	if paddingLen > 0 {
		for index := 0; index < paddingLen; index++ {
			timeKey = timeKey + paddingChar
		}
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
