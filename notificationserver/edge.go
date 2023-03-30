package notificationserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// PushNotificationServer push notification to rest api server. if jsonFormat is set to false, will Marshal useing bson
func PushNotificationServer(edgeURL string, targetMap map[string]string, message interface{}, jsonFormat bool) error {
	var err error

	// setup notification
	notf, err := setNotification(targetMap, message, jsonFormat)
	if err != nil {
		return err
	}

	// push notification
	client := http.Client{}
	for i := 0; i < 3; i++ {
		if err = sendCommandToEdge(&client, edgeURL, notf); err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
		err = fmt.Errorf("error sending url: '%s', reason: %s", edgeURL, err.Error())
	}
	return err

}

// sendCommandToEdge sends the HTTP request
func sendCommandToEdge(client *http.Client, edgeURL string, message []byte) error {

	req, err := http.NewRequest("POST", edgeURL, bytes.NewReader(message))
	req.Close = true
	if err != nil {
		return fmt.Errorf("failed to SendCommandToCluster, url: %s, data: %s, reason: %s", edgeURL, string(message), err.Error())
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to SendCommandToCluster, url: %s, data: %s, reason: %s", edgeURL, string(message), err.Error())
	}
	defer resp.Body.Close()
	respStr, err := httpRespToString(resp)
	if err != nil {
		return fmt.Errorf("failed to SendCommandToCluster, url: %s, data: %s, reason: %s, response: %s", edgeURL, string(message), err.Error(), respStr)
	}
	return nil
}

func setNotification(targetMap map[string]string, message interface{}, jsonFormat bool) ([]byte, error) {
	notification := Notification{
		Target:       targetMap,
		Notification: message,
	}

	var err error
	var m []byte
	if jsonFormat {
		if m, err = json.Marshal(notification); err != nil {
			err = fmt.Errorf("failed marshaling message to bson. message: '%v', reason: '%s'", notification, err.Error())
		}
	} else {

		if m, err = bson.Marshal(notification); err != nil {
			err = fmt.Errorf("failed marshaling message to bson. message: '%v', reason: '%s'", notification, err.Error())
		}
	}
	return m, err
}

// HTTPRespToString parses the body as string and checks the HTTP status code
func httpRespToString(resp *http.Response) (string, error) {
	if resp == nil {
		return "", fmt.Errorf("empty response")
	}
	strBuilder := strings.Builder{}
	if resp.ContentLength > 0 {
		strBuilder.Grow(int(resp.ContentLength))
	}
	_, err := io.Copy(&strBuilder, resp.Body)
	if err != nil {
		return strBuilder.String(), err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		err = fmt.Errorf("response status: %d. content: %s", resp.StatusCode, strBuilder.String())
	}
	return strBuilder.String(), err
}
