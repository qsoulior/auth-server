// Package config provides functions to read env variables
// and set structure fields to variables.
package config

import (
	"bufio"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// EnvGetter is a type of function returning a value by key.
type EnvGetter func(key string) (string, bool)

// setValue parses val and sets fieldVal to val.
// It returns error if parsing failed.
func setValue(fieldVal reflect.Value, val string) error {
	fieldType := fieldVal.Type()
	switch fieldVal.Kind() {
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		fieldVal.SetBool(boolVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(val, 0, fieldType.Bits())
		if err != nil {
			return err
		}
		fieldVal.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(val, 0, fieldType.Bits())
		if err != nil {
			return err
		}
		fieldVal.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(val, fieldType.Bits())
		if err != nil {
			return err
		}
		fieldVal.SetFloat(floatVal)
	case reflect.String:
		fieldVal.SetString(val)
	case reflect.Slice:
		vals := strings.Split(strings.TrimSpace(val), ",")
		sliceVal := reflect.MakeSlice(fieldType, len(vals), len(vals))
		for i, val := range vals {
			if err := setValue(sliceVal.Index(i), val); err != nil {
				return nil
			}
		}
		fieldVal.Set(sliceVal)
	}
	return nil
}

// read gets env variables and sets rv recursively to variables.
// It returns error if setting failed.
func read(rv reflect.Value, getter EnvGetter) error {
	rvType := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fieldVal := rv.Field(i)
		if fieldVal.Kind() == reflect.Struct {
			if err := read(fieldVal, getter); err != nil {
				return err
			}
		} else {
			field := rvType.Field(i)

			key := field.Tag.Get("env")
			val, exists := getter(key)
			if !exists {
				if val, exists = field.Tag.Lookup("default"); !exists {
					return EmptyError(key)
				}
			}

			if err := setValue(fieldVal, val); err != nil {
				return ParseError(key)
			}
		}
	}
	return nil
}

// valueOf checks kind of cfg and initializes reflect.Value.
// It returns error if cfg is not a structure or a pointer to structure.
func valueOf(cfg any) (reflect.Value, error) {
	rv := reflect.ValueOf(cfg)
	if rv.Kind() != reflect.Pointer {
		return reflect.Value{}, ErrNotPointer
	}

	re := rv.Elem()
	if re.Kind() != reflect.Struct {
		return reflect.Value{}, ErrNotStruct
	}

	return re, nil
}

// ReadEnv gets env variables from environment and sets cfg to variables.
// It returns error if cfg has incorrect type or read failed.
func ReadEnv(cfg any) error {
	re, err := valueOf(cfg)
	if err != nil {
		return err
	}
	return read(re, os.LookupEnv)
}

// ReadEnvFile gets env variables from file and sets cfg to variables.
// It returns error if cfg has incorrect type or read failed.
func ReadEnvFile(cfg any, path string) error {
	re, err := valueOf(cfg)
	if err != nil {
		return err
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	env := make(map[string]string, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.SplitN(scanner.Text(), "=", 2)
		key, val := line[0], line[1]
		env[key] = val
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return read(re, func(key string) (string, bool) {
		val, exists := env[key]
		return val, exists
	})
}
