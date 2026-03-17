package folder

import (
	"embed"
	"fmt"
	"os"
	"text/template"
)

// TemplateFile represents a file to be created from a template and data.
type TemplateFile struct {
	// Name of the file to be created.
	Filename string
	// The template to execute into the new file.
	Tmpl *template.Template
	// A struct representing the data that will be used to complete the template.
	TmplData any
}

// CreateFolder creates a folder with the given folder name, writing each template into a new
// file, and each embedded file in the given directory into a new file.
func CreateFolder(folderName string, templates []*TemplateFile, dir string, f *embed.FS) error {
	if err := os.Mkdir(folderName, 0o751); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", folderName, err)
	}

	if err := os.Chdir(folderName); err != nil {
		return fmt.Errorf("failed to change directory to project directory: %w", err)
	}

	if templates != nil {
		if err := createTemplates(templates); err != nil {
			return err
		}
	}

	if f != nil {
		if err := createFiles(dir, f); err != nil {
			return err
		}
	}

	return nil
}

func createTemplates(templates []*TemplateFile) error {
	for _, tmpl := range templates {
		file, err := os.Create(tmpl.Filename)
		if err != nil {
			return fmt.Errorf("failed to create new file '%s' for template: %w", tmpl.Filename, err)
		}

		if err := tmpl.Tmpl.Execute(file, tmpl.TmplData); err != nil {
			return fmt.Errorf("failed to execute template '%s': %w", tmpl.Tmpl.Name(), err)
		}
	}

	return nil
}

func createFiles(dir string, f *embed.FS) error {
	entries, err := f.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read embedded directory 'file': %w", err)
	}

	for _, entry := range entries {
		file, err := os.Create(entry.Name())
		if err != nil {
			return fmt.Errorf("failed to create new file '%s' in project directory: %w", entry.Name(), err)
		}

		data, err := f.ReadFile(dir + "/" + entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read data from embedded file '%s': %w", entry.Name(), err)
		}

		if _, err := file.Write(data); err != nil {
			return fmt.Errorf("failed to write data from embedded file '%s' to new file '%s': %w", entry.Name(), file.Name(), err)
		}
	}

	return nil
}
