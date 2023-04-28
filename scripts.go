package dcshmd

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed scripts/*
var scripts embed.FS

const (
	// exportLuaFileName is a name of lua file that should be modified
	exportLuaFileName = "Export.lua"

	// exportLuaLine is a line to be added to Exports.lua file
	exportLuaLine = "local lfs=require('lfs');dofile(lfs.writedir()..'Scripts/DCSHMD/Export.lua')"
)

// InstallScripts installs the scripts in the specified target directory.
// If verbose is true, the function prints detailed log messages to stdout.
func InstallScripts(scriptsInstallDir string, verbose bool) error {
	// check if scriptsInstallDir exists
	if _, err := os.Stat(scriptsInstallDir); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: '%s'", scriptsInstallDir)
	}

	// replace the scripts in the target directory
	if err := replaceScripts(scriptsInstallDir, verbose); err != nil {
		return err
	}

	// update the Export.lua script in the target directory
	return updateExportScript(scriptsInstallDir, verbose)
}

// replaceScripts replaces the scripts in the specified target directory with the embedded scripts.
// If verbose is true, the function prints detailed log messages to stdout.
func replaceScripts(scriptsInstallDir string, verbose bool) error {
	// get the embedded file system containing the contents of the scripts directory
	scripts, err := fs.Sub(scripts, "scripts")
	if err != nil {
		return fmt.Errorf("failed to read internal scripts directory: %w", err)
	}

	// walk the embedded file system
	return fs.WalkDir(scripts, ".", func(path string, d fs.DirEntry, err error) error {
		// skip root directory to avoid deleting all files and folders in target directory
		if path == "." {
			return nil
		}

		// construct the target path by joining the target directory with the current path
		targetPath := filepath.Join(scriptsInstallDir, path)

		// check if current entry is a directory
		if d.IsDir() {
			if verbose {
				fmt.Printf("removing directory '%s'...\n", targetPath)
			}
			// remove target directory if it exists to clean up outdated scripts
			if err := os.RemoveAll(targetPath); err != nil {
				return fmt.Errorf("failed to clean up directory '%s': %w", targetPath, err)
			}

			if verbose {
				fmt.Printf("creating directory '%s'...\n", targetPath)
			}
			// create target directory
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", targetPath, err)
			}
		} else {
			// current entry is a file
			// read file content from embedded file system
			data, err := fs.ReadFile(scripts, path)
			if err != nil {
				return fmt.Errorf("failed to read internal file '%s': %w", path, err)
			}

			if verbose {
				fmt.Printf("writing file '%s'...\n", targetPath)
			}
			// write file content to target directory
			if err := os.WriteFile(targetPath, data, 0644); err != nil {
				return fmt.Errorf("failed to copy internal file '%s' to '%s': %w", path, targetPath, err)
			}
		}

		// continue walking
		return nil
	})
}

// updateExportScript updates the Export.lua script in the specified target directory.
// If verbose is true, the function prints detailed log messages to stdout.
func updateExportScript(scriptsInstallDir string, verbose bool) error {
	// construct the path to the Export.lua file
	exportFile := filepath.Join(scriptsInstallDir, exportLuaFileName)

	// check if the Export.lua file exists
	if _, err := os.Stat(exportFile); os.IsNotExist(err) {
		if verbose {
			fmt.Printf("creating script file '%s'...\n", exportFile)
		}
		// file does not exist, create it
		if err := os.WriteFile(exportFile, []byte(fmt.Sprintln(exportLuaLine)), 0644); err != nil {
			return fmt.Errorf("failed to create '%s' file: %w", exportFile, err)
		}
	} else {
		if verbose {
			fmt.Printf("checking script file '%s'...\n", exportFile)
		}
		// file exists, check if line is present
		data, err := os.ReadFile(exportFile)
		if err != nil {
			return fmt.Errorf("failed to read the contents of '%s' script: %w", exportFile, err)
		}

		// check if data contains the line
		if !bytes.Contains(data, []byte(exportLuaLine)) {
			if verbose {
				fmt.Printf("prepending line into script file '%s'...\n", exportFile)
			}
			// line is not present, prepend it
			data = append([]byte(fmt.Sprintln(exportLuaLine)), data...)
			if err := os.WriteFile(exportFile, data, 0644); err != nil {
				return fmt.Errorf("failed to prepend line in '%s' script: %w", exportFile, err)
			}
		}
	}

	return nil
}
