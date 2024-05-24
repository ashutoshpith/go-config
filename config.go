package goconfig

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func EnableConfigFile(cfg interface{}, filepath ...string) {
	var fileExist string

	if len(filepath) > 0 {
		fileExist = filepath[0]
	} else {
		fileExist = ".sias"
	}

	file, err := os.Open(fileExist)
	if err != nil {
		log.Println(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	configMap := make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		configMap[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
	}
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		configName := fieldType.Tag.Get("config")
		if configName == "" {
			continue
		}

		configValue, exists := configMap[configName]
		if !exists {
			continue
		}
		switch field.Kind() {
		case reflect.Int:
			intValue, err := strconv.Atoi(configValue)
			if err != nil {
				fmt.Println("Failed to parse ", configName, err)
			}
			field.SetInt(int64(intValue))
		case reflect.String:
			field.SetString(configValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(configValue)
			if err != nil {
				fmt.Println("Failed to parse bool ", configName, err)
			}
			field.SetBool(boolValue)
		default:
			fmt.Println("Unsupported format ")
		}

	}

}
