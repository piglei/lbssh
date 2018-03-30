# lbssh

lbssh is SSH for lazybones, an experimental tool using go-prompt. Main features:

- Auto-complete for your hosts defined in ssh config file
- Using fuzzy search to speed up searching
- Recommend the most possible host for you

<img alt="lbssh-screenshot1 2x" src="https://user-images.githubusercontent.com/731266/38121101-bc18bcc4-33ff-11e8-8611-07c3b614fe9e.gif">

## How to Use

First, you need to install lbssh binary to your system's $BIN path. You could downlod the 
latest binary from [here](https://github.com/piglei/lbssh/releases/tag/v0.0.3).

If you don't have any hosts defined in `~/.ssh/config`, you should try add some:

```
Host us.my-awesome-linode-server-92
HostName 100.101.102.92
Port 22
User piglei
```

For more information on how to write ssh config, you can start from [Simplify Your Life With an SSH Config File](http://nerderati.com/2011/03/17/simplify-your-life-with-an-ssh-config-file/)

Now you are all setup, run `lbssh` and start your lazy life!

## Customize Options

You could change the ssh binary location or ssh config file which lbssh is using, 
below are all supported args:

```console
$ lbssh -h
Usage of lbssh:
      --log-level string         log level (default "INFO")
      --ssh-bin string           ssh binary path (default "/usr/bin/ssh")
      --ssh-config-file string   ssh config file location (default "~/.ssh/config")
      --storage-db-file string   db file location (default "~/.lbssh.db")
      --version                  display version info
```

If you want to use a different SSH_BIN location, you could also try settings up 
`LBSSH_SSH_BIN=/your/ssh/binary` environment variable.

## Future plan

`lbssh` is a simple tool that aims for simplifing the process of finding and logging into 
remote servers. It is best suit for a limited system when tools like "fzf" can not be 
used. I will add more features to make it better.

- Auto-prompt for scp file
- Highlight search result
- Add customize fields
