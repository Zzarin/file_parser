FROM golang:latest
WORKDIR /app
ENV f="./file.txt"
ENV t=2s

#Before we can run go mod download inside our image, we need to get our go.mod and go.sum files copied into it.
#We’ll copy the go.mod and go.sum file into our project directory /app which, owing to our use of WORKDIR, is the current directory (.) inside the image.
COPY go.mod ./
COPY go.sum ./

#Now that we have the module files inside the Docker image that we are building, we can use the RUN command to execute the command go mod download
#Go modules will be installed into a directory inside the image.
RUN go mod download

#COPY command takes two parameters. The first parameter tells Docker what files you want to copy into the image. The last parameter tells Docker where you want that file to be copied to.
# We’ll copy the go.mod and go.sum file into our project directory /app which, owing to our use of WORKDIR, is the current directory (.) inside the image.
COPY ./ ./

RUN go build -o parser ./cmd/cli

EXPOSE 4067

CMD ["/app/parser"]