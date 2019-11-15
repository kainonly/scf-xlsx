package service

import (
	"golang.org/x/crypto/ssh"
	"strconv"
)

type Client struct {
	options map[string]*ConnectOption
	runtime map[string]*ssh.Client
}

func InjectClient() *Client {
	client := Client{}
	if client.options == nil {
		client.options = make(map[string]*ConnectOption)
	}
	if client.runtime == nil {
		client.runtime = make(map[string]*ssh.Client)
	}
	return &client
}

type ConnectOption struct {
	Host       string
	Port       uint
	Username   string
	Password   string
	Key        []byte
	PassPhrase []byte
}

// Factory SSH AuthMethod
func (c *Client) authMethod(option ConnectOption) (auth []ssh.AuthMethod, err error) {
	if option.Key == nil {
		// Password AuthMethod
		auth = []ssh.AuthMethod{
			ssh.Password(option.Password),
		}
	} else {
		// PrivateKey AuthMethod
		var signer ssh.Signer
		if option.PassPhrase != nil {
			// With Passphrase
			if signer, err = ssh.ParsePrivateKeyWithPassphrase(
				option.Key,
				option.PassPhrase,
			); err != nil {
				return
			}
		} else {
			// Without Passphrase
			if signer, err = ssh.ParsePrivateKey(
				option.Key,
			); err != nil {
				return
			}
		}
		auth = []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		}
	}
	return
}

func (c *Client) connect(option ConnectOption) (client *ssh.Client, err error) {
	auth, err := c.authMethod(option)
	if err != nil {
		return
	}
	config := ssh.ClientConfig{
		User:            option.Username,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := option.Host + ":" + strconv.Itoa(int(option.Port))
	client, err = ssh.Dial("tcp", addr, &config)
	return
}

// SSH Connect Testing
func (c *Client) Testing(option ConnectOption) (client *ssh.Client, err error) {
	return c.connect(option)
}

func (c *Client) Put(identity string, option ConnectOption) (client *ssh.Client, err error) {
	client, err = c.connect(option)
	if err != nil {
		return
	}
	c.options[identity] = &option
	c.runtime[identity] = client
	return
}

type GetResult struct {
	Identity  string `json:"identity"`
	Host      string `json:"host"`
	Port      uint   `json:"port"`
	Username  string `json:"username"`
	Connected bool   `json:"connected"`
}

func (c *Client) Get(identity string) (exists bool, result GetResult) {
	exists = c.options[identity] != (&ConnectOption{})
	option := c.options[identity]
	result = GetResult{
		Identity:  identity,
		Host:      option.Host,
		Port:      option.Port,
		Username:  option.Username,
		Connected: true,
	}
	return
}