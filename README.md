# Screenshooter
Coding challenge for Detectify

## Build and start container
```
docker build -t screenshooter -f build/Dockerfile .
docker run -d --rm -p 8080:8080 -v ~/uploads:/root/uploads screenshooter
```

## API
POST /screenshot
- Request body should be a json object with the "urls" parameter set to a list of urls you want to screenshot.
- Example: curl -X POST -d '{"urls": ["https://google.com"]}' localhost:8080/screenshot
- Response will contain a list of urls, one for each screenshot

GET /screenshot/{filename}
- {filename} should be one of the urls returned by the POST endpoint
- Example: curl -X GET localhost:8080/screenshow/https-google.com.png

## Credits:
Inspiration for screenshots: https://github.com/sensepost/gowitness
Installing chrome on alpine: https://github.com/Zenika/alpine-chrome/blob/master/Dockerfile
My fav router: https://github.com/go-chi/chi