package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bakayu/lq/internal/provider"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
)

func main() {
	var (
		fileType string
		selected provider.Template
		confirm  bool
		owner    string
	)

	theme := huh.ThemeCatppuccin()

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("File to generate >").
				Options(
					huh.NewOption("Gitignore", ".gitignore"),
					huh.NewOption("License", "LICENSE"),
				).
				Value(&fileType),
		),
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	var prov provider.Provider
	if fileType == ".gitignore" {
		prov = provider.NewGitignoreProvider()
	} else {
		prov = provider.NewLicenseProvider()
	}

	var templates []provider.Template
	action := func() {
		var err error
		templates, err = prov.List()
		if err != nil {
			log.Fatalf("Error fetching templates: %v", err)
		}
	}
	if err := spinner.New().Title("Fetching templates...").Action(action).Run(); err != nil {
		log.Fatal(err)
	}

	var options []huh.Option[provider.Template]
	for _, t := range templates {
		options = append(options, huh.NewOption(t.Name, t))
	}

	fields := []huh.Field{
		huh.NewSelect[provider.Template]().
			Title("Choose a template").
			Options(options...).
			Value(&selected).
			Height(15).
			Filtering(true),
	}

	if fileType == "LICENSE" {
		fields = append(fields,
			huh.NewInput().
				Title("Who is the copyright holder?").
				Placeholder("e.g. John Doe").
				Value(&owner),
		)
	}

	form = huh.NewForm(
		huh.NewGroup(fields...),
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	form = huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Create %s using %s template?", fileType, selected.Name)).
				Value(&confirm),
		),
	).WithTheme(theme)

	if err := form.Run(); err != nil {
		log.Fatal(err)
	}

	if !confirm {
		fmt.Println("Cancelled")
		os.Exit(0)
	}

	content, err := prov.GetContent(selected.Key)
	if err != nil {
		log.Fatalf("Error fetching content: %v", err)
	}

	if fileType == "LICENSE" {
		content = fillLicensePlaceholders(content, owner)
	}

	fileName := fileType
	if err := os.WriteFile(fileName, []byte(content), 0644); err != nil {
		log.Fatalf("Error writing file: %v", err)
	}

	fmt.Printf("\nSuccessfully created %s!\n", fileName)
}

// fillLicensePlaceholders replaces standard GitHub API placeholders
func fillLicensePlaceholders(text, owner string) string {
	year := fmt.Sprintf("%d", time.Now().Year())

	// Standard GitHub placeholders
	text = strings.ReplaceAll(text, "[year]", year)
	text = strings.ReplaceAll(text, "[fullname]", owner)

	// Common variations found in some raw templates
	text = strings.ReplaceAll(text, "<year>", year)
	text = strings.ReplaceAll(text, "<copyright holders>", owner)

	return text
}
