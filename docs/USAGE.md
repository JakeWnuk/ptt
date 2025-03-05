# Password Transformation Tool (PTT) Usage Guide
> Version 1.0.0
## Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Starter Examples](#starter-examples)
- [Mask Transformation Usage](#mask-transformation-usage)
  - [Mask Creation](#mask-creation)
  - [Mask Matching](#mask-matching)
  - [Removing Characters by Mask](#removing-characters-by-mask)
  - [Creating Retain/Partial Masks](#creating-retainpartial-masks)
- [Rule Transformation Usage](#rule-transformation-usage)
  - [Append Rules](#append-rules)
  - [Prepend Rules](#prepend-rules)
  - [Toggle Rules](#toggle-rules)
  - [Insert Rules](#insert-rules)
  - [Overwrite Rules](#overwrite-rules)
- [Wordlist Creation Usage](#wordlist-creation-usage)
  - [Direct Swapping](#direct-swapping)
  - [Replacing Text and Characters](#replacing-text-and-characters)
  - [Token Popping](#token-popping)
  - [Token Swapping](#token-swapping)
  - [Passphrases](#passphrases)
- [Misc. Transformation Usage](#misc-transformation-usage)
  - [Encoding and Decoding](#encoding-and-decoding)
  - [Hex and Dehex](#hex-and-dehex)
  - [Substrings](#substrings)
  - [Regram](#regram)
  - [Rule Application](#rule-application)
  - [Rule Simplification](#rule-simplification)

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
From `docker` with the program as the entry point:
```
docker run -it -v ${PWD}:/data jwnuk/ptt
``` 
From `git` then build with `docker`:
```
git clone https://github.com/JakeWnuk/ptt && cd ptt && docker build -t ptt . && docker run -it -v ${PWD}:/data ptt
```
### Usage
There are some additional notes when importing data and getting started:
- There are no positional arguments, and every argument is defined after a `flag`.
- When reading from standard input, the tool can detect chaining `ptt` commands when the `-v` flag is used. This can be used to pipe multiple commands together without losing frequency data.
- When reading from files, the tool can detect when `ptt` JSON output is used as input and will parse the JSON data.
- The tool should support multibyte characters and transformations in every mode.
- The `-b` flag can be used to bypass map creation and use `stdout` as the primary output. This can be useful for working with large datasets.
    - If the `-b` flag is used, the final output will be empty, and all filtering and duplication removal will be disabled.
- The `-d [0-2]` flag can be used to enable debug output. This will show the data
  object after all transformations have been applied. There are two levels
  of debug output that can be used.
    - Level 1 will not print each iteration transformation but the overall input and output.
    - Level 2 will print each iteration transformation and the overall input and output.
- The `-tp` flag cannot be used with other transformations at the same time (`-t`).
- The template file should contain a list of transformations and operations to apply
  to the input data. The template file should be in JSON format.
    - See `docs/template.json` ([link](https://github.com/JakeWnuk/ptt/blob/main/docs/template.json)) for an example.
    - See `templates/` ([link](https://github.com/JakeWnuk/ptt/blob/main/templates/)) for more examples.
- The `-f`, `-k`, `-r`, `-tf`, `-tp`, and `-u` flags can be used multiple times and have their collective values combined. The rest of the flags can only be used once. These flags work with files and directories.
- The `-p` flag can be used to change the parsing mode for URLs. The default mode is `0` and will use a narrow character set to parse text from URLs. The `1` mode will use a larger character set to parse text from URLs and include additional parsing by default. The `2` mode will use the same character set as `1` but will also include additional parsing options for maximum parsing, including n-grams and other parsing options.
- The `-i` and `-w` flags can also accept range values in the format of `start-end`. For example, `1-5` will print output for the transformation starting from index 1 to 5. For the `-w` flag, this will be the number of words the output will contain.

> [!CAUTION]
> Ensure input is provided in the correct format and does not contain hidden characters. `Dos2Unix` can be used to convert the file to proper format if needed.
### Starter Examples
#### Input Formats:
- `ptt < input.txt`: Read input from a file.
- `cat input.txt | ptt`: Read input from standard input.
- `ptt -u https://example.com/input.txt`: Read input from a URL.
- `ptt -f input2.txt -f input3.txt -f input4.txt`: Read additional files for input.
- `cat input2.txt | ptt -f input3.txt -u urls.txt`: Read input from standard input and additional files and URLs.
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
- `ppt -l 8`: Only allow items equal to a length for input.
- `ppt -l 8-12`: Keep only items within a range of lengths for input.
- `ptt -m 10`: Keep only items with a minimum frequency from output.
#### Debug Formats:
- `ptt -d 1`: Enable debug mode with verbosity level 1.
- `ptt -d 2`: Enable debug mode with verbosity level 2.
#### Output Formats:
- `ptt -v`: Show verbose output.
- `ptt -vv`: Show statistics output.
- `ptt -vvv`: Show verbose statistics output.
- `ptt -n 50`: Show verbose statistics output with a maximum of 50 items.
- `ptt -o [FILE]`: Show output and save JSON output to a file.
- `ptt -md`: Show output as a Markdown table.
- `ptt -ic`: Ignore case when creating output and convert to lowercase.
- These options are available for all transformations.
#### Rockyou Examples:
`ptt -f rockyou.txt -t pop -l 4-5`:

**Flags**: `-f` to select a *file*, `-t` to use the `pop` *transformation*, and `-l` to give a  *length* range of results.

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

 `ptt -f rockyou.txt -t pop -l 4-5 -v`:
 
 **Flags**: `-f` to select a *file*, `-t` to use the `pop` *transformation*, `-l` to give a  *length* range of results, and `-v` to print *verbose* output.

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

`ptt -f rockyou.txt -t pop -l 4-5 -vv`:
 
 **Flags**: `-f` to select a *file*, `-t` to use the `pop` *transformation*, `-l` to give a  *length* range of results, and `-vv` to print very *verbose* output with an item graph.

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
`ptt -f rockyou.txt -t pop -l 4-5 -vvv`:

**Flags**: `-f` to select a *file*, `-t` to use the `pop` *transformation*, `-l` to give a  *length* range of results, and `-vvv` to print very very *verbose* output with a full report.

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
## Mask Transformation Usage
There are several ways to use masks in PTT:
- `Mask Creation`: Create a mask from a given string.
- `Mask Matching`: Match a mask to a given string.
- `Removing Characters by Mask`: Remove characters from a given string by a mask.
- `Creating Retain/Partial Masks`: Create a mask that retains only certain keywords.
### Mask Creation
Masks replace characters in a string with a common character. The syntax to create a mask is as follows:
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

The default value is `uldsb` for all characters. The `-v` flag is optional and, if provided, will print the length of the original string, the length, the complexity, and the remaining mask keyspace. The format will be `:length:complexity:mask-keyspace` appended to the end of the output. The mask keyspace is the number of possible combinations for the masked portion of the string.
```
$ echo 'HelloWorld!I<3ThePasswordTransformationToolPr0j3ct' | ptt -t mask -rm ds -v
[*] All input loaded.
[*] Task complete with 1 unique results.
1 HelloWorld?sI?s?dThePasswordTransformationToolPr?dj?dct:50:4:94
```
### Mask Matching
Masks can be matched to a given string to determine if the string matches the mask. The syntax to match a mask is as follows:
```
ptt -f <input_file> -t mask-match -tf <mask_file>
```
Where `<mask_file>` is the file containing the mask to match. The mask file should only contain valid masks. The output will be all of the strings that match the masks.
### Removing Characters by Mask
Characters can be removed from a string by a mask. The syntax to remove characters by mask is as follows:
```
ptt -f <input_file> -t mask-remove -rm <mask_characters>
```
Where `<mask_characters>` is the mask to remove from the string. The output will be the string with the characters removed.
### Creating Retain/Partial Masks
Retain masks or partial masks can be created to retain only certain keywords in a string. The `-v` flag is optional and, if provided, will print the length of the original string, the length, the complexity, and the remaining mask keyspace. The syntax to create a retain mask is as follows:
```
ptt -f <input_file> -t mask-retain -rm <mask_characters> -tf <keep_file> -v
```
Where `<mask_characters>` is the mask to retain and `<keep_file>` is the file containing the keywords to retain. The output will be the mask with only the keywords retained.

The `retain` mode can also be used with `-rm` to alter the replacement mask and recieve different output.
```
$ echo 'sp-test1337' | ptt -t retain -tf keep.tmp
[*] Reading files for input.
[*] All input loaded.
[*] Task complete with 2 unique results.
sp-?l?l?l?l?d?d?d?d
?l?l?s?l?l?l?l1337

$ echo 'sp-test1337' | ptt -t retain -tf keep.tmp -rm l
[*] Reading files for input.
[*] All input loaded.
[*] Task complete with 2 unique results.
sp-?l?l?l?l1337
?l?l-?l?l?l?l1337

$ cat keep.tmp
sp-
1337
```

## Rule Transformation Usage
There are several types of rules that can be created using PTT:
- `Append Rules`: Append a string to the end of the password.
- `Append Remove Rules`: Remove characters from the end of the password before appending a string.
- `Prepend Rules`: Prepend a string to the beginning of the password.
- `Prepend Remove Rules`: Remove characters from the beginning of the password before prepending a string.
- `Prepend Toggle Rules`: Toggle the case of the password where a string is prepended.
- `Toggle Rules`: Toggle the case of the password.
- `Insert Rules`: Insert a string at a specific position in the password.
- `Overwrite Rules`: Overwrite a string at a specific position in the password.
### Append Rules
Append rules are used to append a string to the end of the password. The syntax for an append rule is as follows:
```
ptt -f <input_file> -t rule-append
```

The append mode also has additional options:
- `append-remove`: Remove characters from the end of the password before appending a string.

The syntax for an append-remove rule is as follows:
```
ptt -f <input_file> -t rule-append-remove
```
### Prepend Rules
Prepend rules are used to prepend a string to the beginning of the password. The syntax for a prepend rule is as follows:
```
ptt -f <input_file> -t rule-prepend
```

The prepend mode also has two additional options:
- `prepend-remove`: Remove characters from the beginning of the password before prepending a string.

The syntax for a prepend-remove rule is as follows:
```
ptt -f <input_file> -t rule-prepend-remove
```

- `prepend-toggle`: Toggle the case of the password where a string is  prepended. Creating camel and pascal case passwords.

The syntax for a prepend-toggle rule is as follows:
```
ptt -f <input_file> -t rule-prepend-toggle
```
### Toggle Rules
Toggle rules are used to toggle the case of the password. The syntax for a toggle rule is as follows:
```
ptt -f <input_file> -t rule-toggle -i <index>
```
Where `<index>` is the starting index of the toggle pattern. If no index is provided, the toggle pattern will start at the beginning of the password.
### Insert Rules
Insert rules are used to insert a string at a specific position in the password. The syntax for an insert rule is as follows:
```
ptt -f <input_file> -t rule-insert -i <index>
```
Where `<index>` is the position where the string will be inserted. If no index is provided, the string will be inserted at the beginning of the password. The `<index>` can also accept range values in the format of `start-end`. For example, `1-5` will print output for the insert transformation starting from index 1 to 5.
### Overwrite Rules
Overwrite rules are used to overwrite a string at a specific position in the password. The syntax for an overwrite rule is as follows:
```
ptt -f <input_file> -t rule-overwrite -i <index>
```
Where `<index>` is the position where the string will be overwritten. If no index is provided, the string will be overwritten at the beginning of the password. The `<index>` can also accept range values in the format of `start-end`. For example, `1-5` will print output for the overwrite transformation starting from index 1 to 5.
## Wordlist Creation Usage
There are several ways to generate wordlists using PTT:
- `Direct Swapping`: Swapping characters directly with a `:` separated file.
   This is implemented in the `swap-single` module.
- `Replacing Text and Characters`: Replacing text and characters in a string.
  This is implemented in the `replace` module
- `Token Popping`: Generates tokens by popping strings at character boundaries.
  This is implemented in the `pop` module.
- `Token Swapping`: Generates tokens by swapping characters in a string. This is
  implemented in the `mask-swap` module.
- `Passphrases`: Generates passphrases by reforming sentences. This is implemented
  in the `passphrase` module.
### Direct Swapping
The `swap-single` module swaps characters directly with a `:` separated file. The syntax is as follows:
```
ptt -f <input-file> -t swap-single -tf <replacement-file>
```
The replacement file should contain the strings to be transformed as `PRIOR:POST` pairs. The replacements will be applied to all instances in each line, but only one swap is applied at a time. This mode is ideal for substituting words or characters in a string.
### Replacing Text and Characters
The `replace-all` module replaces text and characters in a string. This mode replaces all strings with all matches from a ':' separated file. The syntax is as follows:
```
ptt -f <input-file> -t replace-all -tf <replacement-file>
```
The replacement file should contain the strings to be transformed as `PRIOR:POST` pairs. The replacements will be applied to all instances in each line, and all replacements will be applied to the string. This mode is ideal for replacing all instances of a word or character in a string.
### Token Popping
The `pop` module generates tokens by popping strings at character boundaries. The syntax is as follows:
```
ptt -f <input-file> -t mask-pop -rm <mask-characters>
```
Where `<mask_characters>` can be any of the following:
- `u`: Uppercase characters
- `l`: Lowercase characters
- `d`: Digits
- `s`: Special characters
- `b`: Byte characters
- `t`: Title case words (requires `u` and `l`)
- Multiple characters can be combined to create a mask.

The default value is `uldsbt` for all characters. This mode will create tokens by popping characters from the input string then aggregating the results.
### Token Swapping
The `mask-swap` module generates tokens by swapping characters in a string. The
syntax is as follows:
```
ptt -f <input-file> -t mask-swap -tf <replacement-file>
```
> [!NOTE]
> The input for `mask-swap` is partial masks. This is different from other modes.

The replacement file does not need to be in any specific format. The replacements will be applied to the first instance in each line. The `mask-swap` mode is unique in that it uses partial masks from the `retain` module to generate new candidates. This mode also uses its replacer module (different from the other modes) to generate new candidates by extracting the masks and then matching them to the replacement file.
#### Token Swapping Example
```bash
$ cat pass.lst
love@123
@123love

$ cat retain.txt
love
```

Create retain masks:
```bash
$ ptt -f pass.lst -tf retain.txt -t mask-retain | tee retained.mask
?s?d?d?dlove
love?s?d?d?d
```

Then swap on them with matching values:
```bash
$ cat swap.lst
$333
#888
#123

$ ptt -f retained.mask -tf swap.lst -t mask-swap
#123love
$333love
love$333
love#888
love#123
#888love
```
### Passphrases
The `passphrase` module generates passphrases by reforming sentences. The syntax is as follows:
```
ptt -f <input-file> -t passphrase -w <word-count>
```
The `passphrase` mode will generate new passphrases from the input by reformatting the sentences into new passphrases. The number of words to use in the passphrase is specified by the `-w` flag. The output will be the new passphrases generated from the input with the specified word count.
## Misc. Transformation Usage
There are several types that can be created using PTT:
- `Encoding and Decoding`: This transforms input to and from HTML and Unicode escaped strings.
- `Hex and Dehex`: This transforms input to and from `$HEX[....]` strings.
- `Substrings`: This extracts substrings from the input based on position.
### Encoding and Decoding
This mode allows encoding and decoding of input to and from HTML and Unicode escaped strings.
The syntax is as follows:
```
ptt -f <input_file> -t encode
ptt -f <input_file> -t decode
```

The following table shows the supported transformations:

| Transformation | Description | Input Example | Output Example |
| --- | --- | --- | --- |
| `html` | HTML encoding | `<html>` | `&lt;html&gt;` |
| `unicode` | Unicode encoding | `Hello😎` | `Hello\u1f60e` |
### Hex and Dehex
This mode allows encoding and decoding of input to and from `$HEX[....]` strings.
The syntax is as follows:
```
ptt -f <input_file> -t hex
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
This transformation extracts the substring from the input based on the provided index. If the end index is greater than the length of the input, it will be changed to the length of the input.

This transformation can be used to extract specific parts of the input for
further processing.
### Regram
This mode allows 'regramming' sentences into new n-grams with a given number of words. The syntax is as follows:
```
ptt -f <input_file> -t regram -w <word_count>
```
The `regram` transformation will generate new n-grams from the input by combining words from the input. The number of words to use in the n-gram is specified by the `-w` flag. The output will be the new n-grams generated from the input.

### Rule Application
This mode allows applying rules to the input. The syntax is as follows:
```
ptt -f <input_file> -t rule-apply -tf <rule_file>
```
The `rule-apply` transformation will apply rules from the rule file to the input. The rule file should contain the rules to be applied to the input. The output will be the input with the rules applied. This feature is enabled by the work done on the [HCRE](https://git.launchpad.net/hcre/tree/README.md) project. Please consider visiting and supporting the project.

### Rule Simplification
This mode allows simplifying rules from the input. The syntax is as follows:
```
ptt -f <input_file> -t rule-simplify
```
The `rule-simplify` transformation will simplify rules from the input. The output will be the simplified rules equivalent to the input. This feature is enabled by the work done on the [HCRE](https://git.launchpad.net/hcre/tree/README.md) project. Please consider visiting and supporting the project.
