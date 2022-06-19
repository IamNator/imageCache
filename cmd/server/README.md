## SERVER


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
