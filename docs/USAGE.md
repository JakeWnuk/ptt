# Password Transformation Tool (PTT) Usage Guide
## Version 0.2.5

### Table of Contents
#### Getting Started
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Examples](#examples)
5. [Contributing](#contributing)
6. [License](#license)

#### Mask Creation Guide
1. [Mask Creation Introduction](#mask-creation-introduction)
2. [Mask Creation](#mask-creation)
3. [Mask Matching](#mask-matching)
4. [Removing Characters by Mask](#removing-characters-by-mask)
5. [Creating Retain/Partial Masks](#creating-retainpartial-masks)

### Rule Creation Guide
1. [Rule Creation Introduction](#rule-creation-introduction)
2. [Append Rules](#append-rules)
3. [Prepend Rules](#prepend-rules)
4. [Toggle Rules](#toggle-rules)
5. [Insert Rules](#insert-rules)
6. [Overwrite Rules](#overwrite-rules)

### Wordlist Creation Guide
1. [Wordlist Creation Introduction](#wordlist-creation-introduction)
2. [Direct Swapping](#direct-swapping)
3. [Replacing Text and Characters](#replacing-text-and-characters)
4. [Token Popping](#token-popping)
5. [Token Swapping](#token-swapping)
6. [Passphrases](#passphrases)

### Misc Creation Guide
1. [Misc Creation Introduction](#misc-creation-introduction)
2. [Encoding and Decoding](#encoding-and-decoding)
3. [Hex and Dehex](#hex-and-dehex)
4. [Substrings](#substrings)

## Getting Started

### Introduction
The Password Transformation Tool (PTT) is a command-line utility that allows
users to transform passwords using a variety of methods. This guide will
provide instructions on how to install and use the tool.

The tool was created as a complete solution for password transformation, and
is designed to be easy to use and flexible. PTT is designed around my previous
tools, `maskcat`, `rulecat`, and `mode`, and offers many of the same features
and capabilities with a more user-friendly interface and new functionality.

The tool can read multiple input from standard input, files, or URLs and can
read from multiple sources at the same time. The tool reads all input to
a single data object and then processes the data object with the specified
transformations.

The tool can support multibyte characters in all transformations and does *not*
convert `$HEX[...]` formatted strings to their original characters before use
in a transformation.

The output always contains no duplicates and is sorted by frequency of
occurrence. The output can be shown as is, with frequency counts, as a simple
statistics report, or as a verbose statistics report.

### Installation

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

### Usage

The tool can read input from standard input, files, or URLs and can read from
multiple sources at the same time. The tool can also read additional files in a
similar manner for some other options.

There are some additional notes when importing data:
- Check for hidden characters in files that may cause issues. `Dos2unix` can be used to remove these characters.
- When reading from standard input, the tool can detect chaining `ptt` commands
  when the `-v` flag is used. This can be used to pipe multiple commands together.
- When reading from files, the tool can detect when `ptt` JSON output is used as input and will parse the JSON data.
- The `-b` flag can be used to bypass map creation and use stdout as primary output. This can be useful for working with large datasets. If the `-b` flag is used, the final output will be empty and all filtering and duplication removal will be disabled.
- The `-d [0-2]` flag can be used to enable debug output. This will show the data
  object after all transformations have been applied. There are two (2) levels
  of debug output that can be used.
    - Level 1 will not print each iteration transformation but overall input and output.
    - Level 2 will print each iteration transformation and overall input and output.
- The `-tp` flag can not be used with other transformations at the same time (`-t`). The
  template file should contain a list of transformations and operations to apply
  to the input data. The template file should be in JSON format.
    - See `docs/template.json` ([link](https://github.com/JakeWnuk/ptt/blob/main/docs/template.json)) for an example.
    - See `docs/templates/` ([link](https://github.com/JakeWnuk/ptt/blob/main/templates/)) for more examples.

The `-f`, `-k`, `-r`, `-tf`, `-tp`, and `-u` flags can be used multiple times and have
their collective values combined. The rest of the flags can only be used once.
These flags work with files and directories.

#### Options:
- `-b`: Bypass map creation and use stdout as primary output.
- `-d`: Enable debug mode with verbosity levels [0-2].
- `-f`: Read additional files for input.
- `-i`: Starting index for transformations if applicable. Accepts ranges separated by '-'. (default 0)
- `-k`: Only keep items in a file.
- `-l`: Keeps output equal to or within a range of lengths. Accepts ranges separated by '-'. (default 0)
- `-m`: Minimum numerical frequency to include in output.
- `-n`: Maximum number of items to display in verbose statistics output. (default 25)
- `-o`: Output to JSON file in addition to stdout.
- `-r`: Only keep items not in a file.
- `-rm`: Replacement mask for transformations if applicable. (default "uldsb")
- `-t`: Transformation to apply to input.
- `-tf`: Read additional files for transformations if applicable.
- `-tp`: Read a template file for multiple transformations and operations.
- `-u`: Read additional URLs for input.
- `-v`: Show verbose output when possible.
- `-vv`: Show statistics output when possible.
- `-vvv`: Show verbose statistics output when possible.

#### Transformations:
The following transformations can be used with the `-t` flag:
- `append`: Transforms input into append rules.
- `append-remove`: Transforms input into append-remove rules
- `append-shift`: Transforms input into append-shift rules.
- `prepend`: Transforms input into prepend rules.
- `prepend-remove`: Transforms input into prepend-remove rules.
- `prepend-shift`: Transforms input into prepend-shift rules.
- `insert`: Transforms input into insert rules starting at index.
- `overwrite`: Transforms input into overwrite rules starting at index.
- `toggle`: Transforms input into toggle rules starting at index.
- `encode`: Transforms input by URL, HTML, and Unicode escape encoding.
- `decode`: Transforms input by URL, HTML, and Unicode escape decoding.
- `hex`: Transforms input by encoding strings into $HEX[...] format.
- `dehex`: Transforms input by decoding $HEX[...] formatted
- `mask`: Transforms input by masking characters with provided mask.
- `remove`: Transforms input by removing characters with provided mask characters.
- `substring`: Transforms input by extracting substrings starting at index and ending at index.
- `mask-retain`: Transforms input by creating masks that still retain strings from file.
- `mask-match`: Transforms input by keeping only strings with matching masks from a mask file
- `swap`: Transforms input by swapping tokens with exact matches from a ':' separated file.
- `pop`: Transforms input by generating tokens from popping strings at character boundaries.
- `mask-swap`: Transforms input by swapping tokens from a partial mask file and a input file.
- `passphrase`: Transforms input by randomly generating passphrases with a given number of words and separators from a file.

The modes also have aliases that can be used with the `-t` flag instead of the
keywords above:
- `append`: `a`
- `append-remove`: `ar`
- `append-shift`: `as`
- `prepend`: `p`
- `prepend-remove`: `pr`
- `prepend-shift`: `ps`
- `insert`: `i`
- `overwrite`: `o`
- `toggle`: `t`
- `encode`: `e`
- `decode`: `de`
- `hex`: `h`, `rehex`
- `dehex`: `dh`, `unhex`
- `mask`: `m`, `partial-mask`, `partial`
- `remove`: `rm`, `remove-all`, `delete`, `delete-all`
- `replace`: `rp`, `rep`
- `substring`: `sub`, `sb`
- `retain`: `r`, `retain-mask`,
- `match`: `mt`, `match-mask`
- `swap`: `sw`, `swp`
- `pop`: `po`, `split`, `boundary-split`, `boundary-pop`, `pop-split`, `split-pop`
- `mask-swap`: `ms`, `shuf`, `shuffle`, `token-swap`
- `passphrase`: `pp`, `phrase`

### Examples

#### Input Formats:
- `ptt < input.txt`: Read input from a file.
- `cat input.txt | ptt`: Read input from standard input.
- `ptt -u https://example.com/input.txt`: Read input from a URL.
- `ptt -f input2.txt -f input3.txt -f input4.txt`: Read additional files for input.
- `cat input2.txt | ptt -f input3.txt -u https://example.com/input4.txt`: Read input from standard input and additional files and URLs.

#### Transformation Formats:
- `ptt -t [transformation]`: Apply a transformation to input.
- `ptt -tf file.txt -t [transformation]`: Read file input required for a transformation.
- `ptt -tf file2.txt -tf file3.txt -t [transformation]`: Read multiple files for a transformation.
- `ptt -t [transformation] -rm ulds`: Apply a transformation with a custom mask. Default is all characters.
- `ptt -t [transformation] -i 5`: Apply a transformation starting at a specific index.
- `ptt -i 1-5 -t [transformation]`: Apply a transformation starting at a specific index.
- `ptt -tp template.json`: Apply multiple transformations and operations from a template file.

#### Filter Formats:
- `ptt -k keep.txt`: Keep only items in a file.
- `ptt -r remove.txt`: Keep only items not in a file.
- `ptt -k keep.txt -r remove.txt`: Keep only items in a file and not in another.
- `ppt -l 8`: Keep only items equal to a length.
- `ppt -l 8-12`: Keep only items within a range of lengths.
- `ptt -m 10`: Keep only items with a minimum frequency.

#### Debug Formats:
- `ptt -d 1`: Enable debug mode with verbosity level 1.
- `ptt -d 2`: Enable debug mode with verbosity level 2.

#### Output Formats:
- `ptt -v`: Show verbose output.
- `ptt -vv`: Show statistics output.
- `ptt -vvv`: Show verbose statistics output.
- `ptt -n 50`: Show verbose statistics output with a maximum of 50 items.
- `ptt -o [FILE]`: Show output and save JSON output to a file.
- These options are available for all transformations.

#### Real Examples:
- `ptt -f rockyou.txt -t pop -l 4-5`:
```shell
$ ptt -f rockyou.txt -t pop -l 4-5
1234
2007
2006
love
2008
ever
1994
life
2005
1992
...
```
- `ptt -f rockyou.txt -t pop -l 4-5 -v`:
```shell
$ ptt -f rockyou.txt -t pop -l 4-5 -v
29529 1234
24459 2007
22002 2006
21516 love
20022 2008
17694 ever
14514 1994
14496 life
14300 2005
14159 1992
...
```
- `ptt -f rockyou.txt -t pop -l 4-5 -vv`:
```shell
$ ptt -f rockyou.txt -t pop -l 4-5 -vv
1234 [29529]==================================================
2007 [24459]=========================================
2006 [22002]=====================================
love [21516]====================================
2008 [20022]=================================
ever [17694]=============================
1994 [14514]========================
life [14496]========================
2005 [14300]========================
1992 [14159]=======================
```
- `ptt -f rockyou.txt -t pop -l 4-5 -vvv`:
```shell
[*] Starting statistics generation. Please wait...
Verbose Statistics: max=25
--------------------------------------------------
General Stats:
Total Items: 4695779
Total Unique items: 613210
Total Words: 613206
Largest frequency: 29529
Smallest frequency: 1

Plots:
Item Length: |[|==========]|
Min: 4, Q1: 4, Q2: 4, Q3: 5, Max: 5
Item Frequency: |[|]--------------------------------------------------|
Min: 1, Q1: 1, Q2: 1, Q3: 3, Max: 29529
Item Complexity: |[|]----------------------------------|
Min: 1, Q1: 1, Q2: 1, Q3: 1, Max: 3

Category Counts:
all-uppercase: 58433
non-ASCII: 547
alphanumeric-with-special: 8
alphabetical: 524554
short-non-complex: 613210
numeric: 86028
all-lowercase: 410494
non-complex: 613210
hebrew-characters: 3
hex-string: 11395
thai-characters: 14
arabic-characters: 17
cyrillic-characters: 13
starts-uppercase: 114042
high-numeric-ratio: 86014
greek-characters: 16

--------------------------------------------------
1234  [29529]==================================================
2007  [24459]=========================================
2006  [22002]=====================================
love  [20435]==================================
2008  [20022]=================================
ever  [17605]=============================
1994  [14514]========================
life  [14460]========================
2005  [14300]========================
1992  [14159]=======================
1993  [14070]=======================
12345 [13545]======================
1991  [13038]======================
1995  [12932]=====================
1990  [12336]====================
1989  [11355]===================
1987  [10903]==================
1996  [9929]================
2000  [9801]================
1988  [9718]================
2009  [9257]===============
2004  [9091]===============
yahoo [8942]===============
1986  [8860]===============
1985  [8513]==============
```

### Contributing
Contributions are welcome and encouraged. Please open an issue or pull request
if you have any suggestions or improvements. Please follow the code of conduct
when contributing to this project.

### License
This project is licensed under the MIT License - see the LICENSE file for details.

## Mask Creation Guide

### Mask Creation Introduction
This document describes the ways to use PTT to create `Hashcat` compatible
masks. There are several ways to use masks in PTT:

- `Mask Creation`: Create a mask from a given string.
- `Mask Matching`: Match a mask to a given string.
- `Removing Characters by Mask`: Remove characters from a given string by a mask.
- `Creating Retain/Partial Masks`: Create a mask that retains only certain keywords.

All modes support multibyte characters and can properly convert them. One
transformation can be used at a time.

> [!CAUTION]
> Ensure input is provided in the correct format and does not contain hidden characters. `Dos2Unix` can be used to convert the file to proper format if needed.

### Mask Creation
Masks replace characters in a string with a common character. The syntax to
create a mask is as follows:
```
ptt -f <input_file> -t mask -rm <mask_characters> -v
```
Where `<mask_characters>` can be any of the following:
- `u`: Uppercase characters
- `l`: Lowercase characters
- `d`: Digits
- `s`: Special characters
- `b`: Byte characters
- Multiple characters can be combined to create a mask.

The default value is `uldsb` for all characters. The `-v` flag is optional and
if provided, will print the length of the original string and its character
complexity. The format will be `:length:complexity` appended to the end of the
output.

### Mask Matching
Masks can be matched to a given string to determine if the string matches the
mask. The syntax to match a mask is as follows:
```
ptt -f <input_file> -t match -tf <mask_file>
```
Where `<mask_file>` is the file containing the mask to match. The mask file
should only contain valid masks. The output will be all of the strings that
match the masks.

### Removing Characters by Mask
Characters can be removed from a string by a mask. The syntax to remove
characters by mask is as follows:
```
ptt -f <input_file> -t remove -rm <mask_characters>
```
Where `<mask_characters>` is the mask to remove from the string. The output will
be the string with the characters removed.

### Creating Retain/Partial Masks
Retain masks can be created to retain only certain keywords in a string. The
syntax to create a retain mask is as follows:
```
ptt -f <input_file> -t retain -rm <mask_characters> -tf <keep_file>
```
Where `<mask_characters>` is the mask to retain and `<keep_file>` is the file
containing the keywords to retain. The output will be the mask with only the
keywords retained.

## Rules Creation Guide

### Rule Creation Introduction
This document describes the ways to use PTT to create `Hashcat` compatible
rules. There are several types of rules that can be created using PTT:

- `Append Rules`: Append a string to the end of the password.
- `Append Remove Rules`: Remove characters from the end of the password before appending a string.
- `Append Shift Rules`: Shift the characters of the password to the right before appending a string.
- `Prepend Rules`: Prepend a string to the beginning of the password.
- `Prepend Remove Rules`: Remove characters from the beginning of the password before prepending a string.
- `Prepend Shift Rules`: Shift the characters of the password to the left before prepending a string.
- `Toggle Rules`: Toggle the case of the password.
- `Insert Rules`: Insert a string at a specific position in the password.
- `Overwrite Rules`: Overwrite a string at a specific position in the password.

All modes support multibyte characters and can properly convert them. One
transformation can be used at a time.

> [!CAUTION]
> Ensure input is provided in the correct format and does not contain hidden characters. `Dos2Unix` can be used to convert the file to proper format if needed.

### Append Rules
Append rules are used to append a string to the end of the password. The syntax for an append rule is as follows:
```
ptt -f <input_file> -t append
```

The append mode also has two additional options:
- `append-remove`: Remove characters from the end of the password before appending a string.
- `append-shift`: Shift the characters of the password to the right before appending a string.

The syntax for an append-remove rule is as follows:
```
ptt -f <input_file> -t append-remove
```

The syntax for an append-shift rule is as follows:
```
ptt -f <input_file> -t append-shift
```

### Prepend Rules
Prepend rules are used to prepend a string to the beginning of the password. The syntax for a prepend rule is as follows:
```
ptt -f <input_file> -t prepend
```

The prepend mode also has two additional options:
- `prepend-remove`: Remove characters from the beginning of the password before prepending a string.
- `prepend-shift`: Shift the characters of the password to the left before prepending a string.

The syntax for a prepend-remove rule is as follows:
```
ptt -f <input_file> -t prepend-remove
```

The syntax for a prepend-shift rule is as follows:
```
ptt -f <input_file> -t prepend-shift
```

### Toggle Rules
Toggle rules are used to toggle the case of the password. The syntax for a toggle rule is as follows:
```
ptt -f <input_file> -t toggle -i <index>
```
Where `<index>` is the starting index of the toggle pattern. If no index is provided,
the toggle pattern will start at the beginning of the password. The `<index>`
can also accept range values in the format of `start-end`. For example, `1-5` will
print output for the toggle transformation starting from index 1 to 5.

### Insert Rules
Insert rules are used to insert a string at a specific position in the password. The syntax for an insert rule is as follows:
```
ptt -f <input_file> -t insert -i <index>
```
Where `<index>` is the position where the string will be inserted. If no index is provided,
the string will be inserted at the beginning of the password. The `<index>`
can also accept range values in the format of `start-end`. For example, `1-5` will
print output for the insert transformation starting from index 1 to 5.

### Overwrite Rules
Overwrite rules are used to overwrite a string at a specific position in the password. The syntax for an overwrite rule is as follows:
```
ptt -f <input_file> -t overwrite -i <index>
```
Where `<index>` is the position where the string will be overwritten. If no index is provided,
the string will be overwritten at the beginning of the password. The `<index>`
can also accept range values in the format of `start-end`. For example, `1-5` will
print output for the overwrite transformation starting from index 1 to 5.

## Wordlist Creation Guide

### Wordlist Creation Introduction
This document describes the ways to use PTT to create password cracking
wordlists. There are several ways to generate wordlists using PTT:

- `Direct Swapping`: Swapping characters directly with a `:` separated file.
   This is implemented in the `swap` module.
- `Replacing Text and Characters`: Replacing text and characters in a string.
  This is implemented in the `replace` module
- `Token Popping`: Generates tokens by popping strings at character boundaries.
  This is implemented in the `pop` module.
- `Token Swapping`: Generates tokens by swapping characters in a string. This is
  implemented in the `mask-swap` module.
- `Passphrases`: Generates passphrases by combining words from a wordlist. This
  is implemented in the `passphrase` module.

All modes support multibyte characters and can properly convert them. One
transformation can be used at a time.

> [!CAUTION]
> Ensure input is provided in the correct format and does not contain hidden characters. `Dos2Unix` can be used to convert the file to proper format if needed.

### Direct Swapping
The `swap` module swaps characters directly with a `:` separated file. The
syntax is as follows:
```
ptt -f <input-file> -t swap -tf <replacement-file>
```
The replacement file should contain the strings to be transformed as `PRIOR:POST`
pairs. The replacements will be applied to the all instance in each line but
only one swap is applied at once. This mode is ideal for subsituting words or characters in a string.

### Replacing Text and Characters
The `replace` module replaces text and characters in a string. This mode replaces all strings with all matches from a ':' separated file. The syntax is as follows:
```
ptt -f <input-file> -t replace -tf <replacement-file>
```
The replacement file should contain the strings to be transformed as
`PRIOR:POST` pairs. The replacements will be applied to all instances in each
line and all replacements are applied to the string. This mode is ideal for replacing all instances of a word or character in
a string.

### Token Popping
The `pop` module generates tokens by popping strings at character boundaries.
The syntax is as follows:
```
ptt -f <input-file> -t pop -rm <mask-characters>
```
Where `<mask_characters>` can be any of the following:
- `u`: Uppercase characters
- `l`: Lowercase characters
- `d`: Digits
- `s`: Special characters
- `b`: Byte characters
- `t`: Title case words (requires `u` and `l`)
- Multiple characters can be combined to create a mask.

The default value is `uldsbt` for all characters. This mode will create tokens
by popping characters from the input string then aggregating the results.

### Token Swapping
The `mask-swap` module generates tokens by swapping characters in a string. The
syntax is as follows:
```
ptt -f <input-file> -t mask-swap -tf <replacement-file>
```
> [!NOTE]
> The input for `mask-swap` is partial masks from `retain`! This is different from most other modes.

The replacement file does not need to be in any specific format. The
replacements will be applied to the first instance in each line. The
`mask-swap` mode is unique in that it uses partial masks from the `retain`
module to generate new candidates. This mode also uses its own replacer
module (different from the other modes) to generate new candidates by
extracting the masks and then matching them to the replacement file.

This mode is most similar to token-swapping in that it generates new
candidates by using masks. However, it is unique in that it uses partial
masks to limit the swap positions from prior applications.

### Passphrases
The `passphrase` module generates passphrases by combining words from a wordlist.
The `-w` flag can be used to specify the number of words to use in the passphrase.
The `-tf` flag is optional and can be used to specify a file containing separators
to use between words. The syntax is as follows:
```
ptt -f <input-file> -t passphrase -w <word-count> -tf <separator-file>
```

The passphrases are generated randomly by selecting words and separators from the input.
If no separator file is provided, no separators will be used. The default word count is 0.
The number of passphrases generated is equal to the number of lines in the input file
*including* duplicates. This means that the item count is also used to determine the number
of passphrases generated.

## Misc Creation Guide

### Misc Creation Introduction
This document describes the ways to use PTT to create miscellaneous transformations.
There are several types that can be created using PTT:

- `Encoding and Decoding`: This transforms input to and from URL, HTML, and Unicode escaped strings.
- `Hex and Dehex`: This transforms input to and from `$HEX[....]` strings.
- `Substrings`: This extracts substrings from the input based on position.

All modes support multibyte characters and can properly convert them. One
transformation can be used at a time.

> [!CAUTION]
> Ensure input is provided in the correct format and does not contain hidden characters. `Dos2Unix` can be used to convert the file to proper format if needed.

### Encoding and Decoding
This mode allows encoding and decoding of input to and from URL, HTML, and Unicode escaped strings.
The syntax is as follows:
```
ptt -f <input_file> -t encode
```
or
```
ptt -f <input_file> -t decode
```
The following table shows the supported transformations:

| Transformation | Description | Input Example | Output Example |
| --- | --- | --- | --- |
| `url` | URL encoding | `https://www.example.com` | `https%3A%2F%2Fwww.example.com` |
| `html` | HTML encoding | `<html>` | `&lt;html&gt;` |
| `unicode` | Unicode encoding | `Hello😎` | `Hello\u1f60e` |

### Hex and Dehex
This mode allows encoding and decoding of input to and from `$HEX[....]` strings.
The syntax is as follows:
```
ptt -f <input_file> -t hex
```
or
```
ptt -f <input_file> -t dehex
```
The following table shows the supported transformations:

| Transformation | Description | Input Example | Output Example |
| --- | --- | --- | --- |
| `hex` | Hex encoding | `Hello` | `$HEX[48656c6c6f]` |
| `dehex` | Hex decoding | `$HEX[48656c6c6f]` | `Hello` |

### Substrings
This mode allows extracting substrings from the input based on position. The syntax is as follows:
```
ptt -f <input_file> -t substring -i <start_index>
```
or
```
ptt -f <input_file> -t substring -i <start_index>-<end_index>
```

This transformation extracts the substring from the input based on the provided
index. If the end index is greater than the length of the input, it will be
changed to the length of the input.

This transformation can be used to extract specific parts of the input for
further processing.

