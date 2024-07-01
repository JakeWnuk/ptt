# Password Transformation Tool (PTT) Misc Creation Guide
## Version 0.1.0

### Table of Contents
1. [Introduction](#introduction)
2. [Encoding and Decoding](#encoding-and-decoding)
3. [Hex and Dehex](#hex-and-dehex)
4. [Substrings](#substrings)

### Introduction
This document describes the ways to use PTT to create miscellaneous transformations.
There are several types that can be created using PTT:

- `Encoding and Decoding`: This transforms input to and from URL, HTML, and Unicode escaped strings.
- `Hex and Dehex`: This transforms input to and from `$HEX[....]` strings.

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
