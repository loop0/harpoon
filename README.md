## Harpoon

[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]
[![Go Report Card](https://goreportcard.com/badge/github.com/agrison/harpoon)][goreportcard]
[![Build Status](https://travis-ci.org/agrison/harpoon.svg?branch=master)](https://travis-ci.org/agrison/harpoon)

[godocs]: https://godoc.org/github.com/agrison/harpoon
[goreportcard]: https://goreportcard.com/report/github.com/agrison/harpoon

Harpoon is a simple Go program that listens for GitHub hooks like `push` on specific git repository `refs` and execute
a defined command when it occurs.

It can be configured in a TOML configuration file.

That's nothing fancy, just didn't want to install the whole nodejs + npm or Apache + PHP just for this.

### Install

```sh
go get github.com/agrison/harpoon
go build
```

Add it to your path if you want.

### Runnning

#### Configuration

You must have a `config.toml` file located in your directory.

This is a sample TOML file:

```toml
port = 9001
addr = "0.0.0.0"

[events."push:foo/bar:refs/heads/develop"]
cmd = "echo"
args = "Push!"

[events."watch:foo/bar:refs/heads/develop"]
cmd = "echo"
args = "Watch!"
```

As you can see events and refs can be configured in the `events` TOML table section.

These keys have the following format `events.{event}:{repository}:{ref}` where `{event}` refers to a GitHub WebHook event
like `push`, `watch`, `pull_request`, ...; `{repository}` is the GitHub repository name and `{ref}` refers to a ref in your git repository like `refs/heads/master` or `refs/tags/v0.0.1`.

For the above example, it will echo `Push!` when something has been pushed to the `develop` branch of your GitHub repository. 
It will also echo `Watch!` if someone stars it.

#### Really running

```sh
harpoon
```

or to have it verbose:
```sh
harpoon -v
```

It will output something like this when running:
```
    __
   / /_  ____ __________  ____  ____  ____
  / __ \/ __ `/ ___/ __ \/ __ \/ __ \/ __ \
 / / / / /_/ / /  / /_/ / /_/ / /_/ / / / /
/_/ /_/\__,_/_/  / .___/\____/\____/_/ /_/
                /_/
	Listening on 0.0.0.0:9001
    
push detected on foo/bar with ref refs/heads/develop with the following commits:
	2016-03-08 13:59:38 +0100 CET - adding some awesomeness by Foo Bar
	2016-03-08 13:59:57 +0100 CET - forgot the unicorn png by Foo Bar
	2016-03-08 14:01:12 +0100 CET - so much joy in one commit by Foo Bar
> Push!
```

Obviously, the command should be something like `/path/to/pull-build-n-deploy.sh` where some awesome stuff is going on. 
Like something involving git, maven, rake, npm, bower, I don't know, this is your job :)

### Security

`POST` requests made by GitHub are validated against the `X-Hub-Signature` if the environment variable `GITHUB_HOOK_SECRET_TOKEN` is set.
It must be the same as the one defined on the GitHub's WebHook page.

### Thanks

- [https://mholt.github.io/json-to-go/](https://mholt.github.io/json-to-go/)
- [https://github.com/BurntSushi/toml](https://github.com/BurntSushi/toml)
- [https://github.com/fatih/color](https://github.com/fatih/color)
- [https://github.com/gorilla/mux](https://github.com/gorilla/mux)

### Licence
MIT
