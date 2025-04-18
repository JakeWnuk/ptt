# Password Transformation Tool (PTT) Usage Guide
> Version 1.1.1
## Table of Contents
- [Introduction](#introduction)
- [Installation](#installation)
- [Examples](#examples)

## Introduction
The Password Transformation Tool (PTT) is a command-line utility that allows users to transform passwords using various methods.This guide will provide instructions on how to install and use the tool. The tool was created as a complete solution for password transformation and is designed to be easy to use and flexible. PTT is designed around CLI tools and offers many features and capabilities with a user-friendly interface. The tool focuses on core features and being a memory-light pre-processor. 

> The `rule-simplify` transformation feature is enabled by the work done on the [HCRE](https://git.launchpad.net/hcre/tree/README.md) project.

### Installation
From source with `go`:
```
go install github.com/jakewnuk/ptt@latest
```
From `git` clone then build with `go`:
```
git clone https://github.com/JakeWnuk/ptt && cd ptt && go build ./main.go && mv ./main ~/go/bin/ptt && ptt
```

### Examples
We will be using RockYou as a foundation for examples. To apply
a transformation mode, use the `-t` flag along with a valid option. Multiple
`-t` flags can be used at once.

```
$ cat rockyou.txt | ptt -t mask -v
[*] Starting statistics generation. Please wait...
Verbose Statistics: max=75
--------------------------------------------------
General Stats:
Total Items: 14343188
Total Unique items: 148329
Total Words: 148329
Largest frequency: 687991
Smallest frequency: 1

Plots:
Item Length: |------[=|==]----------------------------------------|
Min: 2, Q1: 14, Q2: 16, Q3: 20, Max: 98
Item Frequency: |[|]--------------------------------------------------|
Min: 1, Q1: 1, Q2: 1, Q3: 2, Max: 687991
Item Complexity: |-------------[|]-------------------------|
Min: 1, Q1: 2, Q2: 2, Q3: 2, Max: 4

Category Counts:
long-non-complex: 147639
short-non-complex: 690
non-ASCII: 177
alphanumeric-with-special: 148323
all-lowercase: 148323
non-complex: 148329

--------------------------------------------------
?l?l?l?l?l?l?l?l                 [687991]==================================================
?l?l?l?l?l?l                     [601152]===========================================
?l?l?l?l?l?l?l                   [585013]==========================================
?l?l?l?l?l?l?l?l?l               [516830]=====================================
?d?d?d?d?d?d?d                   [487429]===================================
?d?d?d?d?d?d?d?d?d?d             [478196]==================================
?d?d?d?d?d?d?d?d                 [428296]===============================
?l?l?l?l?l?l?d?d                 [420318]==============================
?l?l?l?l?l?l?l?l?l?l             [416939]==============================
?d?d?d?d?d?d                     [390529]============================
?d?d?d?d?d?d?d?d?d               [307532]======================
?l?l?l?l?l?d?d                   [292306]=====================
?l?l?l?l?l?l?l?d?d               [273624]===================
?l?l?l?l?l?l?l?l?l?l?l           [267733]===================
?l?l?l?l?d?d?d?d                 [235360]=================
?l?l?l?l?d?d                     [215074]===============
?l?l?l?l?l?l?l?l?d?d             [213109]===============
?l?l?l?l?l?l?d                   [193097]==============
?l?l?l?l?l?l?l?d                 [189847]=============
?l?l?l?l?l?l?l?l?l?l?l?l         [189355]=============
?l?l?l?d?d?d?d                   [178304]============
?l?l?l?l?l?d?d?d?d               [173559]============
?l?l?l?l?l?l?d?d?d?d             [160592]===========
?l?l?l?l?l?l?l?l?d               [160054]===========
?l?l?l?l?l?d?d?d                 [152400]===========
?l?l?l?l?l?l?d?d?d               [132216]=========
?l?l?l?l?l?l?l?l?l?d             [129823]=========
?l?l?l?l?l                       [125731]=========
?l?l?l?l?l?l?l?l?l?l?l?l?l       [119294]========
?l?l?l?l?l?d                     [114732]========
?l?l?l?l?d?d?d                   [111218]========
?d?d?d?d?d?d?d?d?d?d?d           [107862]=======
?l?l?d?d?d?d                     [98305]=======
?l?l?l?d?d?d                     [98183]=======
?l?l?l?l?l?l?l?d?d?d             [87611]======
?l?l?l?l?l?l?l?l?l?d?d           [82654]======
?l?l?l?l?l?l?l?l?l?l?l?l?l?l     [80333]=====
?l?l?l?l?l?l?l?d?d?d?d           [70914]=====
?l?l?l?l?l?l?l?l?l?l?l?l?l?l?l   [55398]====
?l?d?d?d?d?d?d                   [54883]===
?u?u?u?u?u?u                     [51839]===
?l?l?d?d?d?d?d?d                 [48541]===
?l?l?l?l?l?l?l?l?d?d?d?d         [45499]===
?d?d?d?d?d                       [44987]===
?l?l?l?d?d?d?d?d?d               [44792]===
?l?l?l?l?l?l?l?l?d?d?d           [43215]===
?d?d?d?d?d?d?l                   [41557]===
?u?u?u?u?u?u?u                   [40592]==
?u?u?u?u?u?u?u?u                 [39457]==
?d?d?d?d?d?d?d?d?d?d?d?d         [38464]==
?l?l?l?d?d?d?d?d                 [37622]==
?l?l?l?l?l?l?l?l?l?l?d?d         [35980]==
?l?l?l?l?l?l?l?l?l?l?l?l?l?l?l?l [33483]==
?l?l?l?l?d?d?d?d?d?d             [33277]==
?l?l?d?d?d?d?d                   [32540]==
?u?u?u?u?u?u?d?d                 [31373]==
?d?d?d?d?l?l                     [31086]==
?l?d?d?d?d?d?d?d                 [29589]==
?d?d?d?d?d?d?d?d?d?d?d?d?d       [28908]==
?u?l?l?l?l?l?d?d                 [27662]==
?d?d?d?d?l?l?l?l                 [27300]=
?u?u?u?u?u?u?u?u?u               [27019]=
?u?u?u?u?u?d?d                   [26011]=
?l?l?l?l?l?l?l?l?l?d?d?d         [25912]=
?l?l?l?l?l?l?l?l?l?d?d?d?d       [24714]=
?d?d?d?d?d?d?l?l                 [23385]=
?l?l?l?l?l?l?s                   [23126]=
?d?d?d?d?l?l?l                   [23114]=
?d?l?l?l?l?l?l?l                 [22690]=
?d?l?l?l?l?l?l                   [22565]=
?l?l?l?l?d?d?d?d?d               [22362]=
?d?d?d?d?d?d?d?l                 [22313]=
?u?u?u?u?d?d                     [22224]=
?l?l?l?l?l?d?d?d?d?d             [20002]=
?d?d?l?l?l?l?l?l                 [19937]=
```

Some transformations use the `-m` flag to specify a mask to use for that
transformation.
```
$ cat rockyou.txt | ptt -t mask -m d -v
[*] Starting statistics generation. Please wait...
Verbose Statistics: max=75
--------------------------------------------------
[-] Please wait loading. Elapsed: 00:00:10.000. Memory Usage: 2058.85 MB.
General Stats:
Total Items: 14343636
Total Unique items: 8480067
Total Words: 8570299
Largest frequency: 487429
Smallest frequency: 1

Plots:
Item Length: |----[=|==]-------------------------------------------|
Min: 1, Q1: 9, Q2: 11, Q3: 14, Max: 99
Item Frequency: |[|]--------------------------------------------------|
Min: 1, Q1: 1, Q2: 1, Q3: 1, Max: 487429
Item Complexity: |----------[==========|]------------------------------|
Min: 0, Q1: 1, Q2: 2, Q3: 2, Max: 5

Category Counts:
all-uppercase: 256672
URL: 197
hebrew-characters: 75
hiragana-characters: 3
long-complex: 2
numeric: 5
all-lowercase: 7406769
non-complex: 8480065
long-non-complex: 3324687
contains-uppercase: 148573
phrase: 17066
numeric-with-special: 11
complex: 2
korean-characters: 1
alphabetical: 4177376
greek-characters: 88
hex-string: 1539
cyrillic-characters: 66
short-non-complex: 5155378
starts-uppercase: 919213
non-ASCII: 14314
arabic-characters: 268
CJK-characters: 12
alphanumeric: 13
alphanumeric-with-special: 4300954
thai-characters: 4841

--------------------------------------------------
?d?d?d?d?d?d?d                   [487429]==================================================
?d?d?d?d?d?d?d?d?d?d             [478196]=================================================
?d?d?d?d?d?d?d?d                 [428296]===========================================
?d?d?d?d?d?d                     [390529]========================================
?d?d?d?d?d?d?d?d?d               [307532]===============================
?d?d?d?d?d?d?d?d?d?d?d           [107862]===========
?d?d?d?d?d                       [44987]====
?d?d?d?d?d?d?d?d?d?d?d?d         [38464]===
?d?d?d?d?d?d?d?d?d?d?d?d?d       [28908]==
?d?d?d?d?d?d?d?d?d?d?d?d?d?d     [11678]=
a?d?d?d?d?d?d                    [7428]
?d?d?d?d?d?da                    [6999]
?d?d?d?d                         [6359]
?d?d?d?d?d?d?d?d?d?d?d?d?d?d?d?d [6122]
j?d?d?d?d?d?d                    [5481]
?d?d?d?d?d?d?d?d?d?d?d?d?d?d?d   [5403]
m?d?d?d?d?d?d                    [4796]
a?d?d?d?d?d?d?d                  [4261]
s?d?d?d?d?d?d                    [4014]
k?d?d?d?d?d?d                    [3928]
?d?d?d?d?d?d?da                  [3874]
?d?d?d?d?d?dj                    [3822]
?d?d?d?d?d?dm                    [3497]
c?d?d?d?d?d?d                    [3328]
b?d?d?d?d?d?d                    [3117]
d?d?d?d?d?d?d                    [2891]
?d?d?d?d?d?dk                    [2651]
?d?d?d?d?d?ds                    [2637]
a?d?d?d?d?d?d?d?d                [2510]
j?d?d?d?d?d?d?d                  [2453]
love?d?d?d?d                     [2443]
m?d?d?d?d?d?d?d                  [2422]
s?d?d?d?d?d?d?d                  [2412]
t?d?d?d?d?d?d                    [2310]
?d?d?d?d?d?dc                    [2305]
?d?d?d?d?d?db                    [2289]
?d?d/?d?d/?d?d                   [2263]
a?d?d?d?d?d                      [2240]
?d?d?d?d?d?d?d?da                [2213]
l?d?d?d?d?d?d                    [2190]
?d?d?d?d?d?dd                    [2064]
r?d?d?d?d?d?d                    [2057]
?d?d?d?d?d?d?dj                  [1929]
?d?d?d?d?d?dl                    [1836]
e?d?d?d?d?d?d                    [1820]
?d?d?d?d?da                      [1811]
k?d?d?d?d?d?d?d                  [1790]
b?d?d?d?d?d?d?d                  [1769]
c?d?d?d?d?d?d?d                  [1755]
?d?d?d?d?d?d?dm                  [1747]
j?d?d?d?d?d                      [1740]
n?d?d?d?d?d?d                    [1659]
?d?d?d?d?d?dt                    [1641]
?d?d?d?d?d?dr                    [1552]
d?d?d?d?d?d?d?d                  [1517]
m?d?d?d?d?d?d?d?d                [1473]
m?d?d?d?d?d                      [1466]
j?d?d?d?d?d?d?d?d                [1443]
p?d?d?d?d?d?d                    [1418]
g?d?d?d?d?d?d                    [1410]
?d?d?d?d?d?de                    [1373]
angel?d?d?d?d                    [1364]
?d?d?d?d?d?d?dk                  [1348]
k?d?d?d?d?d                      [1315]
may?d?d?d?d                      [1311]
?d?d?d?d?d?d?dc                  [1295]
?d?d?d?d?d?d?ds                  [1287]
s?d?d?d?d?d?d?d?d                [1251]
?d?d/?d?d/?d?d?d?d               [1251]
s?d?d?d?d?d                      [1239]
?d?d?d?d?d?d?db                  [1230]
h?d?d?d?d?d?d                    [1230]
baby?d?d?d?d                     [1202]
may?d?d?d?d?d?d                  [1202]
alex?d?d?d?d                     [1183]
```

Some modes, like many of the phrase options, use `-w` to specify the number of
words that should be returned when transforming. Most options support range
formats seperated by `-`. The length range of the output strings can also be
limited with `-l`.

```
cat rockyou.txt | ptt -t regram -w 2 -l 5
cat rockyou.txt | ptt -t regram -w 2-5 -l 5-12
```

The tool can be used to make rules as well, the `-i` flag will affect the
`overwrite`, `insert`, and `toggle` modes when creating new rules by changing
the default starting index.

```
$ echo 'Test' | ptt -t insert -t overwrite -t toggle -i 4-6
i4T i5e i6s i7t
i5T i6e i7s i8t
i6T i7e i8s i9t
o4T o5e o6s o7t
o5T o6e o7s o8t
o6T o7e o8s o9t
T4
T5
T6
```

Lastly, the tool features some pre-processor options like `token-swap` to
perform token swapping on the input using retain masks to preserve top tokens
before performing token swapping.

```
$ cat test.txt
test123
zest456

$ cat test.txt | ptt -t token-swap
test123
zest456
zest123
test456
```

- https://jakewnuk.com/posts/concept-of-token-swapping-attacks/
- https://jakewnuk.com/posts/improving-token-swapping-attacks/
