package git

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
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

func getGitDir() (error, string) {
	const GIT_DIR = "/.git/"
	err, dir := getRepoDir()
	if err != nil {
		return err, ""
	}
	return nil, dir + GIT_DIR
}

func getRepoDir() (error, string) {
	const GIT_DIR = "/.git/"
	wd, _ := os.Getwd()
	// Traverse all the way to root '/'
	for {
		exist, _ := dirExists(wd + GIT_DIR)
		if exist {
			break
		} else {
			if wd == "/" {
				return errors.New("no git repo found"), ""
			}
			wd = filepath.Dir(wd)
		}
	}
	return nil, wd
}

// Returns all git branch names in repo
func GetGitBranches() (error, []string) {
	err, dir := getRepoDir()
	if err != nil {
		return err, nil
	}
	repo, err := git.PlainOpen(dir)
	if err != nil {
		return err, nil
	}
	refIter, err := repo.Branches()
	if err != nil {
		return err, nil
	}
	var branches []string
	if err := refIter.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, strings.Replace(ref.Name().String(), "refs/heads/", "", 1))
		return nil
	}); err != nil {
		return err, nil
	} else {
		return nil, branches
	}
}

// hierarchical branch name e.g.: feat/create-new-branch
func getHierarchicalBranches(path string, prefix string) (error, []string) {
	files, err1 := ioutil.ReadDir(path)
	if err1 != nil {
		return err1, nil
	}
	fnames := make([]string, 0)
	pre := ""
	if prefix != "" {
		pre = prefix + "/"
	}
	for _, f := range files {
		if n := f.Name(); f.IsDir() {
			err2, h := getHierarchicalBranches(path+"/"+n, pre+n)
			if err2 != nil {
				return err2, nil
			}
			fnames = append(fnames, h...)
		} else {
			fnames = append(fnames, pre+n)
		}
	}
	return nil, fnames
}

func GetCurrentBranch() (error, string) {
	err, gitDir := getGitDir()
	if err != nil {
		return err, ""
	}
	dat, _ := os.ReadFile(gitDir + "HEAD")
	raw := string(dat)
	current := strings.TrimSpace(strings.Replace(raw, "ref: refs/heads/", "", 1))
	return nil, current
}

func DeleteGitBranch(b string) (error, string) {
	_, dir := getRepoDir()
	repo, _ := git.PlainOpen(dir)
	ref := plumbing.NewBranchReferenceName(b)
	err := repo.Storer.RemoveReference(ref)
	if err != nil {
		return err, err.Error()
	}
	return nil, ""
}

func CheckoutGitBranch(b string) (error, string) {
	_, dir := getRepoDir()
	repo, _ := git.PlainOpen(dir)
	worktree, _ := repo.Worktree()
	name := plumbing.ReferenceName("refs/heads/" + b)
	opts := &git.CheckoutOptions{
		Branch: name,
		Create: false,
		Keep:   true,
	}
	err := worktree.Checkout(opts)
	if err != nil {
		return err, err.Error()
	}
	return nil, ""
}

func RenameGitBranch(oldName string, newName string, currentBranch string) (error, string) {
	// Reference: https://github.com/go-git/go-git/issues/233
	_, dir := getRepoDir()
	repo, _ := git.PlainOpen(dir)
	name := plumbing.ReferenceName("refs/heads/" + oldName)
	baseRef, _ := repo.Reference(name, true)

	// Create the new branch
	branchRef := plumbing.NewBranchReferenceName(newName)
	ref := plumbing.NewHashReference(branchRef, baseRef.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		return err, err.Error()
	}

	// Checkout to the new branch
	if currentBranch == oldName {
		worktree, err := repo.Worktree()
		if err != nil {
			return err, err.Error()
		}
		opts := &git.CheckoutOptions{Branch: branchRef, Keep: true}
		if err := worktree.Checkout(opts); err != nil {
			return err, err.Error()
		}
	}

	// Remove the old branch
	if err := repo.Storer.RemoveReference(baseRef.Name()); err != nil {
		return err, err.Error()
	}

	return nil, ""
}

func CreateGitBranch(newBranch string, base string) (error, string) {
	_, dir := getRepoDir()
	repo, _ := git.PlainOpen(dir)
	name := plumbing.ReferenceName("refs/heads/" + base)
	baseRef, _ := repo.Reference(name, true)
	newBranchName := plumbing.ReferenceName("refs/heads/" + newBranch)
	ref := plumbing.NewHashReference(newBranchName, baseRef.Hash())
	if err := repo.Storer.SetReference(ref); err != nil {
		return err, err.Error()
	}
	return nil, ""
}
