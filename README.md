# wpe_merge

Merge account info into one simple CSV!

## Installation
Written in go 1.13

It's recommended you put these files in `$GOPATH/src/github.com/wpe_merge/wpe_merge` in order to build the binary, otherwise you will have to rename all package references for it to work.

## Testing
Tests use the ginkgo framework (github.com/onsi/ginkgo).  You should be able to run `go test ./...` but you can also run `ginkgo -R` for what I would consider as prettier output.

## Additional Options
Basic usage: `wpe_merge <input_file> <output_file>`

However, if you want to get fancy, you can change the url endpoint using the `--url` flag.  The address needs to be prepended with the protocol (https?) in order for it to be parsed correctly.

Another fun flag to try is `--max-concurrent-requests` which limits the number of concurrent requests made to the remote server. There is no specific reason as to why the default is 10, other than the fact that it is better than 1, which means that every request is handled synchronously.
