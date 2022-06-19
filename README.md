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
$./grpcUploadServer
Eg ./grpcUploadServer
Client : Upload all files in the specified directory to the server :
```

```shell
$ ./UploadClient upload  -d <folder containing files to upload>   
Eg  ./UploadClient upload -a localhost:9191 -d /home/
```

```shell
<rest addr>/list 
to list all file in server
```


* for all files in folder

image-server upload ./source/image-1.jpg ./source/image-2.jpg

check if file exist.- > check if file is png or jpg