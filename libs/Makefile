export GO111MODULE=on
export GOPROXY=https://goproxy.cn

default:
	#go build -buildmode=c-shared -o $(shell go list -m).so
	go build -buildmode=c-archive

clean:
	rm -rvf *.so *.a
