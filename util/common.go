package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	"github.com/jinzhu/now"
)

const (
	XHasuraRole       = "x-hasura-role"
	XHasuraUserID     = "x-hasura-user-id"
	XHasuraUserEmail  = "x-hasura-user-email"
	XHasuraIsSAMLUser = "x-hasura-user-saml"
	// Returned only for metrics role
	XHasuraAllowedProjectIDs        = "x-hasura-allowed-project-ids"
	XHasuraAllowedMetricsProjectIDs = "x-hasura-allowed-metrics-project-Ids"
	XHasuraAdminProjectIDs          = "x-hasura-admin-project-ids"
	// Returned for licensing public role
	Authorization     = "Authorization"
	SUBSCRIPTION_DATE = 1
)

// StringMap wrapper for string map
type StringMap map[string]string

// Get get string value from map by key
func (sm StringMap) Get(key string) string {
	if str, ok := sm[key]; ok {
		return str
	}

	return ""
}

// PostgresArrayToStrings convert postgres array to string array
func PostgresArrayToStrings(s string) ([]string, error) {
	if s == "" {
		return []string{}, nil
	}

	if s[0] != '{' || s[len(s)-1] != '}' {
		return nil, errors.New("invalid Postgres array: " + s)
	}

	if strings.TrimSpace(s[1:len(s)-1]) == "" {
		return []string{}, nil
	}

	result := strings.Split(s[1:len(s)-1], ",")
	return result, nil
}

// MapStringToUUID parse array strings to uuids
func MapStringToUUID(arr []string) ([]uuid.UUID, error) {
	results := make([]uuid.UUID, len(arr))

	for i, s := range arr {
		v, err := uuid.Parse(s)
		if err != nil {
			return nil, fmt.Errorf("parse uuid error %w: %q", err, s)
		}

		results[i] = v
	}

	return results, nil
}

// StringList stores string values
type StringList []string

// Add append new value or skip if it's existing
func (sl *StringList) Add(values ...string) {
	for _, s := range values {
		*sl = append(*sl, s)
	}
}

// IsEmpty check if the array is empty
func (sl StringList) IsEmpty() bool {
	return len(sl) == 0
}

// String implements the Stringer interface
func (sl StringList) String() string {
	var nonEmpty []string
	for _, s := range sl {
		if s != "" {
			nonEmpty = append(nonEmpty, s)
		}
	}
	return strings.Join(nonEmpty, ",")
}

// Compact returns a copy of the string list with empty elements removed
func (sl StringList) Compact() StringList {
	nsl := StringList{}
	for _, s := range sl {
		if s != "" {
			nsl = append(nsl, s)
		}
	}
	return nsl
}

// UniqueString creates a string of comma separated
// unique values sorted alphabetically
func (sl StringList) UniqueString() string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range sl {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	sort.Strings(list)
	return strings.Join(list, ",")
}

// ToJson serializes the StringList into Json
func (sl StringList) ToJson() []byte {
	json, _ := json.Marshal(sl)
	return json
}

// JsonList stores strings representing json
type JsonList []string

// Add append new value or skip if it's existing
func (jl *JsonList) Add(values ...string) {
	for _, s := range values {
		*jl = append(*jl, s)
	}
}

// IsEmpty check if the array is empty
func (jl JsonList) IsEmpty() bool {
	return len(jl) == 0
}

// ToJson serializes the JSON List into a valid JSON array. Elements are not
// marshalled, because they are json already
func (jl JsonList) ToJson() []byte {
	var sb strings.Builder
	sb.WriteRune('[')

	i := 0
	for ; i < len(jl)-1; i++ {
		sb.WriteString(jl[i])
		sb.WriteRune(',')
	}

	if i < len(jl) {
		sb.WriteString(jl[i])
	}

	sb.WriteRune(']')
	return []byte(sb.String())
}

// ToJsonSingleElementOrArray serializes a list containing a single element into
// that element's JSON; and into an array if it contains multiple elements.
func (jl JsonList) ToJsonSingleElementOrArray() []byte {
	if len(jl) == 1 {
		return []byte(jl[0])
	} else {
		return jl.ToJson()
	}
}

// Compact returns a copy of the JSON list with empty elements removed
func (jl JsonList) Compact() JsonList {
	nsl := JsonList{}
	for _, s := range jl {
		if s != "" {
			nsl = append(nsl, s)
		}
	}
	return nsl
}

func BeginningOfMonth() time.Time {
	now := time.Now().UTC()
	y, m, _ := now.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, now.Location())
}

func EndOfMonth() time.Time {
	return BeginningOfMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

func NoOfDaysInLastMonth() int {
	return EndOfLastMonth().Day()
}

// GetBeginningOfNextMonth returns the date at which the new billing cycle kicks in
func GetBeginningOfNextMonth() time.Time {
	return now.EndOfMonth().AddDate(0, 0, SUBSCRIPTION_DATE)
}

func GetBeginningOfThisMonth() time.Time {
	return now.BeginningOfMonth().AddDate(0, 0, SUBSCRIPTION_DATE)
}

func GetNextMonthString() string {
	return GetBeginningOfNextMonth().String()
}

func GetCurrentMonthString() string {
	return GetBeginningOfThisMonth().String()
}

func BeginningOfLastMonth() time.Time {
	now := time.Now().UTC()
	y, m, _ := now.Date()
	return time.Date(y, m-1, 1, 0, 0, 0, 0, now.Location())
}

func EndOfLastMonth() time.Time {
	return BeginningOfLastMonth().AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// BeginningOfNextMonth returns the absolute beginning of the next month
func BeginningOfNextMonth() time.Time {
	now := time.Now().UTC()
	y, m, _ := now.Date()
	return time.Date(y, m+1, 1, 0, 0, 0, 0, now.Location())
}

func GetTime() *string {
	b := time.Now().UTC().Format("2006-01-02T15:04:05-0700")
	return &b
}

func GetEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func NoOp() {
	// No op function
}

func BeginningOf28thDayOfMonth(month, year int) time.Time {
	return time.Date(year, time.Month(month), 28, 0, 0, 0, 0, time.UTC)
}

func SetDate(t time.Time, day int) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, day, 0, 0, 0, 0, t.Location())
}

func SetDateAndHour(t time.Time, day int, hour int) time.Time {
	y, m, _ := t.Date()
	return time.Date(y, m, day, hour, 0, 0, 0, t.Location())
}

// Returns no of days in a month
func NoOfDaysInMonth() int {
	return EndOfMonth().Day()
}

func BeginningLastMonth() time.Time {
	now := time.Now().UTC()
	y, m, _ := now.Date()
	return time.Date(y, m-1, 1, 0, 0, 0, 0, now.Location())
}

func BeginningNextMonth() time.Time {
	now := time.Now().UTC()
	y, m, _ := now.Date()
	return time.Date(y, m+1, 1, 0, 0, 0, 0, now.Location())
}

// TimeDiff calculates the absolute difference between 2 time instances in
// years, months, days, hours, minutes and seconds.
//
// For details, see https://stackoverflow.com/a/36531443/1705598
func TimeDiff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

// GetTimeSinceStartOfMonth returns the string of time since
// the start of the month in the same format as `time.Since`
func GetTimeSinceStartOfMonth() string {
	startOfMonth := now.BeginningOfMonth()
	currentTime := time.Now()

	return TimeDiffString(startOfMonth, currentTime)
}

// TimeDiffString returns a string similar to the Duration type. It runs
// the TimeDiff function internally. Please note that this only works for
// returning correct strings for days, hours, minutes and seconds. In case
// you're using this function for Prometheus queries that stretch over
// years or months, please do not use this function since Prometheus
// supports the following durations: milliseconds, seconds, minutes, hours, weeks and years.
// With weeks being involved in between, due to the nature of TimeDiff, this
// function will not provide the correct result.
//
// Examples:
//
//	TimeDiffString(time.Now(), time.Now().Add(-12 * time.Minute)) // 12m
//	TimeDiffString(time.Now().Add(-365 * 24 * time.Hour), time.Now()) // -8760h0m0s
func TimeDiffString(a, b time.Time) (timeString string) {
	_, _, day, hr, min, sec := TimeDiff(a, b)

	if day > 0 {
		timeString += fmt.Sprintf("%dd", day)
	}

	if hr > 0 {
		timeString += fmt.Sprintf("%dh", hr)
	}

	if min > 0 {
		timeString += fmt.Sprintf("%dm", min)
	}

	if sec > 0 {
		timeString += fmt.Sprintf("%ds", sec)
	}

	return
}

func BeginningOfInputMonth(month, year int) time.Time {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
}

func EndOfInputMonth(month, year int) time.Time {
	return BeginningOfInputMonth(month, year).AddDate(0, 1, 0).Add(-time.Nanosecond)
}

// LogFuncExecTime logs the execution time of a function and can (optionally) execute callback functions. Example of a potential callback
// is to raise an alert if the execution time is high. This is best suited to be called as a deferred function at the start of the target
// function. For usage example, refer `CheckDataUsageComponents` in the `api` service.
func LogFuncExecTime(log zerolog.Logger, funcName string, callbacks ...func(startTime time.Time)) func() {
	start := time.Now()
	return func() {
		log.Info().Str("function_name", funcName).Str("execution_time", time.Since(start).String()).Send()
		for _, callback := range callbacks {
			callback(start)
		}
	}
}

// Float64Ptr returns pointer of the float64 input
func Float64Ptr(input float64) *float64 {
	return &input
}

// TimePtr returns pointer of the input time
func TimePtr(input time.Time) *time.Time {
	return &input
}

// IsPtrNil returns true if input to function is nil or nil pointer
func IsPtrNil(v any) bool {
	return v == nil || (reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil())
}

// ToAnyMap a generic function to convert the input map to any map
func ToAnyMap[K comparable, E any](input map[K]E) map[K]any {
	result := make(map[K]any)
	for k, v := range input {
		result[k] = v
	}
	return result
}

func LogInfo(deploymentId uuid.UUID, message string) {
	log.Info().Str("deployment_id", deploymentId.String()).Msgf(message)
}

func LogError(deploymentId uuid.UUID, message string, err error) {
	log.Error().Str("deployment_id", deploymentId.String()).Err(err).Msgf(message)
}
