package main

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the top-level structure of config.yaml
type Config struct {
	Title   string      `yaml:"title"`
	Company Company     `yaml:"company"`
	Login   LoginConfig `yaml:"login"`
}

// LoginConfig represents the social login settings
type LoginConfig struct {
	Title     string     `yaml:"title"`
	Providers []Provider `yaml:"providers"`
}

// Provider represents a single social login provider
type Provider struct {
	Name  string `yaml:"name"`
	Icon  string `yaml:"icon"`
	URL   string `yaml:"url"`
	Class string `yaml:"class"`
}

// Company represents the nested company details
type Company struct {
	ShortName string `yaml:"companyShortName"`
	LongName  string `yaml:"companyLongName"`
	Address   string `yaml:"companyAddress"`
	Email     string `yaml:"companyEmail"`
	Phone     string `yaml:"companyPhone"`
	Year      string `yaml:"companyYear"`
}

func main() {
	// 1. Read config.yaml
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		fmt.Printf("Error reading config.yaml: %v\n", err)
		os.Exit(1)
	}

	var config Config
	if err := yaml.Unmarshal(configFile, &config); err != nil {
		fmt.Printf("Error parsing config.yaml: %v\n", err)
		os.Exit(1)
	}

	// 2. Prepare Web Directory
	if err := os.MkdirAll("web", 0755); err != nil {
		fmt.Printf("Error creating web directory: %v\n", err)
		os.Exit(1)
	}

	// 3. Parse Template
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Printf("Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// 4. Create Output File
	outputFile, err := os.Create("web/index.html")
	if err != nil {
		fmt.Printf("Error creating web/index.html: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// 5. Execute Template
	if err := tmpl.Execute(outputFile, config); err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}

	// 6. Copy Static Assets (from src to web/static)
	// Create web/static directory first
	if err := os.MkdirAll("web/static", 0755); err != nil {
		fmt.Printf("Error creating web/static directory: %v\n", err)
		os.Exit(1)
	}

	if err := copyDir("src", "web/static"); err != nil {
		fmt.Printf("Error copying static assets from src: %v\n", err)
		os.Exit(1)
	}

	// 7. Generate Auth Pages
	if len(config.Login.Providers) > 0 {
		if err := generateAuthPages(config.Login.Providers); err != nil {
			fmt.Printf("Error generating auth pages: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("Successfully generated web/index.html, auth pages, and copied assets.")
}

func generateAuthPages(providers []Provider) error {
	authDir := "web/auth"
	if err := os.MkdirAll(authDir, 0755); err != nil {
		return err
	}

	// Simple template for auth pages
	const authPageTmpl = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Login with {{.Name}}</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css">
    <link rel="stylesheet" href="../static/css/style.css">
    <style>
        body { display: flex; align-items: center; justify-content: center; height: 100vh; text-align: center; }
        .spinner { font-size: 3rem; color: var(--accent-color); margin-bottom: 1rem; }
    </style>
</head>
<body>
    <div class="card">
        <div class="spinner"><i class="fa-solid fa-circle-notch fa-spin"></i></div>
        <h2>Connection to {{.Name}}...</h2>
        <p>This is a static site demo.</p>
        <br>
        <a href="../index.html" class="social-btn" style="display: inline-flex;">
            <i class="fa-solid fa-arrow-left"></i> Back to Home
        </a>
    </div>
    <script>
        // Check for theme preference
        const currentTheme = localStorage.getItem('theme') || 
             (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
        if (currentTheme === 'dark') document.documentElement.setAttribute('data-theme', 'dark');
    </script>
</body>
</html>`

	tmpl, err := template.New("auth").Parse(authPageTmpl)
	if err != nil {
		return err
	}

	for _, provider := range providers {
		// Filename: Google -> web/auth/Google.html
		filename := filepath.Join(authDir, provider.Name+".html")
		f, err := os.Create(filename)
		if err != nil {
			return err
		}

		err = tmpl.Execute(f, provider)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// copyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func copyDir(src string, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Create destination if it doesn't exist (copying permissions of source)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dst, si.Mode())
		if err != nil {
			return err
		}
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			// Skip symlinks
			if entry.Type()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
