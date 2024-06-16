# Development 

## Build
Use the following commands to build the software for the different operating systems
````
GOOS=linux GOARCH=amd64 go build -o tedprocessor   
GOOS=darwin GOARCH=arm64 go build -o tedprocessor-darwin-arm64 
GOOS=windows GOARCH=amd64 go build -o tedprocessor-windows-amd64.exe 
````

