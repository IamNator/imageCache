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
Server : start the server( default destination of files is /tmp) :

```shell
$./grpcUploadServer serve --a <ip:port> -d <destination folder>
Eg ./UploadClient serve -a localhost:9191 -d /media/
Client : Upload all files in the specified directory to the server :
```

```shell
$ ./UploadClient upload  -a <ip:port> -d <folder containing files to upload>   
Eg  ./UploadClient upload -a localhost:9191 -d /home/
```