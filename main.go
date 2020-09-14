// gowatch is a small command-line tool that makes it easy
// to automatically execute/retest your go code whenever a change
// is made
package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/fatih/color"
	"github.com/radovskyb/watcher"
)

func main() {
	args := os.Args[1:]
	os.Exit(gowatch(args))
}

func gowatch(args []string) int {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Rename, watcher.Create, watcher.Write, watcher.Remove)

	r := regexp.MustCompile(".*\\.go")
	w.AddFilterHook(watcher.RegexFilterHook(r, false))
	// execute the first time before going into a loop
	execute(args)

	go func() {
		for {
			select {
			case _ = <-w.Event:
				execute(args)
			case err := <-w.Error:
				color.Red("Error: ", err.Error())
			}

		}
	}()

	if err := w.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}

	// for path, f := range w.WatchedFiles() {
	// 	fmt.Printf("%s: %s\n", path, f.Name())
	// }

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}

	return 0
}

func execute(args []string) {
	args = append([]string{"test"}, args...)
	out, err := exec.Command("go", args...).Output()
	if err != nil {
		color.Red(string(out))
		return
	}
	color.Green(string(out))
}
