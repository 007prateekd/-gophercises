# choose-your-own-adventure

## Usage

Using `go run cmd/web/main.go -h`, we get all the flags:

- arc: The arc from which to start the story (default "intro")
- cli: Whether to run the cli version or not
- port: Port on which to run the server (default "8082")
- story: Path of the story file (default "stories.json")
- template: Path of the template file (default "templates/index.html")

## Example

1. `go run cmd/web/main.go -cli -story="stories.json"` to use the _CLI_ version where the story file is stored in _stories.json_.

2. `go run cmd/web/main.go -port=8082 -story="stories.json" -template="templates/index.html -arc="new-york"` to use the _default GUI_ version running the server on _port 8082_ where the template file is stored in _templates/index.html_ and the starting arc is _new-york_.

## Note

Regarding the 2<sup>nd</sup> bonus of this excercise, in order to support stories starting from a story-defined arc, we can simply add a flag to input the arc name. To change the web (GUI) version of the code, we can add another functional option for the same as can be seen in [this](cmd/util/main.go) file. Change in the CLI code is minimal and easily understable by looking at the code.
