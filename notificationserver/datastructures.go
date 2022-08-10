package notificationserver

// Notification is a notification passed between the Gateway servers
type Notification struct {
	// Target for the notification
	//
	// Describes attributes of the target that should receive the notification
	//
	// Required: true
	// Example: { "cluster": "minikube", "customerGUID": "b5b28ef9-d297-4a93-aec4-22de5b21e802", "component": "websocket" }
	Target map[string]string `json:"target"`
	// Whether to send the message synchronously
	//
	// If `true`, waits for the message to be sent, else the message is sent asynchronously.
	// Example: true
	// Default: false
	SendSynchronicity bool `json:"sendSynchronicity"`
	// Body of the notification
	//
	// Example: { "commands": [ { "CommandName": "kubescapeScan", "args": { "scanV1": { "submit": true, "excludeNamespaces": ["kube-system"] } } } ] }
	Notification interface{} `json:"notification"`
}
