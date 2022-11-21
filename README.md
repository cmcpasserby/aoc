# Advent of Code Input Downloader
aoc is a command line tool for downloading [Advent of Code](https://adventofcode.com) puzzle inputs.

## Usage
- `aoc -h` will print usage instructions and list all available flags
- `aoc -year 2015 -day 1 inputs` would be a common use case and will download the puzzle input for day 1 of 2015 and use the session cookie provided by your `.aocConfig`


## Config file
Config files can be local and live in the working directory, or be global for a user and live in the home directory.
Config filenames are `.aocConfig` and are [toml](https://toml.io/) formatted text files.

#### config file example:
```toml
session_cookie = "your_session_cookie_here"
year = 2015 # if year is not provided in config, it will use year from flags or current date
day = 1     # if day is not provided in config, it will use day from flags or current date
```