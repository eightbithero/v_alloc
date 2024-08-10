package main

// TODO add scan exportzp / export variables

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type MemoryRange struct {
	Label      string
	Start, End int
}

type MemoryAllocation struct {
	Ranges []MemoryRange
}

func (ma *MemoryAllocation) SortRangesByStart() {
	sort.Slice(ma.Ranges, func(i, j int) bool {
		return ma.Ranges[i].Start < ma.Ranges[j].Start
	})
}

func (ma *MemoryAllocation) AddRange(label string, start, end int) {
	ma.Ranges = append(ma.Ranges, MemoryRange{Label: label, Start: start, End: end})
}

func (ma *MemoryAllocation) CheckOverlapping() []MemoryRange {
	var overlapping []MemoryRange
	for i := 0; i < len(ma.Ranges); i++ {
		for j := i + 1; j < len(ma.Ranges); j++ {
			if ma.Ranges[i].End >= ma.Ranges[j].Start && ma.Ranges[i].Start <= ma.Ranges[j].End {
				overlapping = append(overlapping, ma.Ranges[i], ma.Ranges[j])
			}
		}
	}
	return overlapping
}
func (ma *MemoryAllocation) UpdateSize(label string, newSize int) {
	for i := range ma.Ranges {
		if ma.Ranges[i].Label == label {
			ma.Ranges[i].End = ma.Ranges[i].Start + newSize - 1
			return
		}
	}
	fmt.Printf("Range with label %s not found.\n", label)
}

func (ma *MemoryAllocation) Display() {
	fmt.Println("Memory Allocation:")
	for _, r := range ma.Ranges {
		fmt.Printf("%s - %s %s\n", formatHex(r.Start, 6), formatHex(r.End, 6), r.Label)
	}
}

func (ma *MemoryAllocation) DisplayOverlapping() {
	overlapping := ma.CheckOverlapping()
	if len(overlapping) == 0 {
		fmt.Println("No overlap detected.")
	} else {
		fmt.Println("Overlap detected:")
		for _, r := range overlapping {
			fmt.Printf("%s - %s %s\n", formatHex(r.Start, 6), formatHex(r.End, 6), r.Label)
		}
	}
}

func parseLine(line string) (string, string, error) {

	re := regexp.MustCompile(`^al (\w+) (.+)$`)
	matches := re.FindStringSubmatch(line)

	if len(matches) != 3 {
		return "", "", fmt.Errorf("Cannot parse string %s", line)
	}

	address, identifier := matches[1], matches[2][1:]
	addressInt, err := strconv.ParseInt(matches[1], 16, 64)

	if err != nil {
		return "", "", fmt.Errorf("ParseInt error: %s", line)
	}

	if 0 <= addressInt && addressInt <= 0x7FF {
		return address, identifier, nil
	}

	return "", "", fmt.Errorf("Otside the range error %s", line)
}

func formatHex(num int, length int) string {
	hexString := fmt.Sprintf("%X", num)
	for len(hexString) < length {
		hexString = "0" + hexString
	}
	return hexString
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <labels_filepath> <directory_to_analyze> \n", os.Args[0])
		return
	}

	dir := os.Args[2]
	labelsFilepath := os.Args[1]

	ma := &MemoryAllocation{}

	file, err := os.Open(labelsFilepath)
	if err != nil {
		fmt.Println("Cannot open file:", err)
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Cannot close file:", err)
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		address, identifier, err := parseLine(line)
		if err == nil {
			intAddr, _ := strconv.ParseInt(address, 16, 16)
			if strings.HasPrefix(identifier, "FAMISTUDIO_SFX") {
				ma.AddRange(identifier, int(intAddr), int(intAddr)+14) //  FAMISTUDIO_SFX_STRUCT_SIZE = 15
			} else {
				ma.AddRange(identifier, int(intAddr), int(intAddr))
			}
		}
	}

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Cannot open file:", err)
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var prevLine string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, ".res") {
				if prevLine != "" {
					re := regexp.MustCompile(`^\s*([A-Za-z0-9_]+):\s*$`) // label:
					matches := re.FindStringSubmatch(prevLine)
					if len(matches) > 1 {
						identifier := matches[1]
						re = regexp.MustCompile(`\s*\.res\s+(\d+)\s*$`) // .res <number>
						matches = re.FindStringSubmatch(line)
						if len(matches) > 1 {
							size, _ := strconv.Atoi(matches[1])
							ma.UpdateSize(identifier, size)
						}
					}
				}
			}
			prevLine = line
		}
		return nil
	})
	if err != nil {
		fmt.Println("Scanning folder error:", err)
		return
	}

	ma.SortRangesByStart()
	//ma.Display()
	ma.DisplayOverlapping()
}
