# Password Transformation Tool (PTT) Wordlist Creation Guide
## Version 0.1.0

### Table of Contents
1. [Introduction](#introduction)
2. [Direct Swapping](#direct-swapping)
3. [Token Popping](#token-popping)
4. [Token Swapping](#token-swapping)

### Introduction
This document describes the ways to use PTT to create password cracking
wordlists. There are several ways to generate wordlists using PTT:

- `direct-swapping`: Swapping characters directly with a `:` separated file.
   This is implemented in the `swap` module.
- `token-popping`: Generates tokens by popping strings at character boundaries.
  This is implemented in the `pop` module.
- `token-swapping`: Generates tokens by swapping characters in a string. This is
  implemented in the `mask-swap` module.
- `passphrases`: Generates passphrases by combining words from a wordlist. This
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
pairs. The replacements will be applied to the first instance in each line.

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
- Multiple characters can be combined to create a mask.

The default value is `uldsb` for all characters. This mode will create tokens
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
