package utils

import (
	"crypto/rand"
	"fmt"
	"strconv"
	"strings"
)

func  GenerateUniqueAudioFileName() string {
    b := make([]byte, 8)
    _, _ = rand.Read(b)
    return fmt.Sprintf("audio_%x.mp3", b)
}

func GetName(names string) string {
	parts := strings.Split(names, " ")
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}

func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func ParseStringToBool(value *string) (*bool, error) {
	if value == nil || *value == "" {
		return nil, nil
	}

	normalizedValue := strings.ToLower(strings.TrimSpace(*value))
	
	switch normalizedValue {
	case "true":
		result := true
		return &result, nil
	case "false":
		result := false
		return &result, nil
	default:
		result := false
		return &result, nil
	}
}

func ParseStringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("error converting string to int: %w", err)
	}
	return i, nil
}