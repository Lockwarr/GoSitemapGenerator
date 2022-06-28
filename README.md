# SitemapGenerator

## Requirements

Go 1.17+

## Running the cli tool

### Options

1. -parallel= `number of parallel workers to navigate through site`, default is 10
2. -output-file=`output file path`, default is ../sitemap.xml
3. -max-depth= `max depth of url navigation recursion`, default is 10

### Required arguments

1. A valid url in the form of `https://google.com`

### Example commands

1. Navigate inside sitemap-generator-cli and run:

    go run main.go https://google.com

    go run main.go -parallel=5 -max-depth=10 -output-file="../sitemap.xml" https://google.com

## Testing

1. Run all tests:

    make test

## Linux

The code has been tested only on Ubuntu 18.04

## Improvements && Notes

I used one external dependency for testing only.

Of course the program can be expanded in many ways as required, some examples:

    We can add log levels (debug, info, error etc...)
    We could set up Dockerfile
    I don't think this cli tool requires component tests but they could be implemented.
