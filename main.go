package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[0;31m" // using standard red
	Yellow      = "\033[0;33m"
	Green       = "\033[0;32m"
	Blue        = "\033[0;34m"
	White       = "\033[0;37m"
	Gray        = "\033[0;90m"
	BrightWhite = "\033[1;37m"
)

var (
	// Define filtering flags.
	filterName       = flag.String("name", "", "Filter logs by the 'name' field")
	filterNamespace  = flag.String("namespace", "", "Filter logs by the 'namespace' field")
	filterController = flag.String("controller", "", "Filter logs by the 'controller' field")

	// This regex matches the log header.
	logHeaderRegex = regexp.MustCompile(`^([IWEF])(\d{4})\s+([\d:.]+)\s+\d+\s+([^\]]+)]\s+"([^"]+)"(.*)$`)
	// Updated keyValueRegex to capture both key="value" and key={...} patterns.
	// Group 1: key; Group 3: quoted value; Group 4: JSON value.
	keyValueRegex = regexp.MustCompile(`(\w+)=("((?:[^"\\]|\\.)*)"|(\{.*?\}))`)

	// Keys that are removed from the key/value section because they are shown in the headline.
	omitFromKV = map[string]bool{
		"controller":      true,
		"controllerGroup": true,
		"controllerKind":  true,
		"namespace":       true,
		"name":            true,
		"err":             true,
	}
)

// colorForLevel returns an ANSI color code based on the log level.
// For errors, the header components (level, date, time) will be red,
// for warnings yellow, and blue otherwise.
func colorForLevel(level string) string {
	switch level {
	case "E":
		return Red
	case "W":
		return Yellow
	default:
		return Blue
	}
}

func parseLine(line string) {
	matches := logHeaderRegex.FindStringSubmatch(line)
	if matches == nil {
		fmt.Println(line)
		return
	}

	// Parse header parts.
	level, date, timeStr, location, message, kvPart := matches[1], matches[2], matches[3], matches[4], matches[5], matches[6]
	fields := make(map[string]string)

	// Capture key/value pairs.
	kvs := keyValueRegex.FindAllStringSubmatch(kvPart, -1)
	for _, kv := range kvs {
		// If group 3 (quoted value) is non-empty, use it; otherwise use group 4 (the JSON value).
		if kv[3] != "" {
			fields[kv[1]] = kv[3]
		} else {
			fields[kv[1]] = kv[4]
		}
	}

	// Apply filtering.
	// If a filter flag is specified and the corresponding field does not match, skip the line.
	if *filterName != "" && fields["name"] != *filterName {
		return
	}
	if *filterNamespace != "" && fields["namespace"] != *filterNamespace {
		return
	}
	if *filterController != "" && fields["controller"] != *filterController {
		return
	}

	// Build headline fields.
	headerColor := colorForLevel(level)
	controller := fields["controller"]
	nsName := fmt.Sprintf("%s/%s", fields["namespace"], fields["name"])

	// For error logs, capture the "err" field and prepare it for printing in red.
	var errField string
	if level == "E" {
		if errMsg, exists := fields["err"]; exists {
			errField = Red + errMsg + Reset
		}
	}

	// Build the headline (no pipes here).
	// The log level, date, time, and location are colored using headerColor.
	// Controller and nsName are colored in white/green, while the message is bright white.
	headline := fmt.Sprintf("%s %s %s %s@%s %s %s",
		headerColor+level+Reset,
		headerColor+date+Reset,
		headerColor+timeStr+Reset,
		White+controller+Reset,
		White+location+Reset,
		Green+nsName+Reset,
		BrightWhite+message+Reset,
	)

	// If there is an error field (only for error logs), append it after the message.
	if errField != "" {
		headline = headline + " - " + errField
	}

	// Build the structured key/value section.
	// Only the fields after the message (like cluster, AWSCluster, etc.) will be joined with a red pipe.
	redPipe := Red + " | " + Reset
	kvParts := []string{}

	// Use an ordered key list first.
	orderedKeys := []string{"cluster", "AWSCluster", "machinePool", "AWSMachinePool"}
	printed := map[string]bool{}
	for _, k := range orderedKeys {
		if val, ok := fields[k]; ok {
			kvParts = append(kvParts, fmt.Sprintf("%s%s:%s %s%s", Gray, k, Reset, Gray, val+Reset))
			printed[k] = true
		}
	}
	// Then add any remaining key/value pairs (skipping those omitted).
	for k, v := range fields {
		if !omitFromKV[k] && !printed[k] {
			kvParts = append(kvParts, fmt.Sprintf("%s%s:%s %s%s", Gray, k, Reset, Gray, v+Reset))
		}
	}

	// Join the key/value parts with the red pipe separator, but only add the section if there are any fields.
	var kvSection string
	if len(kvParts) > 0 {
		kvSection = redPipe + strings.Join(kvParts, redPipe)
	}

	// Combine the headline with the structured key/value section.
	finalLine := headline + kvSection
	fmt.Println(finalLine)
}

func main() {
	flag.Parse() // Parse filter flags.

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		parseLine(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
