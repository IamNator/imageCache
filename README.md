# imageCache



## REFERENCES   
###### https://docs.buf.build
###### A modification of the code from https://github.com/rickslick/grpcUpload

```markdown
imageCache can be used a CLI tool and a Server.  
The CLI tool uploads files concurrently using grpc to a server instance of the application.

If the server is spurned in a container  
It is Ephemeral, meaning, all data is erased once the container shuts down.  
To prevent this, you can make use of volumes. 
``` 

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
<i>CLI communicates with the grpc server from a terminal</i>
<br clear="left"/>
<br>

---

## Quick Start

#### 1. spin off server
```shell
go run cmd/server/main.go
```
[more ...](./cmd/server/README.md)

#### 2. upload a file
```shell
go run cmd/client/main.go upload 127.0.0.1:4000 -f ./testdata/sample23.png
```
[more ...](./cmd/cli/README.md)

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
docker run --env-file .env -p 9900:<GRPC_PORT> -p 7700:<REST_PORT>  <image Tag Name>

e.g  docker run --env-file .env -p 9900:9900 -p 7700:7700 server  

output:
mac@macs-MBP imageCache %  docker run --env-file .env -p 9900:9900 -p 7700:7700 server
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /data/files/:fileName     --> imageCache/delivery/server.DownloadFileHandler (3 handlers)
[GIN-debug] GET    /list                     --> imageCache/delivery/server.ListFilesHandler (3 handlers)
======> we are running REST APIs @ :9900
=========> We are running GRPC APIs @ 127.0.0.1:7700


```
