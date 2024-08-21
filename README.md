 Password Transformation Tool (`ptt`) is a versatile utility designed for password cracking. It facilitates the creation of custom rules and transformations, as well as the generation of wordlists. This tool supports processing data from files, URLs, and standard input, streamlining cracking workflows.

## Features:
- Transform input strings with various modes.
- Creates `Hashcat` rules and masks from input strings.
- Transforms input strings with a variety of functions for password cracking.
- Accepts input from multiple sources including; standard input, files, and URLs.
- All transformations support multibyte characters.
- Supports JSON output for easy parsing and integration with other tools.
- Supports multiple transformations and operations with a template file.

## Getting Started:

Documentation on usage and examples can be found in the `/docs/USAGE.md` directory or on the repository here: [link](https://github.com/JakeWnuk/ptt/tree/main/docs/USAGE.md)

### Install:

#### Source:
Fast method with Go installed:
```
go install github.com/jakewnuk/ptt@latest
```
Slow method with Go installed:
```
git clone https://github.com/JakeWnuk/ptt && cd ptt && go build ./main.go && mv ./main ~/go/bin/ptt && ptt
```

#### Docker:
Pull the latest image from Docker Hub:
```
docker run -it -v ${PWD}:/data jwnuk/ptt
``` 
Build the Docker image from the Dockerfile:
```
git clone https://github.com/JakeWnuk/ptt && cd ptt && docker build -t ptt . && docker run -it -v ${PWD}:/data ptt
```

---

```
Usage of Password Transformation Tool (ptt) version (0.3.0):

ptt [options] [...]
Accepts standard input and/or additonal arguments.

The -f, -k, -r, -tf, -tp, and -u flags can be used multiple times, together, and with files or directories.
-------------------------------------------------------------------------------------------------------------
Options:
These modify or filter the transformation mode.

  -b    Bypass map creation and use stdout as primary output.
  -d int
        Enable debug mode with verbosity levels [0-2].
  -f value
        Read additional files for input.
  -i value
        Starting index for transformations if applicable. Accepts ranges separated by '-'.
  -k value
        Only keep items in a file.
  -l value
        Keeps output equal to or within a range of lengths. Accepts ranges separated by '-'.
  -m int
        Minimum numerical frequency to include in output.
  -n int
        Maximum number of items to display in verbose statistics output. (default 25)
  -o string
        Output to JSON file in addition to stdout.
  -r value
        Only keep items not in a file.
  -rm string
        Replacement mask for transformations if applicable. (default "uldsbt")
  -t string
        Transformation to apply to input.
  -tf value
        Read additional files for transformations if applicable.
  -tp value
        Read a template file for multiple transformations and operations.
  -u value
        Read additional URLs for input.
  -v    Show verbose output when possible.
  -vv
        Show statistics output when possible.
  -vvv
        Show verbose statistics output when possible.
  -w int
        Number of words to generate for passphrases if applicable.
-------------------------------------------------------------------------------------------------------------
Transformation Modes:
These create or alter based on the selected mode.

  -t decode
        Transforms input by URL, HTML, and Unicode escape decoding.
  -t dehex
        Transforms input by decoding $HEX[...] formatted strings.
  -t encode
        Transforms input by URL, HTML, and Unicode escape encoding.
  -t hex
        Transforms input by encoding strings into $HEX[...] format.
  -t mask -rm [uldsb] -v
        Transforms input by masking characters with provided mask.
  -t mask-match -tf [file]
        Transforms input by keeping only strings with matching masks from a mask file.
  -t mask-pop -rm [uldsbt]
        Transforms input by generating tokens from popping strings at character boundaries.
  -t mask-remove -rm [uldsb]
        Transforms input by removing characters with provided mask characters.
  -t mask-retain -rm [uldsb] -tf [file]
        Transforms input by creating masks that still retain strings from file.
  -t mask-swap -tf [file]
        Transforms input by swapping tokens from a partial mask file and a input file.
  -t passphrase -w [words] -tf [file]
        Transforms input by randomly generating passphrases with a given number of words and separators from a file.
  -t replace -tf [file]
        Transforms input by replacing all strings with all matches from a ':' separated file.
  -t rule-append
        Transforms input into append rules.
  -t rule-append-remove
        Transforms input into append-remove rules.
  -t rule-insert -i [index]
        Transforms input into insert rules starting at index.
  -t rule-overwrite -i [index]
        Transforms input into overwrite rules starting at index.
  -t rule-prepend
        Transforms input into prepend rules.
  -t rule-prepend-remove
        Transforms input into prepend-remove rules.
  -t rule-toggle -i [index]
        Transforms input into toggle rules starting at index.
  -t substring -i [index]
        Transforms input by extracting substrings starting at index and ending at index.
  -t swap -tf [file]
        Transforms input by swapping tokens with exact matches from a ':' separated file.
-------------------------------------------------------------------------------------------------------------
```

