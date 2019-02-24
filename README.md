# Screenshooter
Coding challenge for Detectify

Build and start container
```
docker build -t screenshooter -f build/Dockerfile .
docker run -d --rm -p 8080:8080 -v ~/uploads:/root/uploads screenshooter
```

Credits:
Inspiration for screenshots: https://github.com/sensepost/gowitness
Installing chrome on alpine: https://github.com/Zenika/alpine-chrome/blob/master/Dockerfile
My fav router: https://github.com/go-chi/chi