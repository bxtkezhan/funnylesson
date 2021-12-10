CC=x86_64-w64-mingw32-gcc
CXX=x86_64-w64-mingw32-g++

run:
	go run main.go --port 8888

build:
	go build -o a.out main.go

build_win:
	GOOS=windows \
	GOARCH=amd64 \
	CGO_ENABLED=1 \
	CC=${CC} \
	CXX=${CXX} \
	go build -o a.out.exe
