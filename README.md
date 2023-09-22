# cli

## Installation

### Install via asdf

```shell
asdf plugin add selleo https://github.com/selleo/asdf-cli.git
asdf install selleo latest
asdf global selleo latest

selleo version
```

update:
```shell
asdf install selleo latest
```

### Install via homebrew

If you don't use `asdf`, you can use `brew`.

First time:
```
brew tap Selleo/cli
brew install selleo
```

Upgrade for new release:
```
brew upgrade selleo
```

## Commands

Run `selleo` to see available commands:

```shell
selleo

# example commands
selleo adr new --title "Choose database"
selleo rand uuid4
```

Check `-h` help for command:

```shell
selleo rand bytes -h
selleo crypto hmac sha256 -h
# ...
```

Some 

## About Selleo

![selleo](https://raw.githubusercontent.com/Selleo/selleo-resources/master/public/github_footer.png)

Software development teams with an entrepreneurial sense of ownership at their core delivering great digital products and building culture people want to belong to. We are a community of engaged co-workers passionate about crafting impactful web solutions which transform the way our clients do business.

All names and logos for [Selleo](https://selleo.com/about) are trademark of Selleo Labs Sp. z o.o. (formerly Selleo Sp. z o.o. Sp.k.)

