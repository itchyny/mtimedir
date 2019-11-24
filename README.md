# mtimedir
Change mtime and atime of the directories based on the files inside.

Sometimes macOS creates a `.DS_Store` file and it breaks the directory modified time.
Also, smart sync of Dropbox updates the modified time.
So mtimedir is here to come.

This tool `mtimedir` changes the modified time of the directory to the max the modified time of the files inside.
This procedure is run recursively, so after running mtimedir, the modified time is set to the max modified time inside all the files in any depth.

## Installation
### Build from source
```bash
go get github.com/itchyny/mtimedir
```

## Bug Tracker
Report bug at [Issues・itchyny/mtimedir - GitHub](https://github.com/itchyny/mtimedir/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
