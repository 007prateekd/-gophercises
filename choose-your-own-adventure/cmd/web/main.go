package main

import (
	"cyoa/cmd/util"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	storyFile    string
	templateFile string
	port         string
	useCli       bool
	arc          string
)

func init() {
	flag.StringVar(&storyFile, "story", "stories.json", "Path of the story file")
	flag.StringVar(&templateFile, "template", "templates/index.html", "Path of the template file")
	flag.StringVar(&port, "port", "8082", "Port on which to run the server")
	flag.BoolVar(&useCli, "cli", false, "Whether to run the cli version or not")
	flag.StringVar(&arc, "arc", "intro", "The arc from which to start the story")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	fileBytes, err := ioutil.ReadFile(storyFile)
	checkError(err)
	story, err := util.ParseJson(fileBytes)
	checkError(err)

	if useCli {
		chapter := arc
		chapterId := 1
		for {
			fmt.Printf("\nCHAPTER #%d: %s\n\n", chapterId, story[chapter].Title)
			for _, paragraph := range story[chapter].Story {
				fmt.Println(paragraph)
			}
			if len(story[chapter].Options) == 0 {
				fmt.Println("\nTHE END")
				break
			}
			fmt.Println("\nOPTIONS:")
			for i, option := range story[chapter].Options {
				fmt.Printf("Press %d: %s\n", i+1, option.Text)
			}
			var id int
			fmt.Scan(&id)
			chapter = story[chapter].Options[id-1].Chapter
			chapterId += 1
		}
	} else {
		util.InitializeTemplate(templateFile)
		h := util.NewHandler(
			story,
			util.Template(util.DefaultTemplate),
			util.PathFn(util.DefaultPathFn),
			util.Arc(arc),
		)
		port = ":" + port
		fmt.Printf("Starting server on %s\n", port)
		log.Fatal(http.ListenAndServe(port, h))
	}
}
