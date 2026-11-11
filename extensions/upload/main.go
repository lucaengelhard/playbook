package main

import (
	"bufio"
	"encoding/json/v2"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func onEvent(client *Client, data EventMessage) {
	if d, ok := data.Data.(map[string]any); ok {
		switch data.Event {
		case "upload":
			{
				fmt.Println(d)
				// TODO: client.upload(...) (how to handle multiple files?)
			}
		case "connect":
			{
				host := d["host"].(string)
				port, err := strconv.Atoi(d["port"].(string))
				user := d["user"].(string)
				password := d["password"].(string)

				if err != nil {
					client.error(err.Error())
					return
				}

				if host == "" || user == "" || password == "" {
					client.error("Invalid sftp config :(")
					return
				}

				fmt.Printf("%s %d %s %s\n", host, port, user, password)
				client.connect(host, port, user, password)
			}
		}
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	builder := strings.Builder{}
	for scanner.Scan() {
		builder.WriteString(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	var config Config
	err := json.Unmarshal([]byte(builder.String()), &config)
	if err != nil {
	}

	client := Client{
		url: url.URL{
			Scheme: "ws",
			Host:   "127.0.0.1:" + config.NlPort,
			Path:   "?extensionId=" + config.NlExtensionId + "&connectToken=" + config.NlConnectToken,
		},
		config: &config,
	}

	defer client.quit()
	client.listen(onEvent)
}
