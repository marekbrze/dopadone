package version

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/marekbrze/dopadone/internal/constants"
)

var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if Version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && GitCommit == "unknown" {
				GitCommit = setting.Value
				if len(setting.Value) > 7 {
					GitCommit = setting.Value[:7]
				}
			}
			if setting.Key == "vcs.time" && BuildDate == "unknown" {
				BuildDate = setting.Value
			}
		}
	}
}

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Assets  []struct {
		Name string `json:"name"`
		URL  string `json:"browser_download_url"`
	} `json:"assets"`
}

func BuildInfo() string {
	return fmt.Sprintf("dopa %s\n  Git commit: %s\n  Build date: %s", Version, GitCommit, BuildDate)
}

func CheckLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get("https://api.github.com/repos/marekbrze/dopadone/releases/latest")
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("no releases found")
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	return &release, nil
}

func IsUpdateAvailable() (bool, *GitHubRelease, error) {
	if Version == "dev" {
		return false, nil, nil
	}

	release, err := CheckLatestRelease()
	if err != nil {
		return false, nil, err
	}

	currentVersion := strings.TrimPrefix(Version, "v")
	latestVersion := strings.TrimPrefix(release.TagName, "v")

	return currentVersion != latestVersion, release, nil
}

func getPlatformString() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	if os == "darwin" && arch == "arm64" {
		return "darwin-arm64"
	}
	return fmt.Sprintf("%s-%s", os, arch)
}

func getAssetName() string {
	platform := getPlatformString()
	if runtime.GOOS == constants.OSWindows {
		return fmt.Sprintf("dopa-%s.zip", platform)
	}
	return fmt.Sprintf("dopa-%s.tar.gz", platform)
}

func findAssetURL(release *GitHubRelease, assetName string) (string, error) {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			return asset.URL, nil
		}
	}
	return "", fmt.Errorf("asset %s not found in release", assetName)
}

func downloadFile(url, dest string) error {
	client := &http.Client{Timeout: 5 * time.Minute}

	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, destDir string) (string, error) {
	binaryName := "dopa"
	if runtime.GOOS == constants.OSWindows {
		binaryName = "dopa.exe"
	}

	if strings.HasSuffix(archivePath, ".zip") {
		r, err := zip.OpenReader(archivePath)
		if err != nil {
			return "", err
		}
		defer r.Close()

		for _, f := range r.File {
			if strings.HasSuffix(f.Name, ".exe") || (strings.HasSuffix(f.Name, "dopa") && !strings.Contains(f.Name, "/")) {
				destPath := filepath.Join(destDir, binaryName)
				outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
				if err != nil {
					return "", err
				}
				defer outFile.Close()

				rc, err := f.Open()
				if err != nil {
					return "", err
				}
				defer rc.Close()

				_, err = io.Copy(outFile, rc)
				return destPath, err
			}
		}
		return "", fmt.Errorf("binary not found in archive")
	}

	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if strings.HasSuffix(hdr.Name, "dopa") && hdr.Typeflag == tar.TypeReg {
			destPath := filepath.Join(destDir, binaryName)
			outFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return "", err
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tr); err != nil {
				return "", err
			}
			return destPath, nil
		}
	}

	return "", fmt.Errorf("binary not found in archive")
}

func getCurrentBinaryPath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	resolved, err := filepath.EvalSymlinks(execPath)
	if err != nil {
		return execPath, nil
	}
	return resolved, nil
}

func replaceBinary(newBinary, currentBinary string) error {
	if runtime.GOOS == constants.OSWindows {
		oldPath := currentBinary + ".old"
		if err := os.Rename(currentBinary, oldPath); err != nil {
			return fmt.Errorf("failed to rename old binary: %w", err)
		}
		if err := copyFile(newBinary, currentBinary); err != nil {
			os.Rename(oldPath, currentBinary)
			return fmt.Errorf("failed to copy new binary: %w", err)
		}
		os.Remove(oldPath)
		return nil
	}

	if err := os.Rename(newBinary, currentBinary); err != nil {
		if err := copyFile(newBinary, currentBinary); err != nil {
			return fmt.Errorf("failed to replace binary: %w", err)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

type UpgradeOptions struct {
	DBPath      string
	SkipMigrate bool
}

func PerformUpgrade(opts UpgradeOptions) error {
	fmt.Println("Checking for updates...")

	available, release, err := IsUpdateAvailable()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !available {
		fmt.Printf("Already running the latest version: %s\n", Version)
		return nil
	}

	fmt.Printf("Update available: %s -> %s\n", Version, release.TagName)

	assetName := getAssetName()
	assetURL, err := findAssetURL(release, assetName)
	if err != nil {
		return fmt.Errorf("failed to find download for your platform (%s): %w", getPlatformString(), err)
	}

	tempDir, err := os.MkdirTemp("", "dopa-upgrade-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, assetName)
	fmt.Printf("Downloading %s...\n", assetName)

	if err := downloadFile(assetURL, archivePath); err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}

	fmt.Println("Extracting binary...")
	newBinary, err := extractBinary(archivePath, tempDir)
	if err != nil {
		return fmt.Errorf("failed to extract binary: %w", err)
	}

	currentBinary, err := getCurrentBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to locate current binary: %w", err)
	}

	fmt.Printf("Replacing binary at %s...\n", currentBinary)
	if err := replaceBinary(newBinary, currentBinary); err != nil {
		return fmt.Errorf("failed to replace binary: %w", err)
	}

	fmt.Printf("Successfully upgraded to %s\n", release.TagName)

	if !opts.SkipMigrate && opts.DBPath != "" {
		fmt.Println("\nRunning database migrations...")
		if err := runMigrations(currentBinary, opts.DBPath); err != nil {
			fmt.Printf("Warning: migrations failed: %v\n", err)
			fmt.Println("You may need to run migrations manually: dopa migrate up")
		} else {
			fmt.Println("Migrations completed successfully.")
		}
	}

	fmt.Println("\nUpgrade complete!")
	return nil
}

func runMigrations(binaryPath, dbPath string) error {
	cmd := exec.Command(binaryPath, "migrate", "up", "--db", dbPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
