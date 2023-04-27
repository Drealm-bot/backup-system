package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"math"
	"os"
	"path/filepath"
)

// Constante para el tamaño máximo de cada fragmento en bytes
const fragmentSize = 2 * 1024 * 1024

func Backup(path, backupFolder string) error {
	// Verificar info del archivo
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		// Buscar archivos
		err = filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				return BackupFile(filePath, backupFolder)
			}
			return nil
		})
		if err != nil {
			return err
		}
	} else {
		return BackupFile(path, backupFolder)
	}

	return nil
}

func BackupFile(filePath, backupFolder string) error {

	// Abrir archivo
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener información
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Calcular fragmentos necesarios
	numFragments := int(math.Ceil(float64(fileInfo.Size()) / float64(fragmentSize)))

	// Crear un slice para almacenar los fragmentos
	fragments := make([]string, numFragments)

	// Leer y guardar fragmentos
	for i := 0; i < numFragments; i++ {

		fragmentName := fmt.Sprintf("%s-arrempujala-%d", fileInfo.Name(), i)
		fragmentPath := filepath.Join(backupFolder, fragmentName)
		fragmentFile, err := os.Create(fragmentPath)
		if err != nil {
			return err
		}
		defer fragmentFile.Close()

		// Leer datos del fragmento
		buffer := make([]byte, fragmentSize)
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		// Guardar fragmento
		_, err = fragmentFile.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}

		// Agregar el fragmento al slice de fragmentos
		fragments[i] = fragmentName
	}

	// Leer el archivo JSON existente
	existingBackupInfo := make(map[string]interface{})
	backupInfoFile, err := os.Open(filepath.Join(backupFolder, "backup_info.json"))
	if err == nil {
		defer backupInfoFile.Close()
		decoder := json.NewDecoder(backupInfoFile)
		err = decoder.Decode(&existingBackupInfo)
		if err != nil {
			return err
		}
	}

	// Agregar información copia de seguridad
	newBackupInfo := map[string]interface{}{
		"file":        fileInfo.Name(),
		"size":        fileInfo.Size(),
		"fragments":   fragments,
		"total_parts": numFragments,
	}

	if _, ok := existingBackupInfo["backups"]; !ok {
		existingBackupInfo["backups"] = make([]interface{}, 0)
	}
	existingBackupInfo["backups"] = append(existingBackupInfo["backups"].([]interface{}), newBackupInfo)

	// Crear un archivo JSON con la información actualizada de la copia de seguridad
	jsonData, err := json.MarshalIndent(existingBackupInfo, "", "  ")
	if err != nil {
		return err
	}

	backupInfoFile, err = os.Create(filepath.Join(backupFolder, "backup_info.json"))
	if err != nil {
		return err
	}
	defer backupInfoFile.Close()

	_, err = backupInfoFile.Write(jsonData)
	if err != nil {
		return err
	}

	fmt.Println("Copia de seguridad creada con éxito.")

	return nil
}
