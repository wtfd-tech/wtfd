# WTFd

[![License](https://img.shields.io/github/license/wtfd-tech/wtfd?style=flat-square)](https://github.com/wtfd-tech/wtfd/blob/master/LICENSE)
[![Latest stable version](https://img.shields.io/github/v/tag/wtfd-tech/wtfd?label=Latest%20Version&style=flat-square)](https://github.com/wtfd-tech/wtfd/releases)  
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fwtfd-tech%2Fwtfd%2Fbadge%3Fref%3Dmaster&style=flat-square)](https://actions-badge.atrox.dev/wtfd-tech/wtfd/goto?ref=master)
[![Codecov](https://img.shields.io/codecov/c/github/wtfd-tech/wtfd?style=flat-square&logo=codecov&label=Coverage)](https://codecov.io/gh/wtfd-tech/wtfd)
[![Dependencies](https://img.shields.io/librariesio/github/wtfd-tech/wtfd?style=flat-square&label=Dependencies)](https://libraries.io/github/wtfd-tech/wtfd)
![Repository Size](https://img.shields.io/github/repo-size/wtfd-tech/wtfd?style=flat-square&label=Repo%20Size)  
[![Last Commit](https://img.shields.io/github/last-commit/wtfd-tech/wtfd?style=flat-square&label=Last%20Commit)](https://github.com/wtfd-tech/wtfd/commits/master)
[![Contributors](https://img.shields.io/github/contributors/wtfd-tech/wtfd?style=flat-square&label=Contributors)](https://github.com/wtfd-tech/wtfd/graphs/contributors)
[![Open Issues](https://img.shields.io/github/issues/wtfd-tech/wtfd?style=flat-square&label=Issues)](https://github.com/wtfd-tech/wtfd/issues)
[![Open PRs](https://img.shields.io/github/issues-pr/wtfd-tech/wtfd?style=flat-square&label=Pull%20Requests)](https://github.com/wtfd-tech/wtfd/pulls)
[![Rawsec's CyberSecurity Inventory](https://inventory.rawsec.ml/img/badges/Rawsec-inventoried-FF5050_flat-square.svg)](https://inventory.rawsec.ml/ctf_platforms.html#WTFd)
<!--Micro badger docker image size-->
<!-- Docker hub stars-->

![](https://raw.githubusercontent.com/wtfd-tech/wtfd/master/icon.svg?sanitize=true)

a [CTFd](https://ctfd.io/)-like Server in go

![demo](https://raw.githubusercontent.com/wtfd-tech/wtfd/master/demo.png)

## Configuration

At start, a `config.yaml` is generated. You should edit it  with the settings you need


The Challenge info Dir shall look like that:

```
├── chall-1
│   ├── meta.yaml
│   ├── README.md
│   └── SOLUTION.md
├── chall-2
│   ├── meta.yaml
│   ├── README.md
│   └── SOLUTION.md
```

For each Challenge you need a `meta.yaml`, a `README.md` and a `SOLUTION.md`

The `meta.yaml` shall look like that:

```
points: <How many points the challenge should have>
uri: "<Protocol and user of your ssh Challenges (e.g. `ssh://chall-1@%s`>"
deps: [<Dependencies the Challenge has>]
flag: "<The flag>"
author: "<The author of the challenge>"
title: "(optional) the title of the challenge, else the directory name is used"
```

The `README.md` and `SOLUTION.md` are markdown files ([syntax](https://github.com/gomarkdown/markdown#extensions)).
The `SOLUTION.md` contents can only be seen by users who already solved the challenge

## Building WTFd yourself

You need to have `go`, `sqlite3` and `yarn` installed

```bash
git clone https://github.com/wtfd-tech/wtfd
cd wtfd
make
```

## Running WTFd

Now you can finally start wtfd by downloading it from the [releases](https://github.com/wtfd-tech/wtfd/releases), giving it permissions `chmod +x wtfd` and running it `./wtfd`

WTFd is HTTP only, if you need HTTPS use a reverse proxy like [Traefik](https://traefik.io/) or [nginx](https://nginx.com/)

## Development notes

To make working with the TypeScript easier, you can do 

```bash
make js-run
```

to automatically compile the JS on changes
