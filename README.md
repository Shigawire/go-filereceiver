# go-filereceiver
Receives a file via HTTP and writes to disk.

# Usage
In development, just do `docker-compose up --build`. The image will be built until the `build-env` stage of the Dockerfile and then run `go run main`.

In production, use the registry image:  

`docker run --rm -e ... shigawire/go-filereceiver`.
