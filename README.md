# imageCache



## REFERENCES   
###### https://docs.buf.build
###### A modification of the code from https://github.com/rickslick/grpcUpload

imageCache can be used a CLI tool and a Server.  
The CLI tool uploads files concurrently using grpc to a server instance of the application.

If the server is spurn in a container  
It is Ephemeral, meaning, all data once the container shuts down.  
To prevent this, you can make use of volumes.  

###### by default: all uploaded images are stored in <b>/app/data/files</b> 

---
###### Featues :
* concurrent multi file upload using grpc with concept of chunking
* supports tls (both client and sever )
* Displays progress for each file
* Exposes endpoints for file downloads
* Exposes endpoint that return list of files of uploaded
---
###### How it works

<img width="285" align="left" style="margin-right: 14px; margin-top: 7px;"  alt="Screenshot 2022-06-19 at 13 20 11" src="https://user-images.githubusercontent.com/43158886/174480621-7c487cf7-8eac-46e5-a945-79fc79eb966b.png">
<br><br><br><br><br>
<i>Server Exposes GRPC and REST apis</i>
<br><br><br><br><br><br><br><br><br><br><br><br>
<br><b></b>
<i>CLI communicates with the grpc server from a terminal</i>
<br clear="left"/>
<br>

---

## Usage

### Server
#### 1.  Build Server
```shell
go build -o server cmd/server/main.go
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

Server : destination of files uploaded is /data/files :

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
---



### Working with Container (Docker)  

#### 1.b Build Docker container
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

#### 3.b Run Container (locally)
```shell
docker run --env-file .env imageCache

```
