package dcshmd

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

//go:embed scripts/*
var scripts embed.FS

const (
	// exportLuaFileName is a name of lua file that should be modified
	exportLuaFileName = "Export.lua"

	// exportLuaLine is a line to be added to Exports.lua file
	exportLuaLine = "local lfs=require('lfs');dofile(lfs.writedir()..'Scripts/DCSHMD/Export.lua')"
)

// scriptsFS returns a sub-filesystem of the embedded scripts directory.
func scriptsFS() (fs.FS, error) {
	return fs.Sub(scripts, "scripts")
}

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

func replaceScripts(scriptsInstallDir string, verbose bool) error {
	scripts, err := scriptsFS()
	if err != nil {
		return fmt.Errorf("failed to read embedded scripts directory: %w", err)
	}
	return fs.WalkDir(scripts, ".", func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}
		targetPath := filepath.Join(scriptsInstallDir, path)
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
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", targetPath, err)
			}
		} else {
			data, err := fs.ReadFile(scripts, path)
			if err != nil {
				return fmt.Errorf("failed to read embedded file '%s': %w", path, err)
			}
			if verbose {
				fmt.Printf("writing file '%s'...\n", targetPath)
			}
			if err := os.WriteFile(targetPath, data, 0644); err != nil {
				return fmt.Errorf("failed to copy embedded file '%s' to '%s': %w", path, targetPath, err)
			}
		}
		return nil
	})
}

func updateExportScript(scriptsInstallDir string, verbose bool) error {
	exportFile := filepath.Join(scriptsInstallDir, exportLuaFileName)
	if _, err := os.Stat(exportFile); os.IsNotExist(err) {
		if verbose {
			fmt.Printf("creating script file '%s'...\n", exportFile)
		}
		if err := os.WriteFile(exportFile, []byte(fmt.Sprintln(exportLuaLine)), 0644); err != nil {
			return fmt.Errorf("failed to create '%s' file: %w", exportFile, err)
		}
	} else {
		if verbose {
			fmt.Printf("checking script file '%s'...\n", exportFile)
		}
		data, err := os.ReadFile(exportFile)
		if err != nil {
			return fmt.Errorf("failed to read the contents of '%s' script: %w", exportFile, err)
		}
		if !bytes.Contains(data, []byte(exportLuaLine)) {
			if verbose {
				fmt.Printf("prepending line into script file '%s'...\n", exportFile)
			}
			data = append([]byte(fmt.Sprintln(exportLuaLine)), data...)
			if err := os.WriteFile(exportFile, data, 0644); err != nil {
				return fmt.Errorf("failed to prepend line in '%s' script: %w", exportFile, err)
			}
		}
	}

	return nil
}

// UninstallScripts uninstalls the scripts from the specified target directory.
// If verbose is true, the function prints detailed log messages to stdout.
func UninstallScripts(scriptsInstallDir string, verbose bool) error {
	// Check if the scriptsInstallDir exists.
	if _, err := os.Stat(scriptsInstallDir); os.IsNotExist(err) {
		return fmt.Errorf("folder does not exist: '%s'", scriptsInstallDir)
	}

	// Update the Export.lua script in the target directory by calling the removeExportScriptLine function.
	if err := removeExportScriptLine(scriptsInstallDir, verbose); err != nil {
		return err
	}

	// Remove the scripts in the target directory by calling the removeScripts function.
	return removeScripts(scriptsInstallDir, verbose)
}

// removeScripts removes the scripts from the specified target directory.
// If verbose is true, the function prints detailed log messages to stdout.
func removeScripts(scriptsInstallDir string, verbose bool) error {
	// Read the embedded scripts directory using the scriptsFS function.
	scripts, err := scriptsFS()
	if err != nil {
		return fmt.Errorf("failed to read embedded scripts directory: %w", err)
	}

	// Walk through the embedded scripts directory.
	return fs.WalkDir(scripts, ".", func(path string, d fs.DirEntry, err error) error {
		if path == "." {
			return nil
		}
		targetPath := filepath.Join(scriptsInstallDir, path)
		if d.IsDir() {
			if verbose {
				fmt.Printf("removing directory '%s'...\n", targetPath)
			}
			// Remove the target directory if it exists to clean up outdated scripts.
			if err := os.RemoveAll(targetPath); err != nil {
				return fmt.Errorf("failed to clean up directory '%s': %w", targetPath, err)
			}
			return fs.SkipDir
		} else {
			if verbose {
				fmt.Printf("removing file '%s'...\n", targetPath)
			}
			// Remove the target file if it exists.
			if err := os.Remove(targetPath); err != nil && !os.IsNotExist(err) {
				return fmt.Errorf("failed to remove file '%s': %w", targetPath, err)
			}
		}
		return nil
	})
}

// removeExportScriptLine removes the exportLuaLine from the Export.lua script in the specified target directory.
// If verbose is true, the function prints detailed log messages to stdout.
func removeExportScriptLine(scriptsInstallDir string, verbose bool) error {
	exportFile := filepath.Join(scriptsInstallDir, exportLuaFileName)
	if _, err := os.Stat(exportFile); os.IsNotExist(err) {
		return nil
	}

	if verbose {
		fmt.Printf("checking script file '%s'...\n", exportFile)
	}

	lines, err := readLines(exportFile)
	if err != nil {
		return err
	}

	// in-place filtering
	n := 0
	for _, line := range lines {
		if !strings.Contains(line, exportLuaLine) {
			lines[n] = line
			n++
		}
	}
	if n == len(lines) {
		// no lines have been deleted from the slice
		return nil
	}
	lines = lines[:n]

	var output bytes.Buffer
	for _, line := range lines {
		fmt.Fprintln(&output, line)
	}

	if verbose {
		fmt.Printf("updating script file '%s'...\n", exportFile)
	}
	if err := os.WriteFile(exportFile, output.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to update '%s' file: %w", exportFile, err)
	}

	return nil
}

// readLines reads all the lines from the given file and returns them as a slice of strings.
func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open '%s' file: %w", filename, err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan the contents of '%s' script: %w", filename, err)
	}

	return lines, nil
}
