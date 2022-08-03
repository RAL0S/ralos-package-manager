package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pelletier/go-toml/v2"
)

type PackageDetail struct {
	Name              string `toml:"name"`
	Description       string `toml:"description"`
	Version           string `toml:"version"`
	SourceUrl         string `toml:"source_url"`
	License           string `toml:"license"`
	InstallScript     string `toml:"install_script"`
	InstallScriptType string `toml:"install_script_type"`
}

type PackageDescriptor struct {
	PkgDetail PackageDetail `toml:"package_detail"`
}

type PackageInstaller struct {
	pkgInfo PackageInfo
	tmpDir  string
	cfg     *APMConfig
}

func (pi *PackageInstaller) New(cfg *APMConfig, pkgInfo PackageInfo) {
	pi.cfg = cfg
	pi.pkgInfo = pkgInfo
}

func (pi *PackageInstaller) bootStrap() bool {
	var err error
	pi.tmpDir, err = os.MkdirTemp("", "apm_tmp")
	if err != nil {
		log.Println("Failed to create temporary directory:", err)
		return false
	}

	log.Println("Cloning package repo to", pi.tmpDir)
	_, err = git.PlainClone(pi.tmpDir, false, &git.CloneOptions{
		URL:           pi.pkgInfo.RepoUrl,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: plumbing.NewTagReferenceName(pi.pkgInfo.RepoTag),
	})
	if err != nil {
		log.Println("Failed to clone package repository:", err)
		return false
	}
	return true
}

func (pi *PackageInstaller) processPackageDetails() *PackageDetail {
	pkgTomlPath := filepath.Join(pi.tmpDir, "package.toml")
	pkgTomlBytes, err := os.ReadFile(pkgTomlPath)
	if err != nil {
		log.Println("Failed to read package config file", pkgTomlPath)
		return nil
	}

	var pkgDescriptor PackageDescriptor
	if err := toml.Unmarshal(pkgTomlBytes, &pkgDescriptor); err != nil {
		log.Println("Failed to parse config file", pkgTomlPath)
		return nil
	}

	//TODO: Implement more install script types
	if pkgDescriptor.PkgDetail.InstallScriptType != "shell" {
		log.Println("install_script_type = shell is only supported presently")
		return nil
	}
	return &pkgDescriptor.PkgDetail
}

func (pi *PackageInstaller) cleanup() {
	log.Println("Cleaning up")
	os.RemoveAll(pi.tmpDir)
}

func (pi *PackageInstaller) Install() bool {
	if !pi.bootStrap() {
		log.Println("Failed while bootstrapping")
		return false
	}

	pkgDetails := pi.processPackageDetails()
	if pkgDetails == nil {
		log.Println("Failed while processing package details")
		return false
	}

	installScriptPath := filepath.Join(pi.tmpDir, pkgDetails.InstallScript)
	if _, err := os.Stat(installScriptPath); err != nil {
		log.Println("Package installer script not found:", installScriptPath)
		return false
	}
	// Make install script executable
	os.Chmod(installScriptPath, 0744)

	pkgInstallDir := filepath.Join(pi.cfg.InstallPath, "packages", pkgDetails.Name+"-"+pkgDetails.Version)
	if _, err := os.Stat(pkgInstallDir); err == nil {
		os.RemoveAll(pkgInstallDir)
	}
	err := os.Mkdir(pkgInstallDir, 0755)
	if err != nil {
		log.Println("Failed to create directory:", pkgInstallDir)
		return false
	}

	cmd := exec.Command(installScriptPath)
	cmd.Env = append(
		os.Environ(),
		"APM_TMP_DIR="+pi.tmpDir,
		"APM_PKG_INSTALL_DIR="+pkgInstallDir,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("Running package installation script, please wait...")
	err = cmd.Run()	

	// logFileName := fmt.Sprintf("%s-%s-%d.log", pkgDetails.Name, pkgDetails.Version, time.Now().Unix())
	// os.WriteFile(logFileName, output, 0644)

	if err != nil {
		log.Println("Failed to execute install script")
		return false
	}

	if cmd.ProcessState.ExitCode() != 0 {
		log.Println("Install script returned non zero status code")
		return false
	}
	pi.cleanup()
	return true
}
