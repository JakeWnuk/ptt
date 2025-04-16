# Password Transformation Tool (PTT) Usage Guide
> Version 1.1.0
## Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)

## Introduction
The Password Transformation Tool (PTT) is a command-line utility that allows users to transform passwords using various methods. This guide will provide instructions on how to install and use the tool.

The tool was created as a complete solution for password transformation and is designed to be easy to use and flexible. PTT is designed around CLI tools and offers many features and capabilities with a user-friendly interface.

> The `rule-simplify` transformation feature is enabled by the work done on the [HCRE](https://git.launchpad.net/hcre/tree/README.md) project.

### Installation
From source with `go`:
```
go install github.com/jakewnuk/ptt@latest
```
From `git` clone then build with `go`:
```
git clone https://github.com/JakeWnuk/ptt && cd ptt && go build ./main.go && mv ./main ~/go/bin/ptt && ptt
```
