build:
	mkdir -p bin
	cd src; go build -o ../bin/srv42 .

compile:
	cd src; GOOS=darwin GOARCH=amd64 go build -o ../bin/srv42-darwin-amd64 .
	cd src; GOOS=darwin GOARCH=arm64 go build -o ../bin/srv42-darwin-arm64 .
	cd src; GOOS=linux GOARCH=amd64 go build -o ../bin/srv42-linux-amd64 .
	cd src; GOOS=windows GOARCH=amd64 go build -o ../bin/srv42-windows-amd64.exe .

	tar czf bin/srv42-darwin-amd64.tar.gz bin/srv42-darwin-amd64
	tar czf bin/srv42-darwin-arm64.tar.gz bin/srv42-darwin-arm64
	tar czf bin/srv42-linux-amd64.tar.gz bin/srv42-linux-amd64
	tar czf bin/srv42-windows-amd64.exe.tar.gz bin/srv42-windows-amd64.exe

	sha256sum bin/srv42-darwin-amd64.tar.gz > bin/srv42-darwin-amd64.tar.gz.sha256
	sha256sum bin/srv42-darwin-arm64.tar.gz > bin/srv42-darwin-arm64.tar.gz.sha256
	sha256sum bin/srv42-linux-amd64.tar.gz > bin/srv42-linux-amd64.tar.gz.sha256
	sha256sum bin/srv42-windows-amd64.exe.tar.gz > bin/srv42-windows-amd64.exe.tar.gz.sha256
