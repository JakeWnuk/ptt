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

Documentation on usage and examples can be found in the `/docs/USAGE.md` directory or on the repository here:
- [GitHub Link](https://github.com/JakeWnuk/ptt/tree/main/docs/USAGE.md)

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

### Usage:
```
Usage of Password Transformation Tool (ptt) version (0.3.6):

ptt [options] [...]
Accepts standard input and/or additonal arguments.

The -f, -k, -r, -tf, -tp, and -u flags can be used multiple times, together, and with files or directories.
-------------------------------------------------------------------------------------------------------------
Options:
These modify or filter the transformation mode.

  -b    Bypass map creation and use stdout as primary output. Disables some options.
  -d int
        Enable debug mode with verbosity levels [0-2].
  -f value
        Read additional files for input.
  -i value
        Starting index for transformations if applicable. Accepts ranges separated by '-'.
  -k value
        Only keep items in a file.
  -l value
        Only output items of a certain length (does not adjust for rules). Accepts ranges separated by '-'.
  -m int
        Minimum numerical frequency to include in output.
  -md
        If Markdown format should be used for output instead.
  -n int
        Maximum number of items to return in output.
  -o string
        Output to JSON file in addition to stdout. Accepts file names and paths.
  -p int
        Change parsing mode for URL input. [0 = Strict, 1 = Permissive, 2 = Maximum] [0-2].
  -r value
        Only keep items not in a file.
  -rm string
        Replacement mask for transformations if applicable. (default "uldsbt")
  -t string
        Transformation to apply to input.
  -tf value
        Read additional files for transformations if applicable.
  -tp value
        Read a template file for multiple transformations and operations. Cannot be used with -t flag.
  -u value
        Read additional URLs for input.
  -v    Show verbose output when possible. (Can show additional metadata in some modes.)
  -vv
        Show statistics output when possible.
  -vvv
        Show verbose statistics output when possible.
  -w value
        Number of words for transformations if applicable. Accepts ranges separated by '-'.
-------------------------------------------------------------------------------------------------------------
Transformation Modes:
These create or alter based on the selected mode.

  -t decode
        Transforms input by HTML and Unicode escape decoding.
  -t dehex
        Transforms input by decoding $HEX[...] formatted strings.
  -t encode
        Transforms input by HTML and Unicode escape encoding.
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
  -t passphrase -w [words]
        Transforms input by generating passphrases from sentences with a given number of words.
  -t regram -w [words]
        Transforms input by 'regramming' sentences into new n-grams with a given number of words.
  -t replace-all -tf [file]
        Transforms input by replacing all strings with all matches from a ':' separated file.
  -t rule-append
        Transforms input by creating append rules.
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
  -t rule-toggle -i [index]
        Transforms input by creating toggle rules starting at index.
  -t substring -i [index]
        Transforms input by extracting substrings starting at index and ending at index.
  -t swap-single -tf [file]
        Transforms input by swapping tokens once per string per replacement with exact matches from a ':' separated file.
-------------------------------------------------------------------------------------------------------------
```
