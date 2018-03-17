package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"github.com/xrash/smetrics"
	"strings"
	"regexp"
	"path"
	"os"
)

const (
	directory     = `C:\Users\Laurent\Downloads\JDownloader`
	MinMatchScore = 0.85
)

func main() {
	fds, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal("could not read directory :", err)
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