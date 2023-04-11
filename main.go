package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"time"
)

// Definimos una constante para el tamaño máximo de cada fragmento en bytes
const fragmentSize = 512 * 1024 * 1024

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Debe especificar la ruta de la carpeta a respaldar.")
		return
	}
	if len(os.Args) == 2 {
		// Obtener las rutas de la carpeta a respaldar y la carpeta de backup
		sourceFolder := os.Args[1]
		backupTime := time.Now().Format("2006-01-02_15:04:05")
		backupFolder := sourceFolder + "-backup_" + backupTime

		// Crear la carpeta de backup si no existe
		err := os.MkdirAll(backupFolder, 0755)
		if err != nil {
			fmt.Println("Error al crear la carpeta de backup:", err)
			return
		}

		// Recorrer los archivos de la carpeta a respaldar y hacer backup de cada uno
		err = filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error al acceder a un archivo:", err)
				return err
			}
			if info.IsDir() {
				return nil // Ignorar las carpetas
			}
			err = backupFile(path, backupFolder)
			if err != nil {
				fmt.Println("Error al hacer backup del archivo", path, ":", err)
				return err
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error al recorrer la carpeta de respaldo:", err)
			return
		}

		fmt.Println("Todos los archivos fueron respaldados con éxito.")
	}
	if len(os.Args) == 3 {
		backupFolder := os.Args[1]
		file := os.Args[2]
		err := restoreBackup(backupFolder, file)
		if err != nil {
			fmt.Println("Error al restaurar el archivo:", err)
			return
		}

		fmt.Println("El archivo fue restaurado con éxito.")
	}

}

// Función para dividir el archivo en fragmentos y guardarlos en la carpeta de backup
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
