package env

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

type EnvVar struct {
	Key   string
	Value string
}

var envVars = make([]EnvVar, 0)

func ParseEnv(filename string) error {
	// Exit if no file
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		envVars = append(envVars, EnvVar{
			Key:   strings.TrimSpace(parts[0]),
			Value: strings.TrimSpace(parts[1]),
		})
	}

	return scanner.Err()
}

func ReplaceEnv(input string) string {
	re := regexp.MustCompile(`env\(([A-Z\d_]+)\)`)

	result := re.ReplaceAllStringFunc(input, func(match string) string {
		sub := re.FindStringSubmatch(match)
		if len(sub) != 2 {
			return match
		}

		for _, envVar := range envVars {
			if envVar.Key == sub[1] {
				return envVar.Value
			}
		}

		return ""
	})

	return result
}
