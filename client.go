package main

import (
	"fmt"
	"log"
	"path"

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

	ssh_connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), config)
	if err != nil {
		log.Println(err)
		return
	}
	sftp_client, err := sftp.NewClient(ssh_connection)
	if err != nil {
		log.Println(err)
		return
	}
	client.ssh = ssh_connection
	client.sftp = sftp_client
}

func (client *Client) Upload(localFile []byte, remotePath string) {
	// TODO connect if not connected

	dirPath := path.Dir(remotePath)
	err := client.sftp.MkdirAll(dirPath)
	if err != nil {
		log.Println("Error while creating remote directory:", dirPath, err)
		return
	}

	remoteFile, err := client.sftp.Create(remotePath)
	if err != nil {
		log.Println("Error while creating remote file", remotePath, err)
		return
	}
	defer remoteFile.Close()

	_, err = remoteFile.Write(localFile)
	if err != nil {
		log.Println("Error while writing to remote file", remotePath, err)
		return
	}
	log.Println("File", remoteFile.Name(), "written successfully")
}

func (client *Client) Close() {
	if client.sftp != nil {
		client.sftp.Close()
	}

	if client.ssh != nil {
		client.ssh.Close()
	}
}
