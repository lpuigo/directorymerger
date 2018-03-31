package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"github.com/lpuig/directorymerger/merger/compare"
)

const (
	directory     = `C:/Users/Laurent/Downloads/JDownloader`
	remoteDir     = `V:/Series`
	MinMatchScore = 0.85
)

func main() {
	err := mergeDir(directory)
	if err != nil {
		log.Println(err)
	}
}

func reverse(a []os.FileInfo) {
	for i := len(a)/2-1; i >= 0; i-- {
		opp := len(a)-1-i
		a[i], a[opp] = a[opp], a[i]
	}
}

func mergeDir(dirname string) error {
	fds, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fmt.Errorf("could not read directory :%s", err.Error())
	}
	reverse(fds)
	previous := ""
	for i, fd := range fds {
		if fd.IsDir() {
			name := fd.Name()
			fmt.Printf("%02d: %.3f %s\n", i, compare.MatchScore(name, previous), name)
			if compare.MatchScore(name, previous) > MinMatchScore {
				matchAction(name, previous)
			} else {
				notMAtchAction()
				previous = name
			}
		}
	}
	return nil
}

func compareDir(sourceDir, destDir string) ([]compare.CompList, error) {
	sfds, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("could not read directory :%s", err.Error())
	}
	dfds, err := ioutil.ReadDir(destDir)
	if err != nil {
		return nil, fmt.Errorf("could not read directory :%s", err.Error())
	}

	getFileName := func(fis []os.FileInfo) []string {
		res := []string{}
		for _, fd := range fis {
			if fd.IsDir() {
				res = append(res, fd.Name())
			}
		}
		return res
	}

	sourceDirNames := getFileName(sfds)
	destDirNames := getFileName(dfds)

	cls := make([]compare.CompList, len(sourceDirNames))
	for i, sourceDirName := range sourceDirNames {
		cls[i] = compare.NewCompList(sourceDirName, destDirNames)
	}
	return cls, nil
}

func matchAction(new, old string) {
	fmt.Printf("\t move '%s' content in '%s'\n", new, old)
	targetdir := path.Join(directory, old)
	sourcedir := path.Join(directory, new)

	fds, err := ioutil.ReadDir(sourcedir)
	if err != nil {
		log.Println("could not readDir:", err)
		return
	}
	for _, fd := range fds {
		name := fd.Name()
		err := os.Rename(path.Join(sourcedir, name), path.Join(targetdir, name))
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println("ok")
		}
	}
	fmt.Printf("\tDelete '%s' ...", sourcedir)
	err = os.Remove(sourcedir)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ok")
	}
}

func notMAtchAction() {
	fmt.Println("\tSkip")
}
