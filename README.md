# Screenshooter
Coding challenge for Detectify

## Build
```
$ docker build -t screenshooter -f build/Dockerfile .
$ docker run -d --rm -p 8080:8080 -v ~/uploads:/root/uploads screenshooter
```

... or if you want to run it without docker:
```
$ git clone https://github.com/traed/screenshooter $GOPATH/src/github.com/traed/screenshooter
$ cd $GOPATH/src/github.com/traed/screenshooter
$ go get -d -v ./...
$ go build -o bin/scsh main.go
$ bin/scsh
```

## API
POST /screenshot
- Request body should be a json object with the "urls" parameter set to a list of urls you want to screenshoot.
- Example: curl -X POST -d '{"urls": ["https://google.com"]}' localhost:8080/screenshot
- Response will contain a list of urls, one for each screenshot

GET /screenshot/{filename}
- {filename} should be one of the urls returned by the POST endpoint
- Example: curl -X GET localhost:8080/screenshot/https-google.com.png
- Response is the image as a stream

## Scaling
This project implements a Worker/Job queue in order to handle large amounts of requests. This ensures that the main thread isn't bogged down when downloading screenshots by dispatching that work to a different process. It also makes sure that the server won't crashes if too many request comes in at the same time by only allowing a certain number of workers to run at one time. If multiple instances of the project is run behind a load balancer the max number of jobs in the queue as well as the number of concurrant workers could be read via environment variables but for now they are hard coded in the main.go file.

## Credits:
- Inspiration for screenshots: https://github.com/sensepost/gowitness
- Inspiration for worker package: https://github.com/cahitbeyaz/job-worker
- Installing chrome on alpine: https://github.com/Zenika/alpine-chrome/blob/master/Dockerfile
- My fav router: https://github.com/go-chi/chi
