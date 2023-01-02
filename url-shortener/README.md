# url-shortener

## Usage

Using `./url-shortener -h`, we get all the flags:

- json: Path of JSON file
- port: Port on which the server runs (default "8081")
- yaml: Path of YAML file

For example, we can use  
`./url-shortener -port=8082 -json="mappings/mapping.json"`
to run the server on _port 8082_ where the url mapping is stored in _mappings/mapping.json_