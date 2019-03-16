
//<developer>
//    <name>linapex 曹一峰</name>
//    <email>linapex@163.com</email>
//    <wx>superexc</wx>
//    <qqgroup>128148617</qqgroup>
//    <url>https://jsq.ink</url>
//    <role>pku engineer</role>
//    <date>2019-03-16 19:16:33</date>
//</624450069680099328>


package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/ethereum/go-ethereum/log"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

//ssh client是Go的ssh客户机的一个小包装，有几个实用方法
//在顶部实施。
type sshClient struct {
server  string //没有端口号的服务器名或IP
address string //远程服务器的IP地址
pubkey  []byte //用于对服务器进行身份验证的RSA公钥
	client  *ssh.Client
	logger  log.Logger
}

//拨号使用当前用户和
//
//
func dial(server string, pubkey []byte) (*sshClient, error) {
//
	hostname := ""
	hostport := server
	username := ""
identity := "id_rsa" //违约

	if strings.Contains(server, "@") {
		prefix := server[:strings.Index(server, "@")]
		if strings.Contains(prefix, ":") {
			username = prefix[:strings.Index(prefix, ":")]
			identity = prefix[strings.Index(prefix, ":")+1:]
		} else {
			username = prefix
		}
		hostport = server[strings.Index(server, "@")+1:]
	}
	if strings.Contains(hostport, ":") {
		hostname = hostport[:strings.Index(hostport, ":")]
	} else {
		hostname = hostport
		hostport += ":22"
	}
	logger := log.New("server", server)
	logger.Debug("Attempting to establish SSH connection")

	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	if username == "" {
		username = user.Username
	}
//
	var auths []ssh.AuthMethod

	path := filepath.Join(user.HomeDir, ".ssh", identity)
	if buf, err := ioutil.ReadFile(path); err != nil {
		log.Warn("No SSH key, falling back to passwords", "path", path, "err", err)
	} else {
		key, err := ssh.ParsePrivateKey(buf)
		if err != nil {
			fmt.Printf("What's the decryption password for %s? (won't be echoed)\n>", path)
			blob, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			if err != nil {
				log.Warn("Couldn't read password", "err", err)
			}
			key, err := ssh.ParsePrivateKeyWithPassphrase(buf, blob)
			if err != nil {
				log.Warn("Failed to decrypt SSH key, falling back to passwords", "path", path, "err", err)
			} else {
				auths = append(auths, ssh.PublicKeys(key))
			}
		} else {
			auths = append(auths, ssh.PublicKeys(key))
		}
	}
	auths = append(auths, ssh.PasswordCallback(func() (string, error) {
		fmt.Printf("What's the login password for %s at %s? (won't be echoed)\n> ", username, server)
		blob, err := terminal.ReadPassword(int(os.Stdin.Fd()))

		fmt.Println()
		return string(blob), err
	}))
//
	addr, err := net.LookupHost(hostname)
	if err != nil {
		return nil, err
	}
	if len(addr) == 0 {
		return nil, errors.New("no IPs associated with domain")
	}
//
	logger.Trace("Dialing remote SSH server", "user", username)
	keycheck := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
//
		if pubkey == nil {
			fmt.Println()
			fmt.Printf("The authenticity of host '%s (%s)' can't be established.\n", hostname, remote)
			fmt.Printf("SSH key fingerprint is %s [MD5]\n", ssh.FingerprintLegacyMD5(key))
			fmt.Printf("Are you sure you want to continue connecting (yes/no)? ")

			text, err := bufio.NewReader(os.Stdin).ReadString('\n')
			switch {
			case err != nil:
				return err
			case strings.TrimSpace(text) == "yes":
				pubkey = key.Marshal()
				return nil
			default:
				return fmt.Errorf("unknown auth choice: %v", text)
			}
		}
//
		if bytes.Equal(pubkey, key.Marshal()) {
			return nil
		}
//
		return errors.New("ssh key mismatch, readd the machine to update")
	}
	client, err := ssh.Dial("tcp", hostport, &ssh.ClientConfig{User: username, Auth: auths, HostKeyCallback: keycheck})
	if err != nil {
		return nil, err
	}
//
	c := &sshClient{
		server:  hostname,
		address: addr[0],
		pubkey:  pubkey,
		client:  client,
		logger:  logger,
	}
	if err := c.init(); err != nil {
		client.Close()
		return nil, err
	}
	return c, nil
}

//
//
func (client *sshClient) init() error {
	client.logger.Debug("Verifying if docker is available")
	if out, err := client.Run("docker version"); err != nil {
		if len(out) == 0 {
			return err
		}
		return fmt.Errorf("docker configured incorrectly: %s", out)
	}
	client.logger.Debug("Verifying if docker-compose is available")
	if out, err := client.Run("docker-compose version"); err != nil {
		if len(out) == 0 {
			return err
		}
		return fmt.Errorf("docker-compose configured incorrectly: %s", out)
	}
	return nil
}

//
func (client *sshClient) Close() error {
	return client.client.Close()
}

//
//
func (client *sshClient) Run(cmd string) ([]byte, error) {
//建立一个单一的命令会话
	session, err := client.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

//
	client.logger.Trace("Running command on remote server", "cmd", cmd)
	return session.CombinedOutput(cmd)
}

//
//
func (client *sshClient) Stream(cmd string) error {
//建立一个单一的命令会话
	session, err := client.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

//
	client.logger.Trace("Streaming command on remote server", "cmd", cmd)
	return session.Run(cmd)
}

//
//同时存在文件夹。
func (client *sshClient) Upload(files map[string][]byte) ([]byte, error) {
//建立一个单一的命令会话
	session, err := client.client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

//创建流式处理SCP内容的goroutine
	go func() {
		out, _ := session.StdinPipe()
		defer out.Close()

		for file, content := range files {
			client.logger.Trace("Uploading file to server", "file", file, "bytes", len(content))

fmt.Fprintln(out, "D0755", 0, filepath.Dir(file))             //
fmt.Fprintln(out, "C0644", len(content), filepath.Base(file)) //创建实际文件
out.Write(content)                                            //
fmt.Fprint(out, "\x00")                                       //传输结束为\x00
fmt.Fprintln(out, "E")                                        //离开目录（更简单）
		}
	}()
	return session.CombinedOutput("/usr/bin/scp -v -tr ./")
}

