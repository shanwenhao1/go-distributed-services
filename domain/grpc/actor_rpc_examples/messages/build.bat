protoc -I=. -I=%GOPATH%\src --gogoslick_out=plugins=grpc:. protos.proto
protoc --proto_path=. -I=F:\GoPath\src -I=F:\GoPath\pkg\windows_amd64 --gogoslick_out=plugins=grpc:. test.proto