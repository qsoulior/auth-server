package config

import (
	"bufio"
	"os"
	"reflect"
	"strings"
)

type EnvGetter func(key string) (string, bool)

func read(rv reflect.Value, getter EnvGetter) error {
	rvType := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Field(i)
		if field.Kind() == reflect.Struct {
			if err := read(field, getter); err != nil {
				return err
			}
		} else {
			fieldType := rvType.Field(i)
			key := fieldType.Tag.Get("env")

			if val, valExists := getter(key); valExists {
				field.SetString(val)
				continue
			}

			if def, defExists := fieldType.Tag.Lookup("default"); defExists {
				field.SetString(def)
				continue
			}

			return NewEmptyError(key)
		}
	}
	return nil
}

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

func ReadEnv(cfg any) error {
	re, err := valueOf(cfg)
	if err != nil {
		return err
	}
	return read(re, os.LookupEnv)
}

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
