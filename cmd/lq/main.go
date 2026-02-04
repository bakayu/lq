package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/huh"
)

func main() {
	var (
		fileType string
		template string
		confirm  bool
	)
	err := huh.NewSelect[string]().
		Title("File to generate >").
		Options(
			huh.NewOption("Gitignore", ".gitignore"),
			huh.NewOption("License", "LICENSE"),
		).
		Value(&fileType).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	var options []huh.Option[string]

	if fileType == ".gitignore" {
		options = []huh.Option[string]{
			huh.NewOption("Go", "Go"),
			huh.NewOption("Rust", "Rust"),
			huh.NewOption("Java", "Java"),
			huh.NewOption("Python", "Python"),
		}
	} else {
		options = []huh.Option[string]{
			huh.NewOption("MIT", "MIT"),
			huh.NewOption("Apache 2.0", "Apache-2.0"),
			huh.NewOption("GNU GPLv3", "GPL-3.0"),
		}
	}

	err = huh.NewSelect[string]().
		Title("Choose a template").
		Options(options...).
		Value(&template).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	err = huh.NewConfirm().
		Title(fmt.Sprintf("Creating %s using %s template", fileType, template)).
		Value(&confirm).
		Run()

	if err != nil {
		log.Fatal(err)
	}

	if confirm {
		fmt.Println("Success")
	} else {
		fmt.Println("Cancelled")
		os.Exit(0)
	}
}
