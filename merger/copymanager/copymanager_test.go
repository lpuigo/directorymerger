package copymanager

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"
	"time"
)

const (
	testFile   = `C:\Users\Laurent\Golang\src\github.com\lpuig\directorymerger\merger\copymanager\test\testfile.txt`
	//testFile   = `R:\Test\testfile.txt`
	targetFile = `C:\Users\Laurent\Golang\src\github.com\lpuig\directorymerger\merger\copymanager\test\targetfile.txt`
	//targetFile = `D:\Test\targetfile.txt`
	//targetFile = `V:\Series\Test\targetfile.txt`
	Size       = 1000000
)

func createTestFile(name string, size int) error {
	fd, err := os.Create(name)
	if err != nil {
		return err
	}
	defer fd.Close()

	lorem := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed dapibus lacinia urna vel molestie. Viamus sit amet nibh pharetra, scelerisque libero at, vehicula purus. Maecenas tempor egestas metus, sit amet ullamcorper odio faucibus nec. Cras eget libero feugiat est molestie scelerisque. Integer eget quam sit amet nunc mollis tristique. Nulla egestas finibus nulla ut porta.
Vestibulum sit amet quam sed orci condimentum eleifend ut eu lectus. Pelentesque convallis, metus in laoreet maximus, lacus erat vehicula mi, a rhoncus sapien nulla at odio. Aliquam lacus massa, tincidunt non risus ut, efficitur laoreet quam. Proin maximus, ipsum et lacinia pharetra, sapien enim bibendum sapien, vitae sollicitudin magna felis at dui. Quisque at dui nibh. Sed maximus risus massa, non tincidunt elit tincidunt in. Ut molestie eu odio ut bibendum. Donec tempor vitae justo aliquam gravida. Quisque sit amet velit libero. Cras gravida pellentesque eleifend. Sed elementum accumsan orci. Nullam mi sem, imperdiet et aliquam id amet.
`

	for i := 0; i < size; i++ {
		_, err := fd.WriteString(lorem)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestCreateTestFile(t *testing.T) {
	if err := createTestFile(testFile, Size); err != nil {
		log.Fatal("could not create test file:", err.Error())
	}
}

func TestNewMonitoredReader(t *testing.T) {
	sfd, err := os.Open(testFile)
	if err != nil {
		t.Fatal("could not open:", err.Error())
	}
	defer sfd.Close()

	tfd, err := os.Create(targetFile)
	if err != nil {
		t.Fatal("could not create:", err.Error())
	}
	defer tfd.Close()

	mfd, monitor := NewMonitoredReader(sfd)

	monitor.GetInfo()

	ticks := time.NewTicker(time.Millisecond * 1000/10)
	done := make(chan int)
	go func() {
		for {
			select {
			case <-ticks.C:
				size, dur, atp, utp, inprogress := monitor.GetInfo()
				fmt.Printf("%5.2fs : %10.3f MB (instant: %8.3f MB/s, avg: %8.3f MB/s)\n", dur, float64(size/1024)/1024, utp, atp)
				if !inprogress {
					ticks.Stop()
					fmt.Println("Ticker Done")
					done <- 0
					return
				}
			}
		}
	}()

	_, err = io.Copy(tfd, mfd)
	if err != nil {
		t.Fatal("could not create:", err.Error())
	}
	<-done
}
