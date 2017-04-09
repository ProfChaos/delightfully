package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/profchaos/printy"
)

var (
	match          string
	name           string
	season         string
	nameIndex      int
	episodeIndex   int
	seasonIndex    int
	extensionIndex int
	replacable     [][]string
)

func main() {
	replacable = [][]string{}
	flag.StringVar(&match, "match", `^(.+)\.[sS]([0-9]+)[eE]([0-9]+).*\.(.+)$`, "Regex string")
	flag.StringVar(&name, "name", "", "New name of show")
	flag.StringVar(&season, "season", "", "New season of show")
	flag.IntVar(&nameIndex, "na", 1, "Index of the name in matches")
	flag.IntVar(&seasonIndex, "se", 2, "Index of the season in matches")
	flag.IntVar(&episodeIndex, "ep", 3, "Index of the episode in matches")
	flag.IntVar(&extensionIndex, "ext", 4, "Index of the extension in matches")
	flag.Parse()

	fs, err := ioutil.ReadDir(".")
	if err != nil {
		printy.Err(err)
		return
	}

	printy.Info("Regex: " + match)
	r, err := regexp.Compile(match)
	if err != nil {
		printy.Err(err)
		return
	}

	printy.Info("`Old name` => `New name`")
	renamable := 0
	for _, o := range fs {
		m := r.FindAllStringSubmatch(o.Name(), -1)
		if len(m) > 0 {
			n := m[0]
			if len(n) <= 4 {
				printy.Err("Too few matches in regex. Min 4")
				return
			}
			if season == "" {
				season = n[seasonIndex]
			}

			if name == "" {
				name = n[nameIndex]
			}
			newName := name + "." + "S" + season + "E" + n[episodeIndex] + "." + n[extensionIndex]
			printy.Log("`" + o.Name() + "` => `" + newName + "`")
			matches := []string{o.Name(), newName}
			replacable = append(replacable, matches)
			renamable++
		}
	}

	if renamable == 0 {
		printy.Info("No files to rename")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	printy.Info("Are you sure you want to rename " + strconv.Itoa(renamable) + " file(s)? (Y/n): ")
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	if !(text == "y" || text == "Y" || text == "") {
		printy.Info("No files changed")
		return
	}

	renamedNumber := 0
	for _, arr := range replacable {
		if err = os.Rename(arr[0], arr[1]); err != nil {
			printy.Err("Couldn't rename file: " + arr[0])
			continue
		}
		renamedNumber++
	}

	printy.Info(strconv.Itoa(renamedNumber) + " file(s) renamed")
}
