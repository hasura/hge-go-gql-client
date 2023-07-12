package util

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/now"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestUniqueStrings(t *testing.T) {
	fixtures := []struct {
		Input    []string
		Expected string
	}{
		{Input: []string{}, Expected: ""},
		{Input: []string{"a", "b", "c"}, Expected: "a,b,c"},
		{Input: []string{"a", "b", "b", "c"}, Expected: "a,b,c"},
		{Input: []string{"a", "b", "c", "a"}, Expected: "a,b,c"},
		{Input: []string{"c", "b", "c", "a", "b", "c"}, Expected: "a,b,c"},
	}

	for i, ss := range fixtures {
		sample := StringList{}
		for _, s := range ss.Input {
			sample.Add(s)
		}
		assert.Equal(t, ss.Expected, sample.UniqueString(), i)
	}
}

func TestTimeDiff(t *testing.T) {
	type timeStruct struct {
		year   int
		month  int
		day    int
		hour   int
		minute int
		second int
	}
	fixtures := []struct {
		StartTime time.Time
		EndTime   time.Time
		Expected  timeStruct
	}{
		{
			StartTime: time.Date(2015, 5, 1, 0, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2016, 6, 2, 1, 1, 1, 1, time.UTC),
			Expected:  timeStruct{1, 1, 1, 1, 1, 1},
		},
		{
			StartTime: time.Date(2016, 1, 2, 0, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2016, 2, 1, 0, 0, 0, 0, time.UTC),
			Expected:  timeStruct{0, 0, 30, 0, 0, 0},
		},
		{

			StartTime: time.Date(2015, 2, 11, 0, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2016, 1, 12, 0, 0, 0, 0, time.UTC),
			Expected:  timeStruct{0, 11, 1, 0, 0, 0},
		},
		{
			StartTime: time.Date(2023, 2, 14, 3, 0, 0, 0, time.UTC),
			EndTime:   time.Date(2023, 2, 1, 1, 0, 0, 0, time.UTC),
			Expected:  timeStruct{0, 0, 13, 2, 0, 0},
		},
	}

	for _, ss := range fixtures {
		y, m, d, h, mm, s := TimeDiff(ss.StartTime, ss.EndTime)
		assert.Equal(t, ss.Expected.year, y)
		assert.Equal(t, ss.Expected.month, m)
		assert.Equal(t, ss.Expected.day, d)
		assert.Equal(t, ss.Expected.hour, h)
		assert.Equal(t, ss.Expected.minute, mm)
		assert.Equal(t, ss.Expected.second, s)
	}
}
func TestTimePtr(t *testing.T) {
	for i, fixture := range []time.Time{time.Now(), time.Now().Add(time.Minute)} {
		assert.Equal(t, fixture, *TimePtr(fixture), "%d", i)
	}
}

func TestIsPtrNil(t *testing.T) {
	type Foo interface {
		foo() string
	}

	var foo Foo
	var numPtr *int
	for val, isNil := range map[any]bool{
		nil:    true,
		foo:    true,
		0:      false,
		"":     false,
		numPtr: true,
	} {
		assert.Equal(t, isNil, IsPtrNil(val), "%+v", val)
	}
}

func TestFloat64Ptr(t *testing.T) {
	for i, fixture := range []float64{0, 0.1, 0.001} {
		assert.Equal(t, fixture, *Float64Ptr(fixture), "%d", i)
	}
}

func TestToAnyMap(t *testing.T) {
	strMap := map[string]string{
		"foo": "bar",
	}
	for k, v := range ToAnyMap(strMap) {
		strMap[k] = v.(string)
	}
}

func TestEndOfInputMonth(t *testing.T) {
	for i, fixture := range []struct {
		month int
		year  int
		day   int
	}{
		{1, 2000, 31},
		{2, 2004, 29},
	} {
		assert.Equal(t, fixture.day, EndOfInputMonth(fixture.month, fixture.year).Day(), "%d", i)
	}

	assert.Equal(t, EndOfMonth().Day(), NoOfDaysInMonth())
}

func TestGetTimeSinceStartOfMonth(t *testing.T) {
	assert.Equal(t, TimeDiffString(now.BeginningOfMonth(), time.Now().Local()), GetTimeSinceStartOfMonth())
}

func TestLogFuncExecTime(t *testing.T) {
	isCalled := false
	logger := LogFuncExecTime(log.Logger, "test", func(startTime time.Time) {
		isCalled = true
	})
	logger()
	assert.True(t, isCalled)
}

func TestJsonList(t *testing.T) {
	jsonList := JsonList{}
	jsonList.Add("{}")
	assert.False(t, jsonList.IsEmpty())
	assert.Equal(t, `{}`, string(jsonList.ToJsonSingleElementOrArray()))

	jsonList.Add(`{"foo":"bar"}`)
	assert.Equal(t, `[{},{"foo":"bar"}]`, string(jsonList.ToJson()))
	assert.Equal(t, `[{},{"foo":"bar"}]`, string(jsonList.ToJsonSingleElementOrArray()))
	assert.Equal(t, JsonList(JsonList{"{}", "{\"foo\":\"bar\"}"}), jsonList.Compact())
}

func TestMapStringToUUID(t *testing.T) {
	uuids := []uuid.UUID{
		uuid.New(),
		uuid.New(),
		uuid.New(),
		uuid.New(),
	}

	var strUUIDs []string
	for _, u := range uuids {
		strUUIDs = append(strUUIDs, u.String())
	}

	result, err := MapStringToUUID(strUUIDs)
	assert.NoError(t, err)
	assert.Equal(t, uuids, result)

	_, err = MapStringToUUID([]string{"xxxx"})
	assert.EqualError(t, err, `parse uuid error invalid UUID length: 4: "xxxx"`)
}

func TestPostgresArrayToStrings(t *testing.T) {
	r1, err := PostgresArrayToStrings("")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, r1)

	r2, err := PostgresArrayToStrings("{ }")
	assert.NoError(t, err)
	assert.Equal(t, []string{}, r2)

	r3, err := PostgresArrayToStrings("{a,b,c}")
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, r3)

	_, err = PostgresArrayToStrings("{")
	assert.EqualError(t, err, "invalid Postgres array: {")
}
