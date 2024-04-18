# Password Transformation Tool (PTT) Usage Guide
## Version 0.0.0

### Table of Contents
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Usage](#usage)
4. [Examples](#examples)
5. [Contributing](#contributing)
6. [License](#license)

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
TODO
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

#### Binary:
Download the latest release from the GitHub repository:
```
TODO
```

### Usage

The tool can read input from standard input, files, or URLs and can read from
multiple sources at the same time. The tool can also read additional files in a
similar manner for some other options.

The `-f`, `-k`, `-r`, `-tf`, and `-u` flags can be used multiple times and have
their collective values combined. The rest of the flags can only be used once.

#### Options:
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
- `retain`: Transforms input by creating masks that still retain strings from file.
- `match`: Transforms input by keeping only strings with matching masks from a mask file
- `fuzzy-swap`: Transforms input by swapping tokens with fuzzy matches from another
- `swap`: Transforms input by swapping tokens with exact matches from a ':' separated file.
- `pop`: Transforms input by generating tokens from popping strings at character boundaries.
- `shuf`: Transforms input by shuffling tokens from a partial mask file and a input file.

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
- `retain`: `r`, `retain-mask`,
- `match`: `mt`, `match-mask`
- `fuzzy-swap`: `fs`, `fuzzy`, `fuzzy-replace`, `fuzz`, `mutate`
- `swap`: `s`, `replace`
- `pop`: `po`, `split`, `boundary-split`, `boundary-pop`, `pop-split`, `split-pop`
- `shuffle`: `sh`, `shuf`, `mask-swap`, `token-swap`

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

#### Filter Formats:
- `ptt -k keep.txt`: Keep only items in a file.
- `ptt -r remove.txt`: Keep only items not in a file.
- `ptt -k keep.txt -r remove.txt`: Keep only items in a file and not in another.
- `ppt -l 8`: Keep only items equal to a length.
- `ppt -l 8-12`: Keep only items within a range of lengths.
- `ptt -m 10`: Keep only items with a minimum frequency.

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
$ ptt -f rockyou.txt -t pop -l 4-5 -vvv
Starting statistics generation. Please wait...
Verbose Statistics: max=25
--------------------------------------------------
General Stats:
Total Items: 4730675
Total Unique items: 585203
Total Characters: 2758719
Total Words: 585199
Average Characters Per Item: 4
Average Words Per Item: 0
Largest frequency: 29529
Smallest frequency: 1
Mean frequency: 1
Median frequency: 1
Mode frequency: 29529

Category Counts:
alphabetical: 496547
all-lowercase: 524227
short-non-complex: 585203
high-numeric-ratio: 86021
greek-characters: 16
hex-string: 11208
non-complex: 585203
numeric: 86028
non-ASCII: 566
cyrillic-characters: 15
alphanumeric-with-special: 8
starts-uppercase: 60976
all-uppercase: 149650
arabic-characters: 17
thai-characters: 14
hebrew-characters: 3

--------------------------------------------------
1234  [29529]==================================================
2007  [24459]=========================================
2006  [22002]=====================================
love  [21516]====================================
2008  [20022]=================================
ever  [17694]=============================
1994  [14514]========================
life  [14496]========================
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
yahoo [8953]===============
1986  [8860]===============
1985  [8513]==============
```

### Contributing
Contributions are welcome and encouraged. Please open an issue or pull request
if you have any suggestions or improvements. Please follow the code of conduct
when contributing to this project.

### License
This project is licensed under the MIT License - see the LICENSE file for details.

