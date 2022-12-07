LDFLAGS="-s -w -X main.AppVersion=1-alpha -X main.GitCommit=$(GIT_COMMIT)"

all:
	mkdir release
	go mod tidy
	go mod verify
	GOARCH="amd64" GOHOSTARCH="amd64" GOHOSTOS="linux" GOOS="linux" go build -o release/ham-linux-amd64 -ldflags ${LDFLAGS} cmd/ham/ham.go
	GOARCH="386" GOHOSTARCH="amd64" GOHOSTOS="linux" GOOS="linux" go build -o release/ham-linux-i386 -ldflags ${LDFLAGS} cmd/ham/ham.go
	GOARCH="arm64" GOHOSTARCH="amd64" GOHOSTOS="linux" GOOS="linux" go build -o release/ham-linux-arm64 -ldflags ${LDFLAGS} cmd/ham/ham.go
	GOARCH="arm64" GOHOSTARCH="amd64" GOHOSTOS="linux" GOOS="android" go build -o release/ham-android-arm64 -ldflags ${LDFLAGS} cmd/ham/ham.go
	GOARCH="amd64" GOHOSTARCH="amd64" GOHOSTOS="linux" GOOS="windows" go build -o release/ham-windows-amd64.exe -ldflags ${LDFLAGS} cmd/ham/ham.go
	GOARCH="amd64" GOHOSTARCH="amd64" GOHOSTOS="linux" GOOS="darwin" go build -o release/ham-macos-amd64 -ldflags ${LDFLAGS} cmd/ham/ham.go
	GOARCH="arm64" GOHOSTARCH="amd64" GOHOSTOS="linux" GOOS="darwin" go build -o release/ham-macos-arm64 -ldflags ${LDFLAGS} cmd/ham/ham.go
	GOARCH="amd64" GOHOSTARCH="amd64" GOHOSTOS="linux" GOOS="linux" go build -o release/ham-build-linux-amd64 -ldflags ${LDFLAGS} cmd/ham-build/ham-build.go

clean:
	rm -rf release

