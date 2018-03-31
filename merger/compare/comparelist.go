package compare

import (
	"fmt"
	"sort"
	"strings"
	"regexp"
	"github.com/xrash/smetrics"
)

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
		cl.dists[i] = MatchScore(name, pair)
	}
	sort.Sort(cl)
	return cl
}

func MatchScore(new, old string) (ratio float64) {
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


