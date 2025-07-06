package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ParseJSONResponsePermissive(response string, out any) error {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start == -1 || end == -1 || start > end {
		return fmt.Errorf("no JSON object found in response: %s", response)
	}
	jsonPart := response[start : end+1]

	if err := json.Unmarshal([]byte(jsonPart), out); err != nil {
		return fmt.Errorf("failed to unmarshal accusation response: %w", err)
	}

	return nil
}
