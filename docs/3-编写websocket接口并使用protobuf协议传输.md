## 配置Protobuf

#### 编写proto文件

Protobuf可以将编写的proto文件转变为其他语言可以使用的协议，可以作为设计安全的跨语言RPC接口的基础工具。
proto文件中最基本的数据单元是message，是类似Go语言中结构体的存在。在message中可以嵌套message或其他的基础数据类型的成员。

在 pkg/proto 目录下创建 client_message.proto 文件：

``` protobuf
syntax = "proto3";
package pb;

// 上行数据
message Input {
    string type = 1; // 包的类型
    bytes data = 2; // 数据
    repeated Auth auth = 3;
}

// 下行数据
message Output {
    string type = 1; // 包的类型
    int32 code = 2; // 错误码
    string message = 3; // 错误信息
    bytes data = 4; // 数据
}

// 权限判断
message Auth {
    int64 app_id = 1; // app_id
    int64 device_id = 2; // 设备id
    int64 user_id = 3; // 用户id
    string token = 4; // 秘钥
}

// 获取用户信息
message GetUserInfoReq {
    int64 user_id = 1;
}
message GetUserInfoResp {
    int64 user_id = 1;
    string nickname = 2;
    int32 sex = 3; // 性别 0 未知 1 男 2 女
    string avatar_url = 4;
    string sign = 5;
    string account = 6;
}
```

client_message.proto 文件里定义了很多 message，进入项目查看全部的：

#### 编译proto文件

安装编译工具 protoc，然后使用以下命令将 proto 文件编译为 Go 语言可以使用的协议：

``` shell
# 编译目录下所有protp文件
$ protoc --go_out=. *.proto

# 编译client_message.proto文件
$ protoc -I . --go_out=plugins=grpc:. client_message.proto
```

使用以下命令将 proto 文件编译为 Javascript 语言可以使用的协议：

``` shell
# 将proto编译为js文件
$ protoc --js_out=import_style=commonjs,binary:./ *.proto
# 将编译好的js文件转换为浏览器可以直接使用的js文件
$ browserify client_message_pb.js -o client_message.js
```
---

## 编写Websocket接口

### 