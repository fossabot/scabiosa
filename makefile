all: scabiosa-win32 scabiosa-win64 scabiosa-x86 scabiosa-x64
windows: scabiosa-win32 scabiosa-win64
linux: scabiosa-x86 scabiosa-x64

scabiosa-win32:
	GOOS=windows GOARCH=386 go build -o scabiosa-Win32.exe
scabiosa-win64:
	GOOS=windows GOARCH=amd64 go build -o scabiosa-Win64.exe
scabiosa-x86:
	GOOS=linux GOARCH=386 go build -o scabiosa-x86
scabiosa-x64:
	GOOS=linux GOARCH=amd64 go build -o scabiosa