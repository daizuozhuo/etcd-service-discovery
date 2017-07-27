# etcd-service-discovery
daizuozhuo.github.io/etcd-service-discovery/

### Build
1. Install Golang: sudo yum install golang
2. mkdir $HOME/Go
3. add environment variable to ~/.bashrc `export GOPATH=$HOME/Go` `export PATH=$HOME/Go/bin:$PATH`,
reload it `source ~/.bashrc`
4. put the project inside $HOME/Go/src/github.com/daizuozhuo/
5. `go get github.com/tools/godep`
6. `godep restore`
7. `cd etcd-service-discovery/exmaple"
8. `go build`

### Run
1. run etcd server on the localhost.
2. `./example -role master`
3. `./example -role worker`
