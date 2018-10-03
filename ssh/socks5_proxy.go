package ssh

import (
	"fmt"
	"time"

	"github.com/cloudfoundry/socks5-proxy"
	"golang.org/x/crypto/ssh"
)

// Socks5Proxy is based on github.com/cloudfoundry/socks5-proxy.Socks5Proxy#Dialer
// after establishing a connection it creates a session with a keepalive loop to
// prevent the tunnel closing during inactivity
type Socks5Proxy struct {
	hostKey             proxy.HostKey
	serverAliveInterval time.Duration
	terminate           chan struct{}
	logger              Logger
	client              *ssh.Client
}

func NewSocks5Proxy(hostKey proxy.HostKey, serverAliveInterval time.Duration, logger Logger) *Socks5Proxy {
	return &Socks5Proxy{
		hostKey:             hostKey,
		serverAliveInterval: serverAliveInterval,
		logger:              logger,
	}
}

func (s *Socks5Proxy) Dialer(username, key, url string) (proxy.DialFunc, error) {
	if username == "" {
		username = "jumpbox"
	}

	signer, err := ssh.ParsePrivateKey([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("parse private key: %s", err)
	}

	hostKey, err := s.hostKey.Get(username, key, url)
	if err != nil {
		return nil, fmt.Errorf("get host key: %s", err)
	}

	clientConfig := proxy.NewSSHClientConfig(username, ssh.FixedHostKey(hostKey), ssh.PublicKeys(signer))

	s.client, err = ssh.Dial("tcp", url, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("ssh dial: %s", err)
	}

	err = s.startKeepAliveLoop()
	if err != nil {
		return nil, fmt.Errorf("ssh keepalive: %s", err)
	}

	return s.client.Dial, nil
}

func (s *Socks5Proxy) Close() error {
	close(s.terminate)
	return s.client.Close()
}

func (s *Socks5Proxy) startKeepAliveLoop() error {
	session, err := s.client.NewSession()
	if err != nil {
		return fmt.Errorf("new session for keepalive: %s", err)
	}

	s.terminate = make(chan struct{})
	go func() {
		for {
			select {
			case <-s.terminate:
				return
			default:
				_, err := session.SendRequest("keepalive@bbr", true, nil)
				if err != nil {
					s.logger.Debug("ssh", "keepalive failed: %+v", err)
				}
				time.Sleep(time.Second * s.serverAliveInterval)
			}
		}
	}()

	return nil
}
