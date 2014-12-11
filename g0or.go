package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/user"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

const appName = "go0r"

var errAuthenticationFailed = errors.New(":)")

func logParameters(conn ssh.ConnMetadata) logrus.Fields {
	return logrus.Fields{
		"user":                conn.User(),
		"session_id":          string(conn.SessionID()),
		"client_version":      string(conn.ClientVersion()),
		"server_version":      string(conn.ServerVersion()),
		"remote_addr_network": string(conn.RemoteAddr().Network()),
		"remote_addr_address": string(conn.RemoteAddr().String()),
		"local_addr_network":  string(conn.LocalAddr().Network()),
		"local_addr_address":  string(conn.LocalAddr().String()),
	}
}

func authenticatePassword(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	logrus.WithFields(logParameters(conn)).Info(fmt.Sprintf("Request with password: %s ", password))
	return nil, errAuthenticationFailed
}

func authenticateKey(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	logrus.WithFields(logParameters(conn)).Info(fmt.Sprintf("Request with keytype: %s ", key.Type()))
	return nil, errAuthenticationFailed
}

func init() {
	usr, err := user.Current()
	if err != nil {
		logrus.Warn(err)
	}
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Warn(err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/" + appName)
	viper.AddConfigPath(usr.HomeDir + "/" + appName)
	viper.AddConfigPath(dir + "/configs/")

	viper.BindEnv("port", "PORT")
	viper.SetDefault("port", ":22")

	viper.BindEnv("host_key", "HOST_KEY")
	viper.SetDefault("host_key", "./host_key")

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warn("Can not load config file. defaults are loaded!")
	}
}

func main() {
	config := ssh.ServerConfig{
		PasswordCallback:  authenticatePassword,
		PublicKeyCallback: authenticateKey,
	}

	keyPath := viper.GetString("host_key")
	hostPrivateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		logrus.Panic(err)
	}
	hostPrivateKeySigner, err := ssh.ParsePrivateKey(hostPrivateKey)
	if err != nil {
		logrus.Panic(err)
	}
	config.AddHostKey(hostPrivateKeySigner)
	socket, err := net.Listen("tcp", viper.GetString("port"))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := socket.Accept()
		if err != nil {
			logrus.Panic(err)
		}
		_, _, _, err = ssh.NewServerConn(conn, &config)
		if err == nil {
			logrus.Panic("Successful login? why!?")
		}
	}
}
