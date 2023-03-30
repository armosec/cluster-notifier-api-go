Notification Server
===================

A package for sending notifications to a REST API server. It provides the following functionalities:

1. `PushNotificationServer`: Sends a notification to a given REST API server.
2. `sendCommandToEdge`: Sends the actual HTTP request to the REST API server.
3. `setNotification`: Encodes the notification to JSON or BSON format based on the provided boolean.
4. `httpRespToString`: Reads the response body and returns it as a string.

Installation
------------
```
go get github.com/notificationserver
```


Usage
-----
```go
import "github.com/notificationserver"

// Send a notification to the REST API server
err := notificationserver.PushNotificationServer("http://localhost:8080/" + notificationserver.PathRESTV1, map[string]string{notificationserver.TargetCustomer: "111111-1111-1111-1111"}, notificationserver.Notification{}, true)
if err != nil {
    // Handle error
}
```

Contributing
-----

If you find any bugs or have a feature request, please open an issue or send a pull request.

License
-----

This project is licensed under the MIT License - see the LICENSE file for details.
