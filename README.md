<div align="center">
<h1>ghb (gRPC HTTP Bridge)</h1>
</div>

# About

GHB is a lightweight Go package that allows you to expose gRPC services as HTTP endpoints with minimal configuration. It provides a bridge between gRPC and HTTP, enabling you to serve your gRPC methods through HTTP endpoints.

## Installation

```bash
go get github.com/malayanand/ghb
```

## Usage

### 1. Define HTTP Rules in Proto Files

First, import and use the HTTP rule option in your proto files:

```protobuf
import "github.com/malayanand/ghb/api/http.proto";

service YourService {
    rpc YourMethod(Request) returns (Response) {
        option (ghb.api.http) = {
            path: "/api/your-endpoint"
            method: GET  // Supported methods: GET, POST, HEAD
        };
    }
}
```

### 2. Create and Configure Server

```go
package main

import (
    "net"
    "github.com/malayanand/ghb"
)

func main() {
    // Create a new server
    server := ghb.NewServer()

    // Register your gRPC service
    server.RegisterService(&YourService_ServiceDesc, &yourServiceImpl{})

    // Create listener
    lis, err := net.Listen("tcp", ":8080")
    if err != nil {
        panic(err)
    }

    // Start serving HTTP requests
    if err := server.Serve(lis); err != nil {
        panic(err)
    }
}
```

## Features

- Automatic mapping of gRPC methods to HTTP endpoints
- Support for GET, POST, and HEAD HTTP methods
- JSON request/response handling
- URL parameter extraction
- Custom unmarshalling support
- Simple integration with existing gRPC services

## Examples

### Basic Example

Here's a complete example of how to expose a gRPC service as HTTP endpoints:

```protobuf
// user.proto
syntax = "proto3";

import "github.com/malayanand/ghb/api/http.proto";

service UserService {
    rpc GetUser(GetUserRequest) returns (User) {
        option (ghb.api.http) = {
            path: "/api/users/{id}"  // URL parameter 'id' will be extracted
            method: GET
        };
    }
}

message GetUserRequest {
    string id = 1;
}

message User {
    string id = 1;
    string name = 2;
    int32 age = 3;
}
```

### Custom Unmarshaller Example

You can implement custom unmarshalling logic for your messages by implementing the `Unmarshaler` interface:

```go
// Custom type that needs special unmarshalling
type TimeRange struct {
    StartTime time.Time
    EndTime   time.Time
}

// Implement the Unmarshaler interface
func (t *TimeRange) UnmarshalGHB(data interface{}) error {
    timeStr, ok := data.(string)
    if !ok {
        return fmt.Errorf("expected string, got %T", data)
    }
    
    // Parse time range string (e.g., "2023-01-01/2023-12-31")
    times := strings.Split(timeStr, "/")
    if len(times) != 2 {
        return fmt.Errorf("invalid time range format")
    }
    
    start, err := time.Parse("2006-01-02", times[0])
    if err != nil {
        return fmt.Errorf("invalid start time: %v", err)
    }
    
    end, err := time.Parse("2006-01-02", times[1])
    if err != nil {
        return fmt.Errorf("invalid end time: %v", err)
    }
    
    t.StartTime = start
    t.EndTime = end
    return nil
}

// Usage in proto message
message ReportRequest {
    TimeRange time_range = 1 [(ghb.api.field) = {json_name: "timeRange"}];
}
```

### Custom Marshaller Example

Similarly, you can implement custom marshalling logic for your messages by implementing the `Marshaler` interface:

```go
// Implement the Marshaler interface
func (t *TimeRange) MarshalGHB() (interface{}, error) {
    // Format to "2023-01-01/2023-12-31" format
    formattedTime := fmt.Sprintf("%s/%s", 
        t.StartTime.Format("2006-01-02"),
        t.EndTime.Format("2006-01-02"))
    
    return formattedTime, nil
}
```

### JSON Field Names Example

GHB allows you to customize the JSON field names used in serialization and deserialization:

```protobuf
message User {
    string user_id = 1 [(ghb.api.field) = {json_name: "userId"}];
    string first_name = 2 [(ghb.api.field) = {json_name: "firstName"}];
    string last_name = 3 [(ghb.api.field) = {json_name: "lastName"}];
}
```

This allows your HTTP API to use camelCase JSON naming convention while your protocol buffers follow snake_case convention. When this message is serialized to JSON, it will look like:

```json
{
  "userId": "123",
  "firstName": "John",
  "lastName": "Doe"
}
```

### URL Parameters Example

GHB automatically extracts URL parameters from the path:

```protobuf
service OrderService {
    rpc GetOrder(GetOrderRequest) returns (Order) {
        option (ghb.api.http) = {
            path: "/api/users/{userId}/orders/{orderId}"
            method: GET
        };
    }
}

message GetOrderRequest {
    string user_id = 1;
    string order_id = 2;
}
```

When making a request to `/api/users/123/orders/456`, GHB will automatically populate:
- `user_id` with "123"
- `order_id` with "456"

