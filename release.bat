SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build -o release/qmailstorage-linux-amd64

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=arm
go build -o release/qmailstorage-linux-arm

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build -o release/qmailstorage-darwin-amd64

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -o release/qmailstorage-windows-amd64.exe