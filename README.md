# grpc-api

## Installation

`git clone https://github.com/v1tbrah/grpc-api`

## Getting started server

### Starting

* Open terminal. Go to project working directory. Run grpc server. For example:
 ```
 cd ~/go/src/grpc-api
 go run cmd/server/main.go
 ```

### Options
The following server options are set by default:
```
api server run address: `:8080`
directory with files: `filesSavedInGRPCServer`
log level: `info`
```
* flag options:
```
   -a string
      api server run address
   -d string
      directory with files
   -l string
      log level 
```
* env options you can check in internal/server/config/parse

## Getting started test client

### Save file

* Open terminal. Go to project working directory. Run save file procedure. For example:
 ```
 cd ~/go/src/grpc-api
 go run cmd/userSaveFile/userSaveFile.go
 ```

#### Options
The following options are set by default:
```
grpc api server address: `:8080`
dir with files for sending to server: `filesForSendingToGRPCServer`
```
* flag options:
```
   -a string
      api server run address
   -d string
      dir with files for sending to server
```

### Get files info

* Open terminal. Go to project working directory. Run get files info procedure. For example:
 ```
 cd ~/go/src/grpc-api
 go run cmd/userTestGetFilesInfo/userTestGetFilesInfo.go
 ```

#### Options
The following options are set by default:
```
grpc api server address: `:8080`
count goroutines: 103
```
* flag options:
```
   -a string
      api server run address
   -c int
      count goroutines
```

### Get files

* Open terminal. Go to project working directory. Run get files procedure. For example:
 ```
 cd ~/go/src/grpc-api
 go run cmd/userTestGetFiles/userTestGetFiles.go
 ```

#### Options
The following options are set by default:
```
grpc api server address: `:8080`
count goroutines: 13
```
* flag options:
```
   -a string
      api server run address
   -c int
      count goroutines
```

