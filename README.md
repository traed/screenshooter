# Screenshooter
Coding challenge for Detectify

Build and start container
```
docker build -t screenshooter -f build/Dockerfile .
docker run -d --rm -p 8080:8080 -v ~/uploads:/root/uploads screenshooter
```