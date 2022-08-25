package git

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func dirExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func getGitDir() string {
	const GIT_DIR = "/.git/"
	wd, _ := os.Getwd()
	// Traverse all the way to root '/'
	for {
		exist, _ := dirExists(wd + GIT_DIR)
		if exist {
			break
		} else {
			if wd == "/" {
				panic(errors.New("no git repo found"))
			}
			wd = filepath.Dir(wd)
		}
	}
	return wd + GIT_DIR
}

// Returns all git branch names in repo
func GetGitBranches() []string {
	gitDir := getGitDir()
	files, _ := ioutil.ReadDir(gitDir + "refs/heads")
	fnames := make([]string, len(files))
	for i, f := range files {
		fnames[i] = f.Name()
	}
	return fnames
}

func GetCurrentBranch() string {
	gitDir := getGitDir()
	dat, _ := os.ReadFile(gitDir + "HEAD")
	raw := string(dat)
	current := strings.TrimSpace(strings.Replace(raw, "ref: refs/heads/", "", 1))
	return current
}

func DeleteGitBranch(b string) (error, string) {
	cmd := exec.Command("git", "branch", "-D", b)
	if stdOutErr, err := cmd.CombinedOutput(); err != nil {
		return err, string(stdOutErr)
	}
	return nil, ""
}

func CheckoutGitBranch(b string) (error, string) {
	cmd := exec.Command("git", "checkout", b)
	if stdOutErr, err := cmd.CombinedOutput(); err != nil {
		return err, string(stdOutErr)
	}
	return nil, ""
}

func RenameGitBranch(oldName string, newName string) (error, string) {
	cmd := exec.Command("git", "branch", "-m", oldName, newName)
	if stdOutErr, err := cmd.CombinedOutput(); err != nil {
		return err, string(stdOutErr)
	}
	return nil, ""
}

func CreateGitBranch(newBranch string, base string) (error, string) {
	cmd := exec.Command("git", "checkout", "-b", newBranch, base)
	if stdOutErr, err := cmd.CombinedOutput(); err != nil {
		return err, string(stdOutErr)
	}
	return nil, ""
}
