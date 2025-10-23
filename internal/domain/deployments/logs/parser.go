package logs

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

var (
	dockerHeaderRegex = regexp.MustCompile(`^[\x01\x02]\x00{6}.`)
	isoTimestampRegex = regexp.MustCompile(`\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d+Z\s*`)
	timestampRegex    = regexp.MustCompile(`^\[?(\d{1,2}:\d{2}:\d{2}(?:\s?[AP]M)?)\]?\s*`)
	errorPatterns     = []string{
		"[error]", "ERROR", "error:", "Error:",
		"failed", "Failed", "FAILED",
		"exception", "Exception", "EXCEPTION",
		"fatal", "Fatal", "FATAL",
		"panic", "Panic", "PANIC",
	}
	warnPatterns = []string{
		"[warn]", "[warning]", "WARN", "WARNING",
		"warn:", "warning:", "Warning:",
		"deprecated", "Deprecated", "DEPRECATED",
	}
	successPatterns = []string{
		"✓", "✔", "✅",
		"success", "Success", "SUCCESS",
		"completed", "Completed", "COMPLETED",
		"built successfully", "Built successfully",
		"modules transformed",
		"ready in",
	}
)

func (p *LogParser) ParseLine(line string) *StructuredLog {
	if line == "" {
		return nil
	}

	now := time.Now()
	relativeSeconds := int(now.Sub(p.startTime).Seconds())

	timestamp := extractTimestamp(line)
	if timestamp == "" {
		timestamp = now.Format("3:04:05 PM")
	}

	level := detectLogLevel(line)
	message := cleanMessage(line)

	return &StructuredLog{
		Timestamp:       timestamp,
		RelativeSeconds: relativeSeconds,
		Level:           level,
		Message:         message,
		Raw:             line,
	}
}

func extractTimestamp(line string) string {
	matches := timestampRegex.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func detectLogLevel(line string) LogLevel {
	lowerLine := strings.ToLower(line)

	for _, pattern := range errorPatterns {
		if strings.Contains(line, pattern) || strings.Contains(lowerLine, strings.ToLower(pattern)) {
			return LogLevelError
		}
	}

	for _, pattern := range warnPatterns {
		if strings.Contains(line, pattern) || strings.Contains(lowerLine, strings.ToLower(pattern)) {
			return LogLevelWarn
		}
	}

	for _, pattern := range successPatterns {
		if strings.Contains(line, pattern) {
			return LogLevelSuccess
		}
	}

	return LogLevelInfo
}

func cleanMessage(line string) string {
	cleaned := line

	// Remove Docker stream header (8 bytes: stream type + 6 zeros + 1 variable byte)
	// The format is: [01|02][00 00 00 00 00 00][XX] where XX is variable
	if len(cleaned) >= 8 && (cleaned[0] == '\x01' || cleaned[0] == '\x02') {
		allZeros := true
		for i := 1; i <= 6; i++ {
			if cleaned[i] != '\x00' {
				allZeros = false
				break
			}
		}
		if allZeros {
			// Check if byte 7 is start of ISO timestamp (digit) or byte 8 is
			if cleaned[7] >= '0' && cleaned[7] <= '9' {
				// Byte 7 is first digit of year (e.g., '2' in 2025)
				cleaned = cleaned[7:]
			} else if len(cleaned) > 8 && cleaned[8] >= '0' && cleaned[8] <= '9' {
				// Byte 8 is first digit, byte 7 was a separator
				cleaned = cleaned[8:]
			} else {
				// Skip all 8 bytes
				cleaned = cleaned[8:]
			}
		}
	}

	// Remove ISO timestamps
	cleaned = isoTimestampRegex.ReplaceAllString(cleaned, "")

	// Remove simple timestamps
	cleaned = timestampRegex.ReplaceAllString(cleaned, "")

	cleaned = strings.TrimSpace(cleaned)
	return cleaned
}

func (p *LogParser) ParseAndSerialize(logText string) (string, error) {
	if logText == "" {
		return "", nil
	}

	lines := strings.Split(logText, "\n")
	var structuredLogs []StructuredLog
	seen := make(map[string]bool)

	for _, line := range lines {
		if line == "" {
			continue
		}

		if parsed := p.ParseLine(line); parsed != nil {
			key := parsed.Message
			if !seen[key] {
				seen[key] = true
				structuredLogs = append(structuredLogs, *parsed)
			}
		}
	}

	if len(structuredLogs) == 0 {
		return "", nil
	}

	jsonBytes, err := json.Marshal(structuredLogs)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
