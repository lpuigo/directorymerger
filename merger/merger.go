package main

import (
	"fmt"
	"github.com/xrash/smetrics"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
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

func mergeDir(dirname string) error {
	fds, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fmt.Errorf("could not read directory :%s", err.Error())
	}

	previous := ""
	for i, fd := range fds {
		if fd.IsDir() {
			name := fd.Name()
			fmt.Printf("%02d: %.3f %s\n", i, match(name, previous), name)
			if match(name, previous) > MinMatchScore {
				matchAction(name, previous)
			} else {
				notMAtchAction()
				previous = name
			}
		}
	}
	return nil
}

func compareDir(sourceDir, destDir string) ([]CompList, error) {
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

	cls := make([]CompList, len(sourceDirNames))
	for i, sourceDirName := range sourceDirNames {
		cls[i] = NewCompList(sourceDirName, destDirNames)
	}
	return cls, nil
}

type CompList struct {
	name  string
	pairs []string
	dists []float64
}

func (cl CompList) Len() int {
	return len(cl.pairs)
}

func (cl CompList) Swap(i, j int) {
	cl.pairs[i], cl.pairs[j] = cl.pairs[j], cl.pairs[i]
	cl.dists[i], cl.dists[j] = cl.dists[j], cl.dists[i]
}

func (cl CompList) Less(i, j int) bool {
	return cl.dists[i] > cl.dists[j]
}

func (cl CompList) String() string {
	return cl.StringNearest(0)
}

func (cl CompList) StringNearest(mindist float64) string {
	res := fmt.Sprintf("Comparison with '%s':\n", cl.name)
	for i, pair := range cl.pairs {
		if cl.dists[i] >= mindist {
			res += fmt.Sprintf("\t%0.3f : %s\n", cl.dists[i], pair)
		}
	}
	return res
}

func NewCompList(name string, pairs []string) CompList {
	cl := CompList{
		name:  name,
		pairs: append([]string{}, pairs...),
		dists: make([]float64, len(pairs)),
	}
	for i, pair := range pairs {
		cl.dists[i] = match(name, pair)
	}
	sort.Sort(cl)
	return cl
}

func match(new, old string) (ratio float64) {
	neww := strings.Split(strings.ToUpper(new), " ")
	oldw := strings.Split(strings.ToUpper(old), " ")

	ratio = 0
	nwc := 0
	for i, w := range neww {
		if i >= len(oldw) {
			break
		}
		sefound, err := regexp.MatchString("^S[0-9]{2}E[0-9]{2}", w)
		if err != nil {
			panic(err)
		}
		if sefound {
			//TODO check if season number match
			break
		}
		if w == oldw[i] {
			ratio += 1
		} else {
			ratio += smetrics.JaroWinkler(w, oldw[i], 0.7, 4)
		}
		nwc++
	}

	if nwc > 0 {
		ratio /= float64(nwc)
	}

	return
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
