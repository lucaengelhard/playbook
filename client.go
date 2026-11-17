package main

import (
	"errors"
	"fmt"
	"log"
	"path"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	dialog *dialog.FormDialog
	sftp   *sftp.Client
	ssh    *ssh.Client
}

func NewClient(window fyne.Window) *Client {
	client := &Client{}
	client.dialog = client.createForm(window)
	return client
}

func (client *Client) Connect(host, port, user, password string) error {
	// TODO Check if already connected
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
		return err
	}
	sftp_client, err := sftp.NewClient(ssh_connection)
	if err != nil {
		log.Println(err)
		return err
	}
	client.ssh = ssh_connection
	client.sftp = sftp_client
	return nil
}

func (client *Client) Close() {
	if client.sftp != nil {
		client.sftp.Close()
	}

	if client.ssh != nil {
		client.ssh.Close()
	}
}

func (client *Client) Upload(localFile []byte, remotePath string) error {
	// TODO connect if not connected

	if client.sftp == nil {
		client.dialog.Show()
		return errors.New("Connect before uploading")
	}

	dirPath := path.Dir(remotePath)
	err := client.sftp.MkdirAll(dirPath)
	if err != nil {
		log.Println("Error while creating remote directory:", dirPath, err)
		return err
	}

	remoteFile, err := client.sftp.Create(remotePath)
	if err != nil {
		log.Println("Error while creating remote file", remotePath, err)
		return err
	}
	defer remoteFile.Close()

	_, err = remoteFile.Write(localFile)
	if err != nil {
		log.Println("Error while writing to remote file", remotePath, err)
		return err
	}
	log.Println("File", remoteFile.Name(), "written successfully")
	return nil
}

func (client *Client) createForm(window fyne.Window) *dialog.FormDialog {
	hostEntry := widget.NewEntry()
	portEntry := widget.NewEntry()
	userEntry := widget.NewEntry()
	pwdEntry := widget.NewPasswordEntry()

	items := []*widget.FormItem{
		{Text: "Host", Widget: hostEntry},
		{Text: "Port", Widget: portEntry},
		{Text: "User", Widget: userEntry},
		{Text: "Password", Widget: pwdEntry},
	}

	return dialog.NewForm("Connect to SFTP server", "Confirm", "Cancel", items, func(confirm bool) {
		if !confirm {
			return
		}

		err := client.Connect(hostEntry.Text, portEntry.Text, userEntry.Text, pwdEntry.Text)
		if err != nil {
			// TODO
			return
		}
		hostEntry.SetText("")
		portEntry.SetText("")
		userEntry.SetText("")
		pwdEntry.SetText("")
	}, window)
}
