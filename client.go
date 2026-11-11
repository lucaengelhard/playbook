package main

import (
	"fmt"
	"io"
	"os"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	sftp *sftp.Client
	ssh  *ssh.Client
}

func (client *Client) Connect(host, port, user, password string) {
	// TODO Check if already connected
	// TODO use mutex
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	ssh_connection, _ := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	sftp_client, _ := sftp.NewClient(ssh_connection)
	client.ssh = ssh_connection
	client.sftp = sftp_client
}

func (client *Client) Upload(localFile *os.File, remotePath string) {
	// TODO handle creating directories recursively
	// TODO connect if not connected
	// TODO use mutex

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
}
