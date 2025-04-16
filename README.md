 Password Transformation Tool (`ptt`) is a versatile utility designed for password cracking. It facilitates the creation of custom rules and transformations, as well as the generation of wordlists. This tool supports processing data from files, URLs, and standard input, streamlining cracking workflows.

`PTT` is written in `Go`, is compatible with multiple platforms, and can be easily integrated into existing workflows. The tool is designed to be user-friendly and intuitive, with a wide range of features and options.

## Features:
- **Transformation Modes:** Choose from various transformation modes to
  manipulate input data.
- **Wordlist Generation:** Generate wordlists from input data using custom rules
  and transformations.
- **Rule Creation:** Create custom rules for appending, prepending,
  overwriting, and toggling strings.
- **Hashcat Rule Simplification:** Apply rules to input data and simplify
  rules for optimization by using the [HCRE](https://git.launchpad.net/hcre/tree/README.md) library.

## Getting Started:

Usage documentation can be found in the `/docs/USAGE.md` directory or on the repository here:
- [GitHub Link](https://github.com/JakeWnuk/ptt/tree/main/docs/USAGE.md)

### Install:

From source with `go`:
```
go install github.com/jakewnuk/ptt@latest
```
From `git` clone then build with `go`:
```
git clone https://github.com/JakeWnuk/ptt && cd ptt && go build ./main.go && mv ./main ~/go/bin/ptt && ptt
```

```
Usage of Password Transformation Tool (ptt) version (1.1.0):

ptt [options] [...]
Accepts standard input and/or additonal arguments.

The -t flag can be used multiple times.

Options:
These modify or filter the transformation mode.

  -d    Enable debug mode.
  -i value
        Starting index for transformations if applicable. Accepts ranges separated by '-'.
  -l value
        Only output items of a certain length. Accepts ranges separated by '-'.
  -m string
        Mask for transformations if applicable. (default "uldsbt")
  -t value
        Transformation mode to be used. Can be specified multiple times.
  -v    Show verbose report output. Warning: loads information into memory.
  -w value
        Number of words for transformations if applicable. Accepts ranges separated by '-'.

Transformation Modes:
These create or alter based on the selected mode.

  -t mask -m [uldsb]
        Transforms input by masking characters with provided mask.
  -t mask-pop -m [uldsbt]
        Transforms input by popping tokens from character boundaries using the provided mask.
  -t mask-remove -m [uldsb]
        Transforms input by removing characters with provided mask.
  -t passphrase -w [words]
        Transforms input by generating passphrases from sentences with a given number of words.
  -t regram -w [words]
        Transforms input by regramming sentences into new n-grams with a given number of words.
  -t rule-append                                                                                                                                                                                                                                                                          Transforms input by creating append rules.
  -t rule-append-remove
        Transforms input by creating append-remove rules.
  -t rule-insert -i [index]
        Transforms input by creating insert rules starting at index.
  -t rule-overwrite -i [index]
        Transforms input by creating overwrite rules starting at index.
  -t rule-prepend
        Transforms input by creating prepend rules.
  -t rule-prepend-remove
        Transforms input by creating prepend-remove rules.
  -t rule-prepend-toggle
        Transforms input by creating prepend-toggle rules.
  -t rule-simplify
        Transforms input by simplifying rules to efficient equivalents using the HCRE library.
  -t rule-toggle -i [index]
        Transforms input by creating toggle rules starting at index.
  -t token-swap
        Transforms input by swapping tokens.
```
