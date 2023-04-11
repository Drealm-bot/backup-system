package src

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func restoreBackup(backupFolder, fileName string) error {
	// Leer el archivo JSON de información de la copia de seguridad
	backupInfoFile, err := os.Open(filepath.Join(backupFolder, "backup_info.json"))
	if err != nil {
		return err
	}
	defer backupInfoFile.Close()

	backupInfo := make(map[string]interface{})
	decoder := json.NewDecoder(backupInfoFile)
	err = decoder.Decode(&backupInfo)
	if err != nil {
		return err
	}

	// Buscar la información de la copia de seguridad del archivo especificado
	var backupData map[string]interface{}
	for _, backup := range backupInfo["backups"].([]interface{}) {
		if backup.(map[string]interface{})["file"].(string) == fileName {
			backupData = backup.(map[string]interface{})
			break
		}
	}
	if backupData == nil {
		return fmt.Errorf("no se encontró una copia de seguridad para el archivo '%s'", fileName)
	}

	// Crear el archivo restaurado
	restoredFilePath := filepath.Join(backupFolder, fileName)
	restoredFile, err := os.Create(restoredFilePath)
	if err != nil {
		return err
	}
	defer restoredFile.Close()

	// Escribir los fragmentos en el archivo restaurado
	for _, fragmentName := range backupData["fragments"].([]interface{}) {
		fragmentPath := filepath.Join(backupFolder, fragmentName.(string))
		fragmentFile, err := os.Open(fragmentPath)
		if err != nil {
			return err
		}
		defer fragmentFile.Close()

		_, err = io.Copy(restoredFile, fragmentFile)
		if err != nil {
			return err
		}
	}

	fmt.Println("Archivo restaurado con éxito.")

	return nil
}
