# Rename

Usage: rename [options] files... expression

rename is a utility for renaming multiple files, rename takes the list of files to rename and a search/replace pattern
as the last argument. 

If no filenames are given, filenames will be read from stdin.

It's similar to [Unix rename](http://man7.org/linux/man-pages/man1/rename.1.html) and
[Perl Rename](http://manpages.ubuntu.com/manpages/trusty/man1/prename.1.html) written in Go.

# Examples

Backup files, add ".bak" to all files:  ``` rename path/to/files/* 's/$/.bak//' ```

Backup files, add ".bak" to all files (keeping the original):  ``` rename path/to/files/* 's/$/.bak//' -c ```

Remove (last) extension: ``` rename path/to/files/* 's/\.[a-z]+$//' ```

Swap double extensions : ``` rename path/to/files/* 's/([a-z]+)\.([a-z]+)$/$2.$1/' ```

# Installation

## Homebrew

``` brew install fbergen/tap/rename ```


# Options

    -c, --copy=false: Copy instead of move.
    -f, --force=false: Overwrite existing files.
    -h, --help=false: Show help dialog.
    -i, --interactive=false: Ask for confirmation, before renaming
    -n, --no-action=false: Don't perform any changes. Show what files would have been renamed.
    -v, --verbose=false: Show which files where renamed, if any.



