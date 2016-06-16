package idgen

/**
* Fident uses the Fortifi Open ID structure specification (See 'ID Structure' documentation)
**/

import (
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

// New returns a newly generated ID in Fortifi Open ID format
func New(systemIndicator TypeIndicator, vendor, nType, nSubType, priLocation string) (string, error) {
	timeKey, err := getBase32TimeKey(time.Now())
	if err != nil {
		return "", err
	}

	// System indicator (Replace unknown with 'entity')
	if !isValidIndicator(string(systemIndicator)) || len(systemIndicator) == 0 {
		systemIndicator = IndicatorEntity
	}

	if len(vendor) != 3 {
		return "", errors.New("Vendor must be of length '3'")
	}

	if len(nType) != 2 {
		return "", errors.New("Type must be of length '2'")
	}

	if len(nSubType) != 2 {
		nSubType = nType
	}

	if len(priLocation) != 5 {
		priLocation = unknownLocationValue
	}

	randomString := getRandString(randLen)
	return strings.ToUpper(fmt.Sprintf("%s%s%s%s%s%s%s%s%s%s", timeKey, delimitChar, systemIndicator, vendor, nType, nSubType, delimitChar, priLocation, delimitChar, randomString)), nil
}

// Validate is true if string is a valid Fortifi Open ID
func Validate(id string) (bool, error) {
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

	return true, nil
}

// GetTimeFromID Returns the time from the timekey embedded in ID
func GetTimeFromID(id string) (time.Time, error) {
	validate, err := Validate(id)
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
