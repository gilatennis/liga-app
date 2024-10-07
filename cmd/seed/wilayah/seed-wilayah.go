package wilayah

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
)

func SeedWilayah(db *gorm.DB) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	sqlPath := fmt.Sprintf("%s/cmd/seed/wilayah/wilayah.sql", pwd)

	file, err := os.Open(sqlPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	const maxCapacity int = 100_000
	buf := make([]byte, maxCapacity)

	scanner.Buffer(buf, maxCapacity)

	insertStatements := []string{}

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, ";") {
			lines = append(lines, line)
			insertStatements = append(insertStatements, strings.Join(lines, " "))
			lines = []string{}
			continue
		}
		lines = append(lines, line)
	}

	for _, s := range insertStatements {
		fmt.Println("s: ", s)
	}

	chunkSize := 100

	for i := 0; i < len(insertStatements); i += chunkSize {
		// TODO: chunk insert statements before executing it
		// https://stackoverflow.com/questions/35179656/slice-chunking-in-go
	}

	/*
		err = db.Exec(string(c)).Error
		if err != nil {
			return err
		}
	*/

	return nil
}
