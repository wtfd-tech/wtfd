# WTFd

![Codecov](https://img.shields.io/codecov/c/github/wtfd-tech/wtfd?style=for-the-badge)

![icon]()

a [CTFd](https://ctfd.io/)-like Server in go

![demo](https://raw.githubusercontent.com/wtfd-tech/wtfd/master/demo.png)

## Configuration

You need a `config.json` in the same path as your wtfd binary (or the path you're in if you do `go run ./cmd/wtfd.go`)

It shall look like that:

```
{
        "Port": <The Port WTFd should run (ports <1024 need root (don't run WTFd as root))>,
        "challinfodir": "<The (relative or absolute) directory with your challenge infos>",
       	"sshhost": "<The domain of your ssh challenges>"
}
```

WTFd will also generate the field `Key` in which the cookie session key will be stored


The Challenge info Dir shall look like that:

```
├── chall-1
│   ├── meta.json
│   ├── README.md
│   └── SOLUTION.md
├── chall-2
│   ├── meta.json
│   ├── README.md
│   └── SOLUTION.md
```

For each Challenge you need a `meta.json`, a `README.md` and a `SOLUTION.md`

The `meta.json` shall look like that:

```
{
	"points": <How many points the challenge should have>,
        "uri": "<Protocol and user of your ssh Challenges (e.g. `ssh://chall-1@%s`>",
	"deps": [<Dependencies the Challenge has>],
	"flag": "<The flag>",
	"author": "<The author of the challenge>"
}
```

The `README.md` and `SOLUTION.md` are markdown files ([syntax](https://github.com/gomarkdown/markdown#extensions)).
The `SOLUTION.md` contents can only be seen by users who already solved the challenge

Now you can finally start wtfd by downloading it from the [releases](https://github.com/wtfd-tech/wtfd/releases), giving it permissions `chmod +x wtfd` and running it `./wtfd`
