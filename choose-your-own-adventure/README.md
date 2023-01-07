# choose-your-own-adventure

## Usage

Using `go run cmd/web/main.go -h`, we get all the flags:

- cli: Whether to run the cli version or not
- port: Port on which to run the server (default "8082")
- story: Path of the story file (default "stories.json")
- template: Path of the template file (default "templates/index.html")

## Example

1. `go run cmd/web/main.go -cli -story="stories.json"` to use the _CLI_ version where the story file is stored in _stories.json_.

2. `go run cmd/web/main.go -port=8082 -story="stories.json" -template="templates/index.html"` to use the _default GUI_ version running the server on _port 8082_ where the template file is stored in _templates/index.html_.
