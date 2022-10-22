# `argparser` - Unique(ish) approach to parsing command line args in golang
## What does this do
This package is for parsing command line args. It can perform type conversion to basic types (int, string, etc). It also allows for long and short form of arguments (eg `-a` and `--option-a` mean the same thing). It also allows for multiple short form arguments at once (`-abc` is the same as `-a -b -c`).
## Usage
Look inside the argparser_test.go file for an example usage
## Why is this defferent to other parsing packages (that I have found with a quick google search)
This package does not parse all args at once, but instead parses sections of args. This allows for stuff like this easily:
### Example
Take these two command using a made up utility (one to run a file and one to build it):

`$ dummy run --cache-location /path/to/cache target.dummy`
`$ dummy build --output-location /path/to/output target.dummy`

You can only supply the `--cache-location` parameter to the `run` sub-command, and the `--output-location` to the `build` sub-command. It is invalid to supply the option to the wrong sub-command. The code for this using this package is very simple, due to partial parsing of the arguments at once. See an example of this in the argparser_dummy_test.go file.