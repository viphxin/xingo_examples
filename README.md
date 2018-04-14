# xingo_examples
xingo框架的教程例子

### 安装
```
go get -v github.com/viphxin/xingo_examples
```

------
go get -u github.com/golang/protobuf/protoc-gen-go

Downloads\protoc-3.3.0-win32\bin>protoc.exe --plugin=protoc-gen-go=%GOPATH%\bin\protoc-gen-go.exe --go_out  F:\workspace\src\xingo_examples\helloword\pb -I  F:\workspace\src\xingo_examples\helloword\pb F:\workspace\src\xingo_examples\helloword\pb\msg.proto