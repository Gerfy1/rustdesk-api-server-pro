package api

import (
	"os"
	"path/filepath"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

type DownloadController struct {
	basicController
}

// BeforeActivation registers custom routes
func (c *DownloadController) BeforeActivation(b mvc.BeforeActivation) {
	b.Handle("GET", "/download/list", "HandleList")
	b.Handle("GET", "/download/windows", "HandleWindows")
	b.Handle("GET", "/download/macos", "HandleMacos")
	b.Handle("GET", "/download/linux", "HandleLinux")
	b.Handle("GET", "/download/{filename}", "HandleFile")
}

// HandleWindows GET /api/download/windows
func (c *DownloadController) HandleWindows() {
	c.serveInstaller("windows", "rustdesk-windows.exe")
}

// HandleMacos GET /api/download/macos
func (c *DownloadController) HandleMacos() {
	c.serveInstaller("macos", "rustdesk-macos.dmg")
}

// HandleLinux GET /api/download/linux
func (c *DownloadController) HandleLinux() {
	c.serveInstaller("linux", "rustdesk-linux.deb")
}

// HandleFile GET /api/download/{filename}
func (c *DownloadController) HandleFile(filename string) {
	c.serveInstaller("", filename)
}

// serveInstaller serves the installer file
func (c *DownloadController) serveInstaller(platform, defaultFilename string) {
	installersPath := "./data/installers"

	var filePath string

	if platform != "" {
		// Try to find any file matching the platform
		files, err := os.ReadDir(installersPath)
		if err != nil {
			c.Ctx.StatusCode(iris.StatusNotFound)
			c.Ctx.JSON(iris.Map{"error": "Installers directory not found"})
			return
		}

		for _, file := range files {
			name := file.Name()
			switch platform {
			case "windows":
				if filepath.Ext(name) == ".exe" || filepath.Ext(name) == ".msi" {
					filePath = filepath.Join(installersPath, name)
					break
				}
			case "macos":
				if filepath.Ext(name) == ".dmg" || filepath.Ext(name) == ".pkg" {
					filePath = filepath.Join(installersPath, name)
					break
				}
			case "linux":
				ext := filepath.Ext(name)
				if ext == ".deb" || ext == ".rpm" || ext == ".AppImage" {
					filePath = filepath.Join(installersPath, name)
					break
				}
			}
		}
	} else {
		filePath = filepath.Join(installersPath, defaultFilename)
	}

	if filePath == "" {
		c.Ctx.StatusCode(iris.StatusNotFound)
		c.Ctx.JSON(iris.Map{"error": "Installer not found for " + platform})
		return
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.Ctx.StatusCode(iris.StatusNotFound)
		c.Ctx.JSON(iris.Map{"error": "Installer file not found: " + filepath.Base(filePath)})
		return
	}

	// Set Content-Disposition header to force download with correct filename
	filename := filepath.Base(filePath)
	c.Ctx.Header("Content-Disposition", "attachment; filename=\""+filename+"\"")

	// Serve the file
	c.Ctx.ServeFile(filePath)
}

// HandleList GET /api/download/list - lista os instaladores disponíveis
func (c *DownloadController) HandleList() {
	installersPath := "./data/installers"

	files, err := os.ReadDir(installersPath)
	if err != nil {
		c.Ctx.JSON(iris.Map{
			"installers": []string{},
		})
		return
	}

	type InstallerInfo struct {
		Name     string `json:"name"`
		Platform string `json:"platform"`
		Size     int64  `json:"size"`
		URL      string `json:"url"`
	}

	var installers []InstallerInfo

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		name := file.Name()
		ext := filepath.Ext(name)

		var platform string
		switch ext {
		case ".exe", ".msi":
			platform = "windows"
		case ".dmg", ".pkg":
			platform = "macos"
		case ".deb", ".rpm", ".AppImage":
			platform = "linux"
		default:
			continue
		}

		info, _ := file.Info()
		size := int64(0)
		if info != nil {
			size = info.Size()
		}

		installers = append(installers, InstallerInfo{
			Name:     name,
			Platform: platform,
			Size:     size,
			URL:      "/api/download/" + name,
		})
	}

	c.Ctx.JSON(iris.Map{
		"installers": installers,
	})
}
