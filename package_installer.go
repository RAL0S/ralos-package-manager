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
	installTesting bool
	pkgInfo        PackageInfo
	tmpDir         string
	pkgCloneDir    string
	cfg            *APMConfig
}

func (pi *PackageInstaller) New(cfg *APMConfig, pkgInfo PackageInfo, installTesting bool) {
	pi.cfg = cfg
	pi.pkgInfo = pkgInfo
	pi.installTesting = installTesting
}

func (pi *PackageInstaller) bootstrap() bool {
	var err error
	pi.tmpDir, err = os.MkdirTemp("", "apm_tmp")
	if err != nil {
		log.Println("Failed to create temporary directory:", err)
		return false
	}

	pi.pkgCloneDir = filepath.Join(pi.cfg.InstallPath, CLONE_DIR_NAME, pi.pkgInfo.Name+"-"+pi.pkgInfo.Version)
	os.RemoveAll(pi.pkgCloneDir)

	if err = os.Mkdir(pi.pkgCloneDir, 0755); err != nil {
		log.Println("Failed to create clone directory:", err)
		return false
	}

	log.Println("Cloning package repo to", pi.pkgCloneDir)

	var ref plumbing.ReferenceName
	if !pi.installTesting {
		ref = plumbing.NewTagReferenceName(pi.pkgInfo.RepoTag)
	} else {
		ref = plumbing.NewBranchReferenceName("testing")
	}

	_, err = git.PlainClone(pi.pkgCloneDir, false, &git.CloneOptions{
		URL:           pi.pkgInfo.RepoUrl,
		Depth:         1,
		SingleBranch:  true,
		ReferenceName: ref,
	})

	if err != nil {
		log.Println("Failed to clone package repository:", err)
		os.RemoveAll(pi.pkgCloneDir)
		return false
	}
	return true
}

func (pi *PackageInstaller) processPackageDetails() *PackageDetail {
	pkgTomlPath := filepath.Join(pi.pkgCloneDir, "package.toml")
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
	if !pi.bootstrap() {
		log.Println("Failed while bootstrapping")
		return false
	}

	pkgDetails := pi.processPackageDetails()
	if pkgDetails == nil {
		log.Println("Failed while processing package details")
		return false
	}

	installScriptPath := filepath.Join(pi.pkgCloneDir, pkgDetails.InstallScript)
	if _, err := os.Stat(installScriptPath); err != nil {
		log.Println("Package installer script not found:", installScriptPath)
		return false
	}

	// Make install script executable
	os.Chmod(installScriptPath, 0744)

	pkgInstallDir := filepath.Join(pi.cfg.InstallPath, PKG_DIR_NAME, pkgDetails.Name+"-"+pkgDetails.Version)
	if _, err := os.Stat(pkgInstallDir); err == nil {
		os.RemoveAll(pkgInstallDir)
	}
	err := os.Mkdir(pkgInstallDir, 0755)
	if err != nil {
		log.Println("Failed to create directory:", pkgInstallDir)
		return false
	}

	cmd := exec.Command(installScriptPath, "install")
	cmd.Env = append(
		os.Environ(),
		"APM_TMP_DIR="+pi.tmpDir,
		"APM_PKG_INSTALL_DIR="+pkgInstallDir,
		"APM_PKG_BIN_DIR="+filepath.Join(pi.cfg.InstallPath, BIN_DIR_NAME),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	log.Println("Running package installation script, please wait...")
	err = cmd.Run()

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
