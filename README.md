# imageCache



## REFERENCES   
###### https://docs.buf.build
###### A modification of the code from https://github.com/rickslick/grpcUpload

imageCache is CLI tool that uploads files concurrently using grpc. 

---
###### Featues :
* concurrent multi file upload using grpc with concept of chunking
* supports tls (both client and sever )
* Displays progress for each file
---
###### How it works

<img width="285" alt="Screenshot 2022-06-19 at 13 20 11" src="https://user-images.githubusercontent.com/43158886/174480621-7c487cf7-8eac-46e5-a945-79fc79eb966b.png">


The cli tools communicates with the server via GRPC, the server also exposes rest api's.

---

## Usage
Server : start the server( default destination of files is /data/files) :


### Server
#### 1.  Build Client as executable
```shell
go build -o server cmd/server/main.go
```

#### 1.b Build as a docker container
```shell
docker build . -t imageCache

output:
mac@macs-MBP imageCache % docker build . -t imageCache          
[+] Building 116.0s (8/11)                                                                                                                                                  
 => [internal] load build context                                                                                                                                      2.0s
 => => transferring context: 51.91kB                                                                                                                                   1.6s
 => [builder 1/4] FROM docker.io/library/golang:1.18-alpine@sha256
 ...
```

#### 2. Configure .env file 
```shell
example: 

GRPC_ADDR=127.0.0.1:4000
REST_ADDR=:9900

ps: this is be in the same dir as go.mod (i.e the root directory of the project)
```
#### 3. Run Server
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

#### 3.b RUN Server Container (locally)
```shell
docker run --env-file .env imageCache

```
#### 4. List uploaded files
```shell
<server-ip:port>/list

e.g:
  localhost:9900/list

response:
  ["data/files/abc.png",
  "data/files/main.jpg",
  "data/files/sample.jpeg"]
```
#### 5. Access a file
```shell
<server-ip:port>/data/files/<fileName>

e.g:
  localhost:9900/data/files/flower.jpg
```
---
## Client
#### 1.  Build Client
```shell
go build -o client cmd/cli/main.go
```
#### 2. Upload all files in a folder
```shell
./client upload -a <server address> -d <folder to upload>


e.g:
  ./client upload -a 127.0.0.1:4000 -d ./images/
```

#### 3. Upload a file
```shell
./client upload -a <server address> -f <file to upload>


e.g:
  ./client upload -a 127.0.0.1:4000 -f sample.png
```

#### 4. Upload multiple files
```shell
./client upload -a <server address> -f <file1> -f <file2>


e.g:
  ./client upload -a 127.0.0.1:4000 -f sample.png -f sample.jpeg
```
