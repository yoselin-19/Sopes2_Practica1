package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func getNumericFileInfos() []os.FileInfo {
	files, err := ioutil.ReadDir("/proc")

	if err != nil {
		log.Fatal(err)
	}

	var numberProcs []os.FileInfo

	for _, f := range files {
		_, err := strconv.Atoi(f.Name())
		if err == nil {
			numberProcs = append(numberProcs, f)
		}
	}

	return numberProcs
}

func getLinuxProcesses() []map[string]interface{} {
	var linuxProcesses []map[string]interface{}

	numericFileInfos := getNumericFileInfos()

	for _, f := range numericFileInfos {

		file, err := os.Open("/proc/cpu_grupo14")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		fileContentBytes, err := ioutil.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}

		fileContent := fmt.Sprintf("%s", fileContentBytes)

		processInfo := extractLinuxProcessInfo(f.Name(), fileContent)

		linuxProcesses = append(linuxProcesses, processInfo)
	}

	return linuxProcesses
}

func extractLinuxProcessInfo(pid string, content string) map[string]interface{} {
	processInfo := make(map[string]interface{})

	processInfo["pid"] = pid

	var line = 0
	var saveChars = false
	var value = ""

	for _, c := range content {
		if saveChars && c != '\n' {
			value += string(c)
		}
		if c == ':' {
			saveChars = true
		}
		if c == '\n' {

			switch line {
			case 0:
				processInfo["nombre"] = strings.TrimSpace(value)
				break
			case 2:
				processInfo["estado"] = strings.TrimSpace(value)
				break
			case 8:
				processInfo["uid"] = strings.TrimSpace(value)
				break
			}

			line += 1
			saveChars = false
			value = ""
		}
	}

	return processInfo
}

/*
func killProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	err = process.Signal(syscall.Signal(0)) // if nil then is ok to kill

	if err != nil {
		return err
	}

	err = process.Kill()

	if err != nil {
		return err
	}

	return nil
}

func getCPUUsage() (int, int) {
	statFileContent := getStatFileContent()
	return extractCpuUsageFromStatFileContent(statFileContent)
}

func getRAMUsage() (int, int) {
	fileContent := getMemInfoFileContent()
	return extractRamUsageFromMemInfoContent(fileContent)
}

/**
  1. read the first line of   /proc/stat
  2. discard the first word of that first line   (it's always cpu)
  3. sum all of the times found on that first line to get the total time
  4. divide the fourth column ("idle") by the total time, to get the fraction of time spent being idle
  5. subtract the previous fraction from 1.0 to get the time spent being   not   idle
  6. multiple by   100   to get a percentage
*/

func extractCpuUsageFromStatFileContent(content string) (int, int) {
	numbers := extractNumbersFromLine(content, 0)
	var total = 0
	for _, n := range numbers {
		total += n
	}
	return numbers[3], total
}

func extractRamUsageFromMemInfoContent(content string) (int, int) {
	total := extractNumbersFromLine(content, 0)[0]
	available := extractNumbersFromLine(content, 2)[0]
	used := total - available
	return used, total
}

func extractNumbersFromLine(s string, line int) []int {

	var currentLine = 0
	var numbers []int
	var value = ""
	var firstLine = ""

	for _, c := range s {

		if currentLine == line {
			firstLine += string(c)

			if isNumber(int(c)) {
				value += string(c)
			} else {
				number, err := strconv.Atoi(value)
				if err == nil {
					numbers = append(numbers, number)
				}
				value = ""
			}
		}

		if c == '\n' {
			if currentLine >= line {
				break
			}
			currentLine += 1
		}
	}

	fmt.Printf("Linea: %s", firstLine)
	fmt.Printf("Extraido: %v \n", numbers)

	return numbers
}

func isNumber(c int) bool {
	return c >= '0' && c <= '9'
}

func getStatFileContent() string {
	file, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileContentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	fileContent := fmt.Sprintf("%s", fileContentBytes)
	return fileContent
}
func getMemInfoFileContent() string {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fileContentBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	fileContent := fmt.Sprintf("%s", fileContentBytes)
	return fileContent
}
