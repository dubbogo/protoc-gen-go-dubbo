# protoc-gen-go-dubbo

## Prerequisites
Before using `protoc-gen-go-dubbo`, make sure you have the following prerequisites installed on your system:
- Go (version 1.17 or higher)
- Protocol Buffers (version 3.0 or higher)

## Installation
To install `protoc-gen-go-dubbo`, you can use `go install`:
```shell
go install github.com/dubbogo/protoc-gen-go-dubbo
```

Or you can clone this repository and build it manually:
```shell
git clone github.com/dubbogo/protoc-gen-go-dubbo
cd protoc-gen-go-dubbo
go build
```

## Usage
To generate dubbo code using `protoc-gen-go-dubbo`, you can use the following command:
```shell
protoc -I ./ \
  --go-dubbo_out=./ --go-dubbo_opt=paths=source_relative \
  ./greet.proto
```
The `--go-dubbo_out` option specifies the output directory for the generated code,
and `--go-dubbo_opt=paths=source_relative` sets the output path to be relative to the source file.

You can also use the flags `in_file_method_protocol_spec=true` or `in_file_method_protocol_spec=true` to enable in-file transport protocol specification, then protoc-gen-go-dubbo will only generate the service/method with the option "DUBBO". For example:

```shell
protoc -I ./ \
  --go-hessian2_out=./ --go-hessian2_opt=paths=source_relative \
  --go-dubbo_out=./ --go-dubbo_opt=paths=source_relative,in_file_method_protocol_spec=true \
  ./greet.proto
```

## Example
To generate Dubbo code, you can create a `.proto` file with the following content:
```protobuf
syntax = "proto3";

package greet;

option go_package = "some_path/greet;greet";

import "unified_idl_extend/unified_idl_extend.proto";

message GreetRequest {
  string name = 1;

  option (unified_idl_extend.message_extend) = {
    java_class_name: "org.apache.greet.GreetRequest";
  };
}

message GreetResponse {
  string greeting = 1;

  option (unified_idl_extend.message_extend) = {
    java_class_name: "org.apache.greet.GreetResponse";
  };
}

service GreetService {
  rpc Greet(GreetRequest) returns (GreetResponse) {
    option (unified_idl_extend.method_extend) = {
      method_name: "greet";
    };
  }
}
```

Extra options with transport protocol specification look like the follows:

```proto
//method protocol spec
service GreetingsService  {
  rpc Greet(GreetRequest) returns (GreetResponse) {
    option (unified_idl_extend.method_protocols) = {
      protocol_names: ["DUBBO"];
    };
    option (unified_idl_extend.method_extend) = {
      method_name: "greet";
    };
  }
  rpc Greet2(GreetRequest) returns (GreetResponse) {
    option (unified_idl_extend.method_protocols) = {
      protocol_names: ["TRIPLE","DUBBO"];
    };
    option (unified_idl_extend.method_extend) = {
      method_name: "greet2";
    };
  }
}
```

```proto
service GreetingsService  {
//service protocol spec
option (unified_idl_extend.service_protocol) = {
    protocol_name: "DUBBO";
  };
  rpc Greet(GreetRequest) returns (GreetResponse) {
    option (unified_idl_extend.method_extend) = {
      method_name: "greet";
    };
  }
  rpc Greet2(GreetRequest) returns (GreetResponse) {
    option (unified_idl_extend.method_extend) = {
      method_name: "greet2";
    };
  }
}
```

Note that you need to import the `unified_idl_extend` package to use the `method_extend` options to extend the service.
And Dubbo protocol must be used with Hessian2 serialization, so you need to use the `message_extend` options to extend 
the message.

Then, you can run the following command to generate the Dubbo and Hessian2 code:
```shell
protoc -I ./ \
  --go-hessian2_out=./ --go-hessian2_opt=paths=source_relative \
  --go-dubbo_out=./ --go-dubbo_opt=paths=source_relative \
  ./greet.proto
```
This will generate the `greet.dubbo.go` and `greet.hessian2.go` file in the same directory as your `greet.proto` file:
```shell
.
├── greet.hessian2.go
├── greet.dubbo.go
├── greet.proto
└── unified_idl_extend
    ├── unified_idl_extend.pb.go
    └── unified_idl_extend.proto
```

The content of the `greet.hessian2.go` file will be:
```go
// Code generated by protoc-gen-go-dubbo. DO NOT EDIT.

// Source: greet.proto
// Package: greet

package greet

import (
	"context"

	"dubbo.apache.org/dubbo-go/v3"
	"dubbo.apache.org/dubbo-go/v3/client"
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/server"
)

const (
	// GreetServiceName is the fully-qualified name of the GreetService service.
	GreetServiceName = "greet.GreetService"

	// These constants are the fully-qualified names of the RPCs defined in this package. They're
	// exposed at runtime as procedure and as the final two segments of the HTTP route.
	//
	// Note that these are different from the fully-qualified method names used by
	// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
	// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
	// period.
	// GreetServiceGreetProcedure is the fully-qualified name of the GreetService's Greet RPC.'
	GreetServiceGreetProcedure = "/greet.GreetService/greet"
)

var (
	_ GreetService = (*GreetServiceImpl)(nil)
)

type GreetService interface {
	Greet(ctx context.Context, req *GreetRequest, opts ...client.CallOption) (*GreetResponse, error)
}

// NewGreetService constructs a client for the greet.GreetService service
func NewGreetService(cli *client.Client, opts ...client.ReferenceOption) (GreetService, error) {
	conn, err := cli.DialWithInfo("greet.GreetService", &GreetService_ClientInfo, opts...)
	if err != nil {
		return nil, err
	}
	return &GreetServiceImpl{
		conn: conn,
	}, nil
}

func SetConsumerService(srv common.RPCService) {
	dubbo.SetConsumerServiceWithInfo(srv, &GreetService_ClientInfo)
}

// GreetServiceImpl implements GreetService
type GreetServiceImpl struct {
	conn *client.Connection
}

func (c *GreetServiceImpl) Greet(ctx context.Context, req *GreetRequest, opts ...client.CallOption) (*GreetResponse, error) {
	resp := new(GreetResponse)
	if err := c.conn.CallUnary(ctx, []interface{}{req}, resp, "greet", opts...); err != nil {
		return nil, err
	}
	return resp, nil
}

var GreetService_ClientInfo = client.ClientInfo{
	InterfaceName: "greet.GreetService",
	MethodNames: []string{
		"greet",
	},
	ConnectionInjectFunc: func(dubboCliRaw interface{}, conn *client.Connection) {
		dubboCli := dubboCliRaw.(*GreetServiceImpl)
		dubboCli.conn = conn
	},
}

// GreetServiceHandler is an implementation of the greet.GreetService service.
type GreetServiceHandler interface {
	Greet(ctx context.Context, req *GreetRequest) (*GreetResponse, error)
}

func RegisterGreetServiceHandler(srv *server.Server, hdlr GreetServiceHandler, opts ...server.ServiceOption) error {
	return srv.Register(hdlr, &GreetService_ServiceInfo, opts...)
}

func SetProviderService(srv common.RPCService) {
	dubbo.SetProviderServiceWithInfo(srv, &GreetService_ServiceInfo)
}

var GreetService_ServiceInfo = server.ServiceInfo{
	InterfaceName: "greet.GreetService",
	ServiceType:   (*GreetServiceHandler)(nil),
	Methods: []server.MethodInfo{
		{
			Name: "greet",
			Type: constant.CallUnary,
			ReqInitFunc: func() interface{} {
				return new(GreetResponse)
			},
			MethodFunc: func(ctx context.Context, args []interface{}, handler interface{}) (interface{}, error) {
				req := args[0].(*GreetRequest)
				res, err := handler.(GreetServiceHandler).Greet(ctx, req)
				return res, err
			},
		},
	},
}
```