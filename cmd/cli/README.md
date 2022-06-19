## Client

#### 1.  Build Client
```shell
go build -o client cmd/cli/main.go
```
#### 2. Upload a file
```shell
./client upload -a <server address> -f <file to upload>


e.g:
  ./client upload -a 127.0.0.1:4000 -f sample.png
```

#### 3. Upload multiple files
```shell
./client upload -a <server address> -f <file1> -f <file2>


e.g:
  ./client upload -a 127.0.0.1:4000 -f sample.png -f sample.jpeg
```

#### 4. Upload all files in a folder
```shell
./client upload -a <server address> -d <folder to upload>


e.g:
  ./client upload -a 127.0.0.1:4000 -d ./images/
```


