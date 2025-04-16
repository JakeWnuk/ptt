# Password Transformation Tool (PTT) Usage Guide
> Version 1.1.0
## Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)

## Introduction
The Password Transformation Tool (PTT) is a command-line utility that allows users to transform passwords using various methods. This guide will provide instructions on how to install and use the tool.

The tool was created as a complete solution for password transformation and is designed to be easy to use and flexible. PTT is designed around my previous tools, `maskcat`, `rulecat`, and `mode`, and offers many of the same features and capabilities with a more user-friendly interface and new functionality.

The tool can read multiple inputs from standard input, files, or URLs and can read from multiple sources at the same time. The tool reads all input into a single data object and then processes the data object with the specified transformations.

The output contains no duplicates and is sorted by frequency of occurrence. The output can be shown as is, with frequency counts, as a simple statistics report, or as a verbose statistics report. The tool also supports template files, loading directories and files, chaining input from multiple sessions, JSON output, debugging levels, and other quality of life features.

### Installation
From source with `go`:
```
go install github.com/jakewnuk/ptt@latest
```
From `git` clone then build with `go`:
```
git clone https://github.com/JakeWnuk/ptt && cd ptt && go build ./main.go && mv ./main ~/go/bin/ptt && ptt
```
