package modifier

import (
	"fmt"
	"io/ioutil"
	"strings"
	"regexp"
	"time"
)

type Modifier struct {
	fileName string
	delay    int
}

const SrtPeriodSeparator = " --> "
const SrtModifiedTimeFormat = "15:04:05.000"

func formatPeriodLine(start string, end string) string {
	return start + SrtPeriodSeparator + end
}

func NewModifier(fileName string, addOrSub string, nbSeconds int) *Modifier {
	delay := nbSeconds
	if addOrSub == "-" {
		delay = -delay
	}
	return &Modifier{fileName, delay}
}

func (m Modifier) addDelayToMarkerString(marker string) (error, string) {
	//fmt.Printf("add delay to marker %v\n", marker)
	marker = strings.Replace(marker, ",", ".", 1)
	if parsed, err := time.Parse(SrtModifiedTimeFormat, marker); err != nil {
		return err, ""
	} else {
		parsed := parsed.Add(time.Second * time.Duration(m.delay))
		markerUpdated := parsed.Format(SrtModifiedTimeFormat)
		return nil, strings.Replace(markerUpdated, ".", ",", 1)
	}
}

func (m Modifier) processLine(line string) string {
	//fmt.Printf("process line %v\n", line)
	timePattern := "([0-9]{2}:[0-9]{2}:[0-9]{2},[0-9]{3})"
	pattern := timePattern + " --> " + timePattern
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		errStart, start := m.addDelayToMarkerString(matches[1])
		errEnd, end := m.addDelayToMarkerString(matches[2])
		if errStart != nil || errEnd != nil {
			fmt.Printf("unable to modify line %v\n", line)
			return line
		} else {
			return formatPeriodLine(start, end)
		}
	}
	return line
}

func (m Modifier) Process() error {
	//fmt.Printf("processing %v\r\n", m.fileName)
	content, err := ioutil.ReadFile(m.fileName)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for idx, line := range lines {
		lines[idx] = m.processLine(line)
	}
	//fmt.Printf("%v", lines)
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(m.fileName, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}
