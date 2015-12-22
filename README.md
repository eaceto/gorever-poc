# gorever (Proof of Concept)

## Dependencies

    go get gopkg.in/alecthomas/kingpin.v2 github.com/Sirupsen/logrus github.com/inconshreveable/go-update

## Steps to run

### Generate executable file
   
    go build
   
### Start an http server at port 8078 which will serve the executable file

    http-server -p 8078
    
### Execute gorever
    
    ./gorever-poc
    
    
