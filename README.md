# imageCache



## REFERENCES   
###### https://docs.buf.build
###### A modification of the code from https://github.com/rickslick/grpcUpload



## FROM https://github.com/rickslick/grpcUpload


grpcUpload is CLI tool that uploads files concurrently using grpc. Featues :

concurrent multi file upload using grpc with concept of chunking
supports tls (both client and sever )
Displays progress for each file


Usage
Server : start the server( default destination of files is /data/files) :


## Server
### 1.  Build Client as executable
```shell
go build -o server cmd/server/main.go
```

### 1.b Build as a docker container
```shell
docker build . -t imageCache
```

### 2. Configure .env file 
```shell
example: 

GRPC_ADDR="127.0.0.1:4000"
REST_ADDR=":9900"

ps: this is be in the same dir as go.mod (i.e the root directory of the project)
```
### 3. Run Server
```shell
./server

output: 

[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

=========> We are running GRPC @ 127.0.0.1:4000
[GIN-debug] GET    /data/files/:fileName     --> imageCache/delivery/server.getSingleFileHandler (3 handlers)
[GIN-debug] GET    /list                     --> imageCache/delivery/server.listFilesHandler (3 handlers)
======> we are running REST @ :9900

``` 
### 4. List uploaded files
```shell
<server-ip:port>/list

e.g:
  localhost:9900/list

response:
  ["data/files/abc.png",
  "data/files/main.jpg",
  "data/files/sample.jpeg"]
```
### 5. Access a file
```shell
<server-ip:port>/data/files/<fileName>

e.g:
  localhost:9900/data/files/flower.jpg
```
---
## Client
### 1.  Build Client
```shell
go build -o client cmd/cli/main.go
```
### 2. Upload a file
```shell
./client upload -a <server address> -d <folder to send>


e.g:
  ./client upload -a 127.0.0.1:4000 -d ./images/
```

