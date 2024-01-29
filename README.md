# codemetagenerator
This project is a [CodeMeta](https://codemeta.github.io/) project description generator written in [Go](https://go.dev/).

## Installation
To install this project, you need to have [Go](https://go.dev/) installed on your machine. Once you have Go installed, you can clone this repository and build the project.

```bash
git clone https://github.com/cacoco/codemetagenerator.git
cd codemetagenerator
go build
```

Then install via `go install`

## Usage
To run this project, you can use the following commands:

```bash 
codemetagenerator --help
```

### Commands
```bash
Available Commands:
  add         Adds resources [authors, contributors, keywords] to the in-progress codemeta.json file
  clean       Clean the $HOME/.codemetagenerator directory
  edit        Edit an existing property value by key, e.g., 'foo' or 'foo.bar' or 'foo[1]', or 'foo[1].bar' in the in-progress codemeta.json file
  generate    Generate the final codemeta.json file to the optional output file or to the console
  help        Help about any command
  insert      Insert a new property with the given value by key, e.g., 'foo' or 'foo.bar' or 'foo[1]', or 'foo[1].bar' into the in-progress codemeta.json file
  licenses    List (or refresh cached) SPDX license IDs
  new         Start a new codemeta.json file. When complete, run "codemetagenerator generate" to generate the final codemeta.json file
  remove      Remove a property by key, e.g., 'foo' or 'foo.bar' or 'foo[1]', or 'foo[1].bar' from the in-progress codemeta.json file
```

#### New
To start a new `codemeta.json` file:

```bash
codemetagenerator new
```

This will walk you through an interactive session and will store an "in-progress" `codemeta.json` file. The expectation is that
you will continue to add more metadata, e.g., `author`, `contributor`, or `keyword`.

#### Add
For example to add an `author`:

```bash
codemetagenerator add author
```

This will walk you through an interactive session to create a [`Person`](https://schema.org/Person) or [`Organization`](https://schema.org/Organization) author. This command can be run multiple times to add more authors. Similarly, this can also be done for adding one or more contributors:

```bash
codemetagenerator add contributor
```

Again, this command can be run multiple times to add more contributors.

To add one or more keywords:

```bash
codemetagenerator add keyword "Java" "JVM" "etc"
```

The `keyword` command accepts multiple terms but the command can also be run multiple times to add more keywords.

#### Remove
Properties can be removed by running the `remove` command. This allows for removing *any property* in the `codemeta.json` file including an author, contributor, or keyword.

```bash
codemetagenerator remove "version"
```

#### Insert
New properties can be inserted into object values in the `codemeta.json` file (arrays cannot be added to via insert). This allows for adding any other `Codemeta` key and value not covered in the generator.

```bash
codemetagenerator insert "relatedLink" "https://thisisrelated.org"
```

Note: only a single JSON value can be inserted (e.g. no objects or arrays) for a given property key.

#### Generate
Once creating and editing are done, you can run `generate` to produce a final `codemeta.json`. Generation accepts an optional `-o | --output` flag that allows for specifying an output file. If this flag is not provided, the output is sent to the console.

```bash
codemetagenerator generate
```

## CodeMeta
[CodeMeta](https://codemeta.github.io) is a [JSON-LD](https://json-ld.org/) file format used to describe software projects. See a [full example](https://github.com/ropensci/codemetar/blob/main/codemeta.json).

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[Apache License 2.0](https://spdx.org/licenses/Apache-2.0.html)

Copyright 2024 Christopher Coco

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.