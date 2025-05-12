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
	filterLevel      = flag.String("level", "", "Filter logs by level: info, warning, error, debug")

	// This regex matches the log header.
	logHeaderRegex = regexp.MustCompile(`^([IWEF])(\d{4})\s+([\d:.]+)\s+\d+\s+([^\]]+)]\s+"([^"]+)"(.*)$`)
	// Updated keyValueRegex to capture key=<...>, key="value" and key={...} patterns.
	// Group 1: key; Group 2: angle-bracket value; Group 3: quoted value; Group 4: JSON value.
	keyValueRegex = regexp.MustCompile(`(\w+)=(?:<([^>]+)>|"((?:[^"\\]|\\.)*)"|(\{.*?\}))`)

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
	level, date, timeStr, location, message, kvPart :=
		matches[1], matches[2], matches[3], matches[4], matches[5], matches[6]

	// Apply level filter
	if *filterLevel != "" {
		want := strings.ToLower(*filterLevel)
		var wantChar string
		switch want {
		case "info":
			wantChar = "I"
		case "warning", "warn":
			wantChar = "W"
		case "error":
			wantChar = "E"
		case "debug":
			wantChar = "D"
		default:
			// unrecognized, skip filtering
			wantChar = ""
		}
		if wantChar != "" && level != wantChar {
			return
		}
	}

	fields := make(map[string]string)

	// Capture key/value pairs.
	for _, kv := range keyValueRegex.FindAllStringSubmatch(kvPart, -1) {
		switch {
		case kv[2] != "":
			fields[kv[1]] = strings.TrimSpace(kv[2])
		case kv[3] != "":
			fields[kv[1]] = kv[3]
		default:
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

	// For error logs, merge the "err" field into the message and color it red.
	if level == "E" {
		if errMsg, exists := fields["err"]; exists {
			message = fmt.Sprintf("%s: %s%s%s", message, Red, errMsg, Reset)
		}
	}

	// Build headline fields.
	headerColor := colorForLevel(level)
	controller := fields["controller"]
	nsRaw := fmt.Sprintf("%s/%s", fields["namespace"], fields["name"])

	// Set namespace/name and message colors.
	nsColor := Green
	msgColor := BrightWhite
	if level == "E" {
		nsColor = Red
		msgColor = Red
	}

	nsName := nsColor + nsRaw + Reset
	coloredMsg := msgColor + message + Reset

	// Build the headline (no pipes here).
	headline := fmt.Sprintf("%s %s %s %s@%s %s %s",
		headerColor+level,
		date,
		timeStr,
		controller,
		location+Reset,
		nsName,
		coloredMsg,
	)

	// Build the structured key/value section.
	redPipe := Red + " | " + Reset
	kvParts := []string{}

	// Use an ordered key list first.
	orderedKeys := []string{"cluster", "AWSCluster", "machinePool", "AWSMachinePool"}
	printed := map[string]bool{}
	for _, k := range orderedKeys {
		if val, ok := fields[k]; ok {
			kvParts = append(kvParts,
				fmt.Sprintf("%s%s:%s %s%s", Gray, k, Reset, Gray, val+Reset))
			printed[k] = true
		}
	}
	// Then add any remaining key/value pairs (skipping those omitted).
	for k, v := range fields {
		if !omitFromKV[k] && !printed[k] {
			kvParts = append(kvParts,
				fmt.Sprintf("%s%s:%s %s%s", Gray, k, Reset, Gray, v+Reset))
		}
	}

	// Combine and print.
	if len(kvParts) > 0 {
		fmt.Println(headline + redPipe + strings.Join(kvParts, redPipe))
	} else {
		fmt.Println(headline)
	}
}

func main() {
	flag.Parse() // Parse filter flags.

	scanner := bufio.NewScanner(os.Stdin)

	var buf string
	var collecting bool

	for scanner.Scan() {
		line := scanner.Text()

		// detect start of a multi-line err block
		if !collecting && strings.Contains(line, `err=<`) && !strings.Contains(line, ">") {
			collecting = true
			buf = line
			continue
		}

		if collecting {
			// append the continued lines
			buf += " " + strings.TrimSpace(line)
			// detect end of err block
			if strings.Contains(line, ">") {
				collecting = false
				parseLine(buf)
			}
			continue
		}

		// normal single-line log
		parseLine(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}
}
