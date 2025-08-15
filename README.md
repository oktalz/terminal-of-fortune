# terminal of fortune

```
┌┬┐┌─┐┬─┐┌┬┐┬┌┐┌┌─┐┬    ┌─┐┌─┐  ┌─┐┌─┐┬─┐┌┬┐┬ ┬┌┐┌┌─┐
 │ ├┤ ├┬┘│││││││├─┤│    │ │├┤   ├┤ │ │├┬┘ │ │ ││││├┤
 ┴ └─┘┴└─┴ ┴┴┘└┘┴ ┴┴─┘  └─┘└    └  └─┘┴└─ ┴ └─┘┘└┘└─┘
```

tool for selecting random person/item from list in terminal
items are stored in txt file

## Installation
Use the following command to download and install this tool:
```sh
go install github.com/oktalz/terminal-of-fortune@latest
```

## Binaries
  prebuilt binaries can be found on [releases](https://github.com/oktalz/terminal-of-fortune/releases) page

## Running

![Example](demo.gif)

```sh
terminal-of-fortune file.txt
```

if more txt files exist, start the program and pick correct one.
if one txt file exists it will be autoselected
```sh
terminal-of-fortune
```

## limitations

number of items does not necessarity fit on screen (if to many are in list)

## todo

- option to remove items that were selected
- option to sort items by current percentage
