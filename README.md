# lbssh

lbssh is SSH for lazybones, an experimental tool using go-prompt. Main features:

- Auto-complete for your hosts defined in ssh config file
- [MAYBE] Quick add/edit/remove host entries
- [MAYBE] Manage hosts with groups/tags
- [MAYBE] Even more fancy interfaces

<img width="546" alt="lbssh-screenshot1 2x" src="https://user-images.githubusercontent.com/731266/32082491-79800690-baee-11e7-945a-de4e33578ac7.png">

## Install

Download the latest release in [release page](https://github.com/piglei/lbssh/releases/tag/v0.0.1).

## Use

- Define some Host in your ~/.ssh/config
- run lbssh enter lbssh prompt mode
- enter "go HOSTNAME" to connect to your host
