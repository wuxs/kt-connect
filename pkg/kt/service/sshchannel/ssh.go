package sshchannel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/util"
	"io"
	"io/ioutil"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/wzshiming/socks5"
	"golang.org/x/crypto/ssh"
)

type SocksLogger struct {}

func (s SocksLogger) Println(v ...any) {
	_, _ = util.BackgroundLogger.Write([]byte(fmt.Sprint(v...) + util.Eol))
}

// StartSocks5Proxy start socks5 proxy
func (c *Cli) StartSocks5Proxy(privateKey, sshAddress, socks5Address string) (err error) {
	conn, err := createSshConnection(privateKey, sshAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	svc := &socks5.Server{
		Logger: SocksLogger{},
		ProxyDial: func(ctx context.Context, network string, address string) (net.Conn, error) {
			return conn.Dial(network, address)
		},
	}
	return svc.ListenAndServe("tcp", socks5Address)
}

// RunScript run the script on remote host.
func (c *Cli) RunScript(privateKey, sshAddress, script string) (result string, err error) {
	conn, err := createSshConnection(privateKey, sshAddress)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create ssh tunnel")
		return "", err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to create ssh session")
		return "", err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	err = session.Run(script)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to run ssh script")
		return "", err
	}
	output := stdoutBuf.String()
	return output, nil
}

// ForwardRemoteToLocal forward remote request to local
func (c *Cli) ForwardRemoteToLocal(privateKey, sshAddress, remoteEndpoint, localEndpoint string) error {
	// Handle incoming connections on reverse forwarded tunnel
	conn, err := createSshConnection(privateKey, sshAddress)
	if err != nil {
		log.Debug().Err(err).Msgf("Failed to create ssh tunnel")
		return err
	}

	// Listen on remote server port of shadow pod, via ssh connection
	listener, err := conn.Listen("tcp", remoteEndpoint)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to listen remote endpoint")
		_ = conn.Close()
		disconnectRemotePort(privateKey, sshAddress, remoteEndpoint, c)
		return err
	}

	log.Info().Msgf("Reverse tunnel %s -> %s established", remoteEndpoint, localEndpoint)
	for {
		if err = handleRequest(listener, localEndpoint); errors.Is(err, io.EOF) {
			_ = listener.Close()
			_ = conn.Close()
			return err
		}
	}
}

func disconnectRemotePort(privateKey, sshAddress, remoteEndpoint string, c *Cli) {
	remotePort := strings.Split(remoteEndpoint, ":")[1]
	out, err := c.RunScript(privateKey, sshAddress, fmt.Sprintf("/disconnect.sh %s", remotePort))
	if out != "" {
		_, _ = util.BackgroundLogger.Write([]byte(out + util.Eol))
	}
	if err != nil {
		log.Warn().Err(err).Msgf("Failed to disconnect remote port %s", remotePort)
	}
}

func handleRequest(listener net.Listener, localEndpoint string) error {
	defer func() {
		if r := recover(); r != nil {
			log.Error().Msgf("Failed to handle request: %v", r)
		}
	}()

	// Wait requests from remote endpoint
	client, err := listener.Accept()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to accept remote request")
		if !errors.Is(err, io.EOF) {
			time.Sleep(2 * time.Second)
		}
		return err
	}

	// Open a (local) connection to localEndpoint whose content will be forwarded to remoteEndpoint
	local, err := net.Dial("tcp", localEndpoint)
	if err != nil {
		_ = client.Close()
		log.Error().Err(err).Msgf("Local service error")
		return err
	}

	// Handle request in individual coroutine, current coroutine continue to accept more requests
	go handleClient(client, local)
	return nil
}

func createSshConnection(privateKey, address string) (*ssh.Client, error) {
	key, err := ioutil.ReadFile(privateKey)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
	}

	return ssh.Dial("tcp", address, config)
}

func handleClient(client net.Conn, remote net.Conn) {
	done := make(chan int)

	// Start remote -> local data transfer
	remoteReader := util.NewInterpretableReader(remote)
	go func() {
		defer handleBrokenTunnel(done)
		if _, err := io.Copy(client, remoteReader); err != nil {
			log.Warn().Err(err).Msgf("Error while copy remote->local")
		}
		done<-1
	}()

	// Start local -> remote data transfer
	localReader := util.NewInterpretableReader(client)
	go func() {
		defer handleBrokenTunnel(done)
		if _, err := io.Copy(remote, localReader); err != nil {
			log.Warn().Err(err).Msgf("Error while copy local->remote")
		}
		done<-1
	}()

	<-done
	remoteReader.Cancel()
	localReader.Cancel()
	_ = remote.Close()
	_ = client.Close()
}

func handleBrokenTunnel(done chan int) {
	if r := recover(); r != nil {
		log.Error().Msgf("Ssh tunnel broken: %v", r)
		done<-1
	}
}
