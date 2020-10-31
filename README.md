# Roll20 d20 stat generator
This is a frivolous bit of code to show I always get worse numbers on roll20 than everyone else. It takes a chat archive data dump and parses that to produce a set of roll20 dice results which it orders by time and then outputs stats suitable for throwing into a discord channel to show results.

## Building

Clone this repository, then you can

`go build ./cmd/stat_dumper`

This will place the output binary in the current directory. See the `go` tools build help for more information on how to install or build elsewhere. 

## Testing

`go test ./pkg/roll20msg`

Though honestly the testing is super basic; just enough to make sure (some of) the main json we care about is not compeltely broken. This could do with some proper testing.

## Usage

Download your game's chat archive to disk as an HTML file, and extract the data that is assigned to
`var msgdata =`

Don't include the quotes around the (probably) enormous data block.

Save this to another file, and the run this tool on it 

`stat_dumper --file /path/to/data`

The tool wil run and produce output.


## Contributing

Any useful (as I deem it) contributions are welcome.
