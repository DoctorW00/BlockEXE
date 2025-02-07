package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Bitte einen Verzeichnis-Pfad angeben oder den Ordner auf die App ziehen.")
		return
	}

	root := os.Args[1]

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".exe" {
			fmt.Println("Gefunden: ", path)
			blockInFirewall(path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Fehler beim Durchsuchen des Verzeichnisses: %v\n", err)
	}
}

func blockInFirewall(exePath string) {
	cmdIn := exec.Command("netsh", "advfirewall", "firewall", "add", "rule", "name=\"BlockIn"+filepath.Base(exePath)+"\"", "dir=in", "action=block", "program="+exePath, "enable=yes")
	if err := cmdIn.Run(); err != nil {
		fmt.Printf("Fehler beim Blockieren der eingehenden Verbindung: %v\n", err)
	}

	cmdOut := exec.Command("netsh", "advfirewall", "firewall", "add", "rule", "name=\"BlockOut"+filepath.Base(exePath)+"\"", "dir=out", "action=block", "program="+exePath, "enable=yes")
	if err := cmdOut.Run(); err != nil {
		fmt.Printf("Fehler beim Blockieren der ausgehenden Verbindung: %v\n", err)
	}
}
