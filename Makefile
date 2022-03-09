name=ormat

execveFile:=${name} # 设置固件名称

# 路径相关
ProjectDir=.

# 编译平台
platform = CGO_ENABLED=0
# 编译选项,如tags,多个采用','分开 sqlite3,noswag
opts = -trimpath -tags=sqlite3
# 编译flags
path = github.com/things-go/ormat/pkg/builder
flags = -ldflags "-w -s" # -s 引起gops无法识别go版本号,upx压缩也同样

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${opts} ${flags} -o ${execveFile} main.go
windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${opts} ${flags} -o ${execveFile}.exe main.go
mac:
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build ${opts} ${flags} -o ${execveFile} main.go

clear:
	test ! -d model/ || rm -rf  model/*
	test ! -f ormat || rm ormat
	test ! -f ormat.exe || rm ormat.exe
	test ! -f ormat_linux.zip || rm ormat_linux.zip
	test ! -f ormat_mac.zip || rm ormat_mac.zip
	test ! -f ormat_windows.zip || rm ormat_windows.zip

all: # 构建
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${opts} ${flags} -o ${execveFile} main.go
	tar czvf ormat_linux.zip ${execveFile}
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build ${opts} ${flags} -o ${execveFile}.exe main.go
	tar czvf ormat_windows.zip ormat.exe
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build ${opts} ${flags} -o ${execveFile} main.go
	tar czvf ormat_mac.zip ${execveFile}