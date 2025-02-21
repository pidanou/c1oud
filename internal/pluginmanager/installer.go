package pluginmanager

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/pidanou/c1-core/internal/constants"
	"github.com/pidanou/c1-core/pkg/plugin"
)

func downloadFromVCS(plug *plugin.Plugin) error {
	downloadPath := path.Join(constants.Envs["C1_DIR"], plug.Name)
	err := os.MkdirAll(downloadPath, 0755)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", plug.URI, ".")
	cmd.Dir = downloadPath
	err = cmd.Run()
	if err != nil {
		return err
	}

	if plug.InstallCommand != "" {
		cmd := exec.Command("sh", "-c", plug.InstallCommand)
		cmd.Dir = downloadPath
		err := cmd.Run()
		if err != nil {
			return err
		}
		return err
	}

	return nil
}

func updateFromVCS(plug *plugin.Plugin) error {
	if plug.UpdateCommand != "" {
		cmd := exec.Command("sh", "-c", plug.UpdateCommand)
		cmd.Dir = path.Join(constants.Envs["C1_DIR"], plug.Name)
		if err := cmd.Run(); err != nil {
			return err
		}
		return nil
	}

	repoPath := path.Join(constants.Envs["C1_DIR"], plug.Name)

	cmd := exec.Command("git", "pull", "origin", "HEAD")
	cmd.Dir = repoPath
	err := cmd.Run()
	if err != nil {
		return err
	}

	return err
}

func downloadFromHTTP(plug *plugin.Plugin) error {
	parts := strings.Split(plug.URI, "/")
	resourceName := parts[len(parts)-1]
	downloadDir := path.Join(constants.Envs["C1_DIR"], plug.Name)
	err := os.MkdirAll(downloadDir, 0755)
	out, err := os.Create(path.Join(downloadDir, resourceName))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(plug.URI)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	if plug.InstallCommand != "" {
		cmd := exec.Command("sh", "-c", plug.InstallCommand)
		cmd.Dir = downloadDir
		err := cmd.Run()
		if err != nil {
			return nil
		}
		return err
	}

	return nil
}

func updateFromHTTP(plug *plugin.Plugin) error {
	if plug.UpdateCommand == "" {
		err := DeletePlugin(plug.Name)
		if err != nil {
			return err
		}
		return downloadFromHTTP(plug)
	}

	cmd := exec.Command("sh", "-c", plug.UpdateCommand)
	cmd.Dir = path.Join(constants.Envs["C1_DIR"], plug.Name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func downloadFromLocal(plug *plugin.Plugin) error {
	downloadDir := path.Join(constants.Envs["C1_DIR"], plug.Name)
	err := os.MkdirAll(downloadDir, 0755)
	if err != nil {
		return err
	}

	parts := strings.Split(plug.URI, "/")
	fileName := parts[len(parts)-1]

	info, err := os.Stat(plug.URI)
	if err != nil {
		return err
	}

	if info.IsDir() {
		err := copyDir(plug.URI, downloadDir)
		if err != nil {
			return err
		}
	}

	if info.Mode().IsRegular() {
		err := copyFile(plug.URI, path.Join(downloadDir, fileName))
		if err != nil {
			return err
		}
	}

	if plug.InstallCommand != "" {
		cmd := exec.Command("sh", "-c", plug.InstallCommand)
		cmd.Dir = downloadDir
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			if exitErr, ok := err.(*exec.ExitError); ok {
				fmt.Println("Exit Code:", exitErr.ExitCode())
			}
			fmt.Println("Output:", string(out))
			return nil
		}
		return err
	}

	return nil
}

func updateFromLocal(plug *plugin.Plugin) error {
	if plug.UpdateCommand == "" {
		err := DeletePlugin(plug.Name)
		if err != nil {
			return err
		}
		return downloadFromLocal(plug)
	}

	cmd := exec.Command("sh", "-c", plug.UpdateCommand)
	cmd.Dir = path.Join(constants.Envs["C1_DIR"], plug.Name)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	err = os.MkdirAll(dst, 0755)
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

	err = out.Sync()
	if err != nil {
		return err
	}

	return nil
}

func DeletePlugin(name string) error {
	return os.RemoveAll(path.Join(constants.Envs["C1_DIR"], name))
}
