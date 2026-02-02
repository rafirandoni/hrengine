package utility

import (
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
)

func ToPointer[T any](val T) *T {
	return &val
}

func ToValue[T any](val *T, defaultValue ...T) T {
	if val == nil {
		return defaultValue[0]
	}

	return *val
}

func IsZeroValue[T any](val T) bool {
	ref := reflect.ValueOf(val)

	if ref.IsZero() {
		return true
	}

	switch ref.Kind() {
	case reflect.Pointer:
		return ref.IsNil()
	case reflect.Slice, reflect.Array, reflect.Map:
		return ref.Len() < 1
	}

	return false
}

func Ternary[T any](condition bool, a T, b T) T {
	if condition {
		return a
	}
	return b
}

func TernaryFA[T any](condition bool, a func() T, b T) T {
	if condition {
		return a()
	}

	return b
}

func TernaryFB[T any](condition bool, a T, b func() T) T {
	if condition {
		return a
	}

	return b()
}

func GetRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)

	return basepath + "/../.."
}

func RandomGenerator(letters string, length int) string {
	var str strings.Builder
	str.Grow(length)

	const (
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	src := rand.NewSource(time.Now().UnixNano())
	letterLength := len(letters)

	for i, cache, remain := 0, src.Int63(), letterIdxMax; i < length; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}

		if idx := int(cache & int64(letterIdxMask)); idx < letterLength {
			str.WriteByte(letters[idx])
			i++
		}
		cache >>= letterIdxBits
		remain--
	}

	return str.String()
}

func RandomString(length int) string {
	letters := `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ`
	str := RandomGenerator(letters, length)
	return str
}

func RandomNumeric(length int) int {
	letters := `123456789`
	randomStr := RandomGenerator(letters, length)

	randomNumber, _ := strconv.Atoi(randomStr)
	return randomNumber
}

func CleanUrlScheme(url string) string {
	url = strings.ReplaceAll(url, "https://", "")
	url = strings.ReplaceAll(url, "http://", "")

	return url
}

func CleanPhoneNumber(phoneNumber string) string {
	if phoneNumber[0:1] == "+" {
		phoneNumber = phoneNumber[1:]
	}

	if phoneNumber[0:1] == "0" {
		phoneNumber = "62" + phoneNumber[1:]
	}

	phoneNumber = strings.ReplaceAll(phoneNumber, "-", "")
	phoneNumber = strings.ReplaceAll(phoneNumber, " ", "")

	return phoneNumber
}

func TimeNow(timezone ...string) time.Time {
	if len(timezone) == 0 {
		defaultLoc := time.UTC.String()

		if os.Getenv("TIMEZONE") != "" {
			defaultLoc = os.Getenv("TIMEZONE")
		}

		timezone = []string{
			defaultLoc,
		}
	}

	loc, err := time.LoadLocation(timezone[0])
	if err != nil {
		log.Println("load location err, ", err.Error())
	}

	return time.Now().In(loc)
}

func TimeParse(layout string, value string, timezone ...string) time.Time {
	tm, _ := time.Parse(layout, value)
	return TimeAddLocation(tm)
}

func TimeAddLocation(tm time.Time, timezones ...string) time.Time {
	timezone := ""
	if timezone == "" {
		timezone = time.UTC.String()
	}

	if os.Getenv("TIMEZONE") != "" {
		timezone = os.Getenv("TIMEZONE")
	}

	if len(timezones) > 0 {
		timezone = timezones[0]
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Println("load location err, ", err.Error())
	}

	return time.Date(tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond(), loc)
}

func TimeIn(tm time.Time, timezones ...string) time.Time {
	timezone := ""
	if timezone == "" {
		timezone = time.UTC.String()
	}

	if os.Getenv("TIMEZONE") != "" {
		timezone = os.Getenv("TIMEZONE")
	}

	if len(timezones) > 0 {
		timezone = timezones[0]
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		log.Println("load location err, ", err.Error())
	}

	return tm.In(loc)
}

func TimeLastDateOnMonth(tm time.Time) int {
	newTm := time.Date(tm.Year(), tm.Month(), 1, 0, 0, 0, 0, tm.Location())
	nextMonth := newTm.AddDate(0, 1, 0)
	_, _, date := nextMonth.AddDate(0, 0, -1).Date()
	return date
}

func TimeStartEndOnMonth(tm time.Time) (time.Time, time.Time) {
	startTm := time.Date(tm.Year(), tm.Month(), 1, 0, 0, 0, 0, tm.Location())
	endTm := time.Date(tm.Year(), tm.Month(), TimeLastDateOnMonth(tm), 23, 59, 59, 0, tm.Location())

	return startTm, endTm
}

func InArray[T comparable](search T, arr []T) bool {
	for _, item := range arr {
		if search == item {
			return true
		}
	}

	return false
}

func GenerateUuid() string {
	uid, err := uuid.NewV7()
	if err != nil {
		uid, _ = uuid.NewV7()
	}
	return uid.String()
}

func GenerateXid() string {
	id := xid.New()
	return id.String()
}
