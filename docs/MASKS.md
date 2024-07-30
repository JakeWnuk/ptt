# Password Transformation Tool (PTT) Mask Creation Guide

### Table of Contents
1. [Introduction](#introduction)
2. [Mask Creation](#mask-creation)
3. [Mask Matching](#mask-matching)
4. [Removing Characters by Mask](#removing-characters-by-mask)
5. [Creating Retain/Partial Masks](#creating-retainpartial-masks)

### Introduction
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
