package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Mobrick/name-shortener/internal/models"
)

func LoadURLRecords(file *os.File) ([]models.URLRecord, error) {
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

func UploadNewURLRecord(newRecord models.URLRecord, file *os.File) error {
	// На случай если флаг адреса файла был пуст, пропускаем добавление в файл
	if file == nil {
		return nil
	}
	encoder := json.NewEncoder(file)
	err := encoder.Encode(newRecord)
	if err != nil {
		return err
	}
	return nil
}

func MakeFile(fileName string) (*os.File, error) {	
	// На случай если флаг адреса файла был пуст, открытие файла
	if len(fileName) == 0 {
		return nil, nil
	}

	err := os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
