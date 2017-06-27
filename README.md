# Description

ansigo converts ANSi art files to PNG files. ansigo is a Go port of [ansilove/C][1] that aims to have as few external dependencies as possible. It provides both a command-line application and a package interface to allow you to integrate it with your own applications.

ansigo supports all of the options from ansilove/C with the addition of 24bit ANSi support.

# Installation

After cloning this repo or downloading the code, you'll need to ensure that any required dependencies are installed.

This package uses [github.com/nfnt/resize](https://github.com/nfnt/resize). You can either manually place this within the
`vendor` folder, or use a dependency management package.  

While it is still in development, we recommend using the `dep` tool. You can get it from [github.com/golang/dep](https://github.com/golang/dep) or via the Go toolchain:

       $ go get -u github.com/golang/dep/...
       
After installing `dep` you can use it to fetch the dependencies automatically:

       $ dep init
       $ dep ensure -update

# Features

Rendering of all known ANSi / ASCII art file types:

- ANSi (.ANS)
- Binary (.BIN)
- Artworx (.ADF)
- iCE Draw (.IDF)
- Xbin (.XB) [details](http://www.acid.org/info/xbin/xbin.htm)
- PCBoard (.PCB)
- Tundra (.TND) [details](https://sourceforge.net/projects/tundradraw/)
- ASCII (.ASC)
- Release info (.NFO)
- Description in zipfile (.DIZ)

Files with custom suffix default to the ANSi renderer (e.g. ICE or CIA).

ansigo is capabable of processing:

- SAUCE records
- DOS and Amiga fonts (embedded binary dump)
- iCE colors

Even more:

- Output files are highly optimized 4-bit PNGs.
- Optionally generates additional (and proper) Retina @2x PNG.
- You can use custom options for adjusting output results.
- Built-in support for rendering Amiga ASCII.
- Support for 24-bit ANSi

# Documentation

## Synopsis

       ansigo [options] file
       ansigo -e | -h | -v

## Options

       -b bits     set to 9 to render 9th column of block characters (default: 8)
       -c columns  adjust number of columns for BIN files (default: 160)
       -e          print a list of examples
       -f font     select font (default: 80x25)
       -h          show help
       -i          enable iCE colors
       -m mode     set rendering mode for ANS files:
                     ced            black on gray, with 78 columns
                     transparent    render with transparent background
                     workbench      use Amiga Workbench palette
       -o file     specify output filename/path
       -r          creates additional Retina @2x output file
       -s          show SAUCE record without generating output
       -v          show version information

There are certain cases where you need to set options for proper rendering. However, this is occasionally. Results turn out well with the built-in defaults. You may launch ansigo with the option `-e` to get a list of basic examples. Note that columns is restricted to `BIN` and `TND` files, it won't affect other file types.

## Fonts

ansigo inherits all the embedded fonts from ansilove/C as binary data, so the most popular typefaces for rendering ANSi / ASCII art are available at your fingertips.

PC fonts can be (all case-sensitive):

- `80x25` (code page 437)
- `80x50` (code page 437, 80x50 mode)
- `baltic` (code page 775)
- `cyrillic` (code page 855)
- `french-canadian` (code page 863)
- `greek` (code page 737)
- `greek-869` (code page 869)
- `hebrew` (code page 862)
- `icelandic` (Code page 861)
- `latin1` (code page 850)
- `latin2` (code page 852)
- `nordic` (code page 865)
- `portuguese` (Code page 860)
- `russian` (code page 866)
- `terminus` (modern font, code page 437)
- `turkish` (code page 857)

AMIGA fonts can be (all case-sensitive):

- `amiga` (alias to Topaz)
- `microknight` (Original MicroKnight version)
- `microknight+` (Modified MicroKnight version)
- `mosoul` (Original mO'sOul font)
- `pot-noodle` (Original P0T-NOoDLE font)
- `topaz` (Original Topaz Kickstart 2.x version)
- `topaz+` (Modified Topaz Kickstart 2.x+ version)
- `topaz500` (Original Topaz Kickstart 1.x version)
- `topaz500+` (Modified Topaz Kickstart 1.x version)

## Bits

`bits` can be (all case-sensitive):

- `8` (8-bit)
- `9` (9-bit)

Setting the bits to `9` will render the 9th column of block characters, so the output will look like it is displayed in real textmode.

## Rendering Mode

`mode` can be (all case-sensitive):

- `ced`
- `transparent`
- `workbench`

Setting the mode to `ced` will cause the input file to be rendered in black on gray, and limit the output to 78 columns (only available for `ANS` files). Used together with an Amiga font, the output will look like it is displayed on Amiga.

Setting the mode to `workbench` will cause the input file to be rendered using Amiga Workbench colors (only available for `ANS` files).

Settings the mode to `transparent` will produce output files with transparent background (only available for `ANS` files).

## iCE Colors

iCE colors are disabled by default, and can be enabled by specifying the `-i` option.

When an ANSi source was created using iCE colors, it was done with a special mode where the blinking was disabled, and you had 16 background colors available. Basically, you had the same choice for background colors as for foreground colors, that's iCE colors.

## Columns

`columns` is only relevant for .BIN files, and even for those files is optional. In most cases conversion will work fine if you don't set this flag, the default value is `160` then. So please pass `columns` only to `BIN` files and only if you exactly know what you're doing.

## SAUCE records

You can use ansigo as SAUCE reader without generating any output, just use option `-s` for this purpose.

# License

ansigo is released under the BSD 3-Clause license. See `LICENSE` file for details.

# Author

ansigo is developed by [Pete Garcin](http://rawktron.com).  

Based on [ansilove/C][1]. ansilove is developed by Stefan Vogt, Brian Cassidy, [Frederic Cambus](http://www.cambus.net).

# Resources

Project Homepage : [https://github.com/ActiveState/ansigo](https://github.com/ActiveState/ansigo)


[1]: https://github.com/ByteProject/ansilove-C
