# DADP.FlacTool
A tool manipulates structure of flac file(s) according to the given task list.  
Compatible with Golang >= 1.12, Windows, macOS, and Linux.

## Features
Manipulate METADATA_BLOCK:  
- `Sort(Filter)`, `Delete`

Modify METADATA_BLOCK_DATA:  
- METADATA_BLOCK_VORBIS_COMMENT
  - Refer: `Set` 
  - Comments: `Add`, `Set`, `Delete`, `Sort(Filter)`, `Dump`, `Import`
- METADATA_BLOCK_PICTURE
  - `Set` type (ID3v2 APIC)
  - `Set` description
  - `Load` binary data from picture file
  - Automatically parse properties from binary data

Configure Task List:  
- Call third-party executables
- Argument Support  
  - Environmental variables
  - Values from METADATA_BLOCK_DATA

## Usage
Download binary files from [Releases](https://gitlab.com/KTGWKenta/DADP.FlacTool/releases) or
build with following command:
```bash
# Clone this repository
$ git clone git@gitlab.com:KTGWKenta/DADP.FlacTool.git

# Go into the repository
$ cd DADP.FlacTool

# Build (for Windows)
$ build.cmd

# Build (for Linux & macOS)
$ chmod +x ./build.sh
$ ./build.sh
```
To process single flac file, call `./bin/main` in the following format:
```bash
$ ./bin/main -input "source.flac" -output "output.flac" -task "taskList.yaml"
```
To process a bunch of flac file, place all source files in the same directory and call
`./bin/loop` in the following format:
```bash
$ ./bin/loop -exe ./bin/main -inputDir "sourceDir" -outputDir "outputDir" -task "taskList.yaml"
``` 
~~Format specification of task list is available at wiki, 
For more advanced usages, please refer to examples here~~(Comming soon). 

Feel free to commit issues if you come across any problem.  




## License
Copyright (c) 2019 Digimon Audiopedia Database Project.  
Licensed under GNU General Public License v3.0.  
See the [LICENSE](https://gitlab.com/KTGWKenta/DADP.FlacTool/blob/master/LICENSE.md) for more information.