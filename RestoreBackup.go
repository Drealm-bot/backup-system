package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func RestoreBackup(backupFolder, fileName, restoredFolder string) error {
	// Leer el archivo JSON de información de la copia de seguridad
	backupInfoFile, err := os.Open(filepath.Join(backupFolder, "backup_info.json"))
	if err != nil {
		return fmt.Errorf("No se ha encontrado el archivo JSON en la carpeta de respaldo.")
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
		if fileName == "All" {
			backupFileName := backup.(map[string]interface{})["file"].(string)
			err := RestoreBackup(backupFolder, backupFileName, restoredFolder)
			if err != nil {
				return err
			}
		} else {
			if backup.(map[string]interface{})["file"].(string) == fileName {
				backupData = backup.(map[string]interface{})
				break
			}
		}
	}
	if backupData == nil {
		if fileName == "All" {
			return fmt.Errorf("Se han restaurado todos los archivos")
		}
		return fmt.Errorf("no se encontró una copia de seguridad para el archivo '%s'", fileName)
	}

	// Crear el archivo restaurado
	restoredFilePath := filepath.Join(restoredFolder, fileName)
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
