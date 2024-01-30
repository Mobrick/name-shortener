package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Mobrick/name-shortener/internal/models"
)

func LoadURLRecords(fileName string) ([]models.URLRecord, error) {

	err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var records []models.URLRecord

	for scanner.Scan() {
		line := scanner.Bytes()
		var record models.URLRecord
		err := json.Unmarshal(line, &record)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func UploadNewURLRecord(newRecord models.URLRecord, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return err
    }
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(newRecord)
	if err != nil {
		return err
	}
	return nil
}