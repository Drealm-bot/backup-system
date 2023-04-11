package src

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
)

// Definimos una constante para el tamaño máximo de cada fragmento en bytes
const fragmentSize = 512 * 1024 * 1024

func backupFile(filePath, backupFolder string) error {

	// Abrir el archivo original para leer los datos
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener información del archivo original
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	// Calcular el número total de fragmentos necesarios
	numFragments := int(math.Ceil(float64(fileInfo.Size()) / float64(fragmentSize)))

	// Crear un slice para almacenar los fragmentos
	fragments := make([]string, numFragments)

	// Leer y guardar cada fragmento
	for i := 0; i < numFragments; i++ {
		// Definir el nombre del fragmento
		fragmentName := fmt.Sprintf("%s-%d", fileInfo.Name(), i)

		// Crear el archivo del fragmento
		fragmentPath := filepath.Join(backupFolder, fragmentName)
		fragmentFile, err := os.Create(fragmentPath)
		if err != nil {
			return err
		}
		defer fragmentFile.Close()

		// Leer los datos del fragmento
		buffer := make([]byte, fragmentSize)
		bytesRead, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}

		// Guardar el fragmento
		_, err = fragmentFile.Write(buffer[:bytesRead])
		if err != nil {
			return err
		}

		// Agregar el fragmento al slice de fragmentos
		fragments[i] = fragmentName
	}

	// Leer el archivo JSON existente (si existe)
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

	// Agregar la información de la nueva copia de seguridad
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
