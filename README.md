# `argparser` - Unique(ish) approach to parsing command line args in golang
## What does this do
This package is for parsing command line args. It can perform type conversion to basic types (int, string, etc). It also allows for long and short form of arguments (eg `-a` and `--option-a` mean the same thing). It also allows for multiple short form arguments at once (`-abc` is the same as `-a -b -c`).
## Usage
The `argparser_test.go` has an example of usage in it
## Why is this different to other parsing packages (that I have found with a quick google search)
This package does not parse all args at once, but instead parses sections of args.