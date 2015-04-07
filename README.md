
# Simple ssh honypot in Golang

A simple, around 100 line of code, ssh honeypot written in golang.

## The name
Its base on persian proverb "گریه کردن بر روی گور بدون مرده " (cry over the empty grave) and the word goor (گور means the Grave in persian).

## How to use it
This is the (wrong) steps :

```bash
go get -u -v github.com/fzerorubigd/go0r
# config folder could be $HOME/.config/go0r or /etc/go0r or ./config
CONFIG_FOLDER=/etc/go0r
# create a host key, password and anything is not important at all! just hit enter
ssh-keygen -f $CONFIG_FOLDER/host_key
# go0r port to use, normally :22 :) must run with sudo in that case, and do not forget :
echo "port=\":22\"" > $CONFIG_FOLDER/config.toml
# address of host key we create in secound step
echo "host_key=\"$CONFIG_FOLDER/host_key\"" >> $CONFIG_FOLDER/config.toml
# run the application!
$GOPATH/bin/go0r
```

And then try to login into ssh server on "port" and watch the output :)
Also, you can use GOOR_PORT , GOOR_HOST_KEY environment variables to set the config values.

## Note
Running this as root is dangerous. run it as nobody, on some port > 1024, then use iptable to redirect traffic from 
port 22 to this app port.
