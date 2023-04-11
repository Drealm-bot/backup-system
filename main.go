package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

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
			err = BackupFile(path, backupFolder)
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
		err := RestoreBackup(backupFolder, file)
		if err != nil {
			fmt.Println("Error al restaurar el archivo:", err)
			return
		}

		fmt.Println("El archivo fue restaurado con éxito.")
	}

}
