# Password Transformation Tool (PTT) Rules Creation Guide
## Version 0.1.0

### Table of Contents
1. [Introduction](#introduction)
2. [Append Rules](#append-rules)
3. [Prepend Rules](#prepend-rules)
4. [Toggle Rules](#toggle-rules)
5. [Insert Rules](#insert-rules)
6. [Overwrite Rules](#overwrite-rules)

### Introduction
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

