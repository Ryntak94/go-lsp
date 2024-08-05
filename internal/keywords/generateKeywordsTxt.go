package keywords

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func generateKeywordsTxt(logger *log.Logger) *os.File {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	keywordFilePath := fmt.Sprintf("%s/internal/keywords/generatedWords.txt", cwd)
	keywordFile, err := os.Open(keywordFilePath)
	if err == nil {
		return keywordFile
	}
	keywordFile, err = os.Create(keywordFilePath)
	if err != nil {
		panic(err)
	}
	apiPath := "/usr/local/go/api"
	files, err := os.ReadDir(apiPath)
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.Name() == "README" {
			continue
		}
		openedFile, err := os.Open(fmt.Sprintf("%s/%s", apiPath, file.Name()))
		if err != nil {
			print(err)
			logger.Fatal(err)
			openedFile.Close()
			continue
		}
		scanner := bufio.NewScanner(openedFile)
		println(file.Name())
		for scanner.Scan() {
			line := scanner.Text()
			splitLine := strings.Split(line, " ")
			if len(splitLine) < 2 || splitLine[0] == "#" {
				continue
			}

			if _, err := keywordFile.WriteString(fmt.Sprintf("%s\n", line)); err != nil {
				panic(err)
			}
		}
		if err := scanner.Err(); err != nil {
			logger.Fatal(err)
		}
	}
	keywordFile.Close()

	keywordFile, err = os.Open(keywordFilePath)
	if err != nil {
		panic(err)
	}

	return keywordFile
}
