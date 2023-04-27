package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	switch len(os.Args) {
	case 1:
		fmt.Println("Debe especificar la ruta de la carpeta a respaldar.")
		return
	case 2:
		sourceFolder := os.Args[1]
		backupTime := time.Now().Format("2006-01-02_15:04:05")
		backupFolder := sourceFolder + "-backup_" + backupTime

		err := os.MkdirAll(backupFolder, 0755)
		if err != nil {
			fmt.Println("Error al crear la carpeta de backup:", err)
			return
		}

		err = Backup(sourceFolder, backupFolder)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Todos los archivos fueron respaldados con éxito.")
	case 4:
		backupFolder := os.Args[1]
		file := os.Args[2]
		restoredFolder := os.Args[3]

		err := os.MkdirAll(restoredFolder, 0755)
		if err != nil {
			fmt.Println("Error al crear la carpeta de restauración:", err)
			return
		}

		err = RestoreBackup(backupFolder, file, restoredFolder)
		if err != nil {
			fmt.Println(err)
			return
		}

	default:
		fmt.Println("Número incorrecto de argumentos.")
		return
	}
}
