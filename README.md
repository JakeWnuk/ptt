 `ptt` or Password Transformation Tool is a multi-tool for working with text for password cracking. It can be used to transform text in a variety of ways for password cracking such as:

```
$ ptt -h
Usage of Password Transformation Tool (ptt) version (0.0.0):

ptt [options] [...]
Accepts standard input and/or additonal arguments.

Options:
  -f value
        Read additonal files for input.
  -i value
        Starting index for transformations if applicable. Accepts ranges separated by '-'. (default 0)
  -k value
        Only keep items in a file.
  -m int
        Minimum numerical frequency to include in output.
  -r value
        Only keep items not in a file.
  -rm string
        Replacement mask for transformations if applicable. (default "uldsb")
  -t string
        Transformation to apply to input.
  -tf value
        Read additonal files for transformations if applicable.
  -v    Show verbose output when possible.

The '-f', '-k', '-r', and '-tf' flags can be used multiple times.

Transformation Modes:
  -t append
        Transforms input into append rules.
  -t append-remove
        Transforms input into append-remove rules.
  -t append-shift
        Transforms input into append-shift rules.
  -t prepend
        Transforms input into prepend rules.
  -t prepend-remove
        Transforms input into prepend-remove rules.
  -t prepend-shift
        Transforms input into prepend-shift rules.
  -t insert -i [index]
        Transforms input into insert rules starting at index.
  -t remove -i [index]
        Transforms input into remove rules starting at index.
  -t toggle -i [index]
        Transforms input into toggle rules starting at index.
  -t encode
        Transforms input by  URL, HTML, and Unicode escape encoding.
  -t mask -rm [uldsb] -v
        Transforms input by masking characters with provided mask.
  -t dehex
        Transforms input by decoding $HEX[...] formatted strings.
  -t hex
        Transforms input by encoding strings into $HEX[...] format.
  -t remove -rm [uldsb] -v
        Transforms input by removing characters with provided mask characters.
  -t retain -rm [uldsb] -tf [file]
        Transforms input by creating masks that still retain strings from file.
  -t match -tf [file]
        Transforms input by keeping only strings with matching masks from a mask file.
  -t fuzzy-swap -tf [file]
        Transforms input by swapping tokens with fuzzy matches from another file.
  -t swap -tf [file]
        Transforms input by swapping tokens with exact matches from a ':' separate file.
```

## Getting Started:

>[!NOTE]
> This tool is still in development and considered early access. Please report any issues, bugs, or feature requests to the GitHub repository.

### Install from Go
```
TODO
```

### Install from Source
```
git clone https://github.com/JakeWnuk/ptt && cd ptt && go build ./main.go && mv ./main ~/go/bin/ptt
```

## Docker:
```
TODO
``` 
