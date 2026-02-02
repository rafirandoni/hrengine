package utility

import (
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"strings"
)

type Env interface {
	GetString(key string, defaultValue ...string) string
	GetBool(key string) bool
	GetSliceString(key, delimiter string, defaultValues ...string) []string
	GetInt(key string, defaultValue ...int) (int, error)
	GetInt64(key string, defaultValue ...int64) (int64, error)
	GetFloat64(key string, defaultValue ...float64) (float64, error)
	GetBase64(key string) ([]byte, error)
}

type env struct{}

func NewEnv() Env {
	return &env{}
}

func (slf *env) GetString(key string, defaultValue ...string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val != "" {
		return val
	}

	if len(defaultValue) > 0 {
		return defaultValue[0]
	}

	return ""
}

func (slf *env) GetBool(key string) bool {
	val := strings.TrimSpace(strings.ToUpper(os.Getenv(key)))
	if val != "" && val == "TRUE" {
		return true
	}

	return false
}

func (slf *env) GetSliceString(key, delimiter string, defaultValues ...string) []string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return defaultValues
	}

	vals := strings.Split(key, delimiter)
	return vals
}

func (slf *env) GetInt(key string, defaultValue ...int) (int, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val != "" {
		intVal, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}

		return intVal, nil
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return 0, errors.New("value is not int")
}

func (slf *env) GetInt64(key string, defaultValue ...int64) (int64, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val != "" {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return 0, err
		}

		return intVal, nil
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return 0, errors.New("value is not int64")
}

func (slf *env) GetFloat64(key string, defaultValue ...float64) (float64, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val != "" {
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return 0, err
		}

		return floatVal, nil
	}

	if len(defaultValue) > 0 {
		return defaultValue[0], nil
	}

	return 0, errors.New("value is not float64")
}

func (slf *env) GetBase64(key string) ([]byte, error) {
	base64Val := slf.GetString(key)
	return base64.StdEncoding.DecodeString(base64Val)
}
