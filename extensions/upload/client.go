package main

import (
	"encoding/json/v2"
	"fmt"
	"io"
	"net/url"
	"os"
	"uuid"

	"github.com/gorilla/websocket"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	NlPort         string `json:"nlPort"`
	NlToken        string `json:"nlToken"`
	NlExtensionId  string `json:"nlExtensionId"`
	NlConnectToken string `json:"nlConnectToken"`
}

type EventMessage struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type DataPacketData struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

type DataPacket struct {
	Id          string         `json:"id"`
	Method      string         `json:"method"`
	AccessToken string         `json:"accessToken"`
	Data        DataPacketData `json:"data"`
}

type Client struct {
	url    url.URL
	config *Config
	socket *websocket.Conn
	ssh    *ssh.Client
	sftp   *sftp.Client
}

func (client *Client) send(event string, data map[string]interface{}) {
	msg := DataPacket{
		Id:          uuid.New().String(),
		Method:      "app.broadcast",
		AccessToken: client.config.NlToken,
		Data: DataPacketData{
			Event: event,
			Data:  data,
		},
	}

	msg_str, _ := json.Marshal(msg)

	client.socket.WriteMessage(websocket.TextMessage, []byte(msg_str))
}

func (client *Client) error(msg string) {
	out := make(map[string]interface{})
	out["msg"] = msg
	client.send("error", out)
}

func (client *Client) listen(callback func(client *Client, message EventMessage)) {
	client.socket, _, _ = websocket.DefaultDialer.Dial(client.url.String(), nil)

	for {
		_, msg, err := client.socket.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		d := EventMessage{}
		json.Unmarshal(msg, &d)

		if d.Event == "windowClose" || d.Event == "appclose" {
			client.quit()
		}

		callback(client, d)
	}
}

func (client *Client) connect(host string, port int, user string, password string) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client.ssh, _ = ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	client.sftp, _ = sftp.NewClient(client.ssh)
	client.send("sftp_connected", map[string]any{}) // TODO maybe transfer error state or smth?
}

func (client *Client) quit() {
	if client.socket != nil {
		client.socket.Close()
	}

	if client.sftp != nil {
		client.sftp.Close()
	}

	if client.ssh != nil {
		client.ssh.Close()
	}

	pid := os.Getpid()
	process, _ := os.FindProcess(pid)
	process.Signal(os.Kill)
}

func (client *Client) upload(localPath, remotePath string) {
	localFile, err := os.Open(localPath)
	if err != nil {
		fmt.Println("Failed to open local file:", err)
		return
	}
	defer localFile.Close()

	remoteFile, err := client.sftp.Create(remotePath)
	if err != nil {
		fmt.Println("Failed to create remote file:", err)
		return
	}
	defer remoteFile.Close()

	_, err = io.Copy(remoteFile, localFile)
	if err != nil {
		fmt.Println("Failed to upload file:", err)
		return
	}

	client.send("upload", map[string]any{
		"from": localPath,
		"to":   remotePath,
	})
}
