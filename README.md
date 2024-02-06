# codemetagenerator

[![GitHub](https://img.shields.io/github/license/cacoco/codemetagenerator)](https://github.com/cacoco/codemetagenerator/blob/main/LICENSE) [![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active) [![go ci build](https://github.com/cacoco/codemetagenerator/actions/workflows/ci.yaml/badge.svg)](https://github.com/cacoco/codemetagenerator/actions/workflows/ci.yaml) [![GitHub Release](https://img.shields.io/github/release/cacoco/codemetagenerator?style=flat)](https://github.com/cacoco/codemetagenerator/releases/latest)

This project is a [CodeMeta](https://codemeta.github.io/) project description generator written in [Go](https://go.dev/).

## Installation

### Manually
To install this project, you need to have [Go](https://go.dev/) installed on your machine. Once you have Go installed, you can clone this repository and build the project.

```bash
git clone https://github.com/cacoco/codemetagenerator.git
cd codemetagenerator
go build
```

Then install via `go install`

### Via [Homebrew](https://brew.sh/)

```bash
brew tap cacoco/tap
brew install codemetagenerator
```

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
  delete      Delete an arbitrary key and its value from the in-progress codemeta.json file.
  generate    Generate the final codemeta.json file to the optional output file or to the console
  help        Help about any command
  licenses    List (or refresh cached) SPDX license IDs
  new         Start a new codemeta.json file. When complete, run "codemetagenerator generate" to generate the final codemeta.json file
  set         Set the value of an arbitrary key in the in-progress codemeta.json file.
```

#### New
To start a new `codemeta.json` file:

```bash
codemetagenerator new
```

This will walk you through an interactive session and will store an "in-progress" `codemeta.json` file. The expectation is that
you will continue to add more metadata, e.g., `author`, `contributor`, or `keyword`. 

Optionally, the `--input` (or `-i`) can be passed with a path to a file to load as the new in-progress `codemeta.json` file. This 
allows for editing and updating an existing file which you can then choose to generate back into the original location.

#### Add
Add helps with the addition of 3 specific fields: `author`, `contributor`, or `keyword`. The command provides a wizard for adding these values. 

For example to add an `author`:

```bash
codemetagenerator add author
```

will walk you through an interactive session to create a [`Person`](https://schema.org/Person) or [`Organization`](https://schema.org/Organization) author. This command can be run multiple times to add more authors. Similarly, this can also be done for adding one or more contributors:

```bash
codemetagenerator add contributor
```

Again, this command can be run multiple times to add more contributors.

To add one or more keywords:

```bash
codemetagenerator add keyword "Java" "JVM" "etc"
```

The `keyword` command accepts multiple terms but can also be run multiple times to add more keywords.

#### Delete
Properties can be removed by running the `delete` command. This allows for removing *any value* in the `codemeta.json` file for a given key specified via the [Path Syntax](#path-syntax).

```bash
codemetagenerator delete "version"
```

#### Set
New properties can be inserted or exisiting properties can be updated in the `codemeta.json` file with the `set` command. Keys must be specified via the [Path Syntax](#path-syntax).

```bash
codemetagenerator set "relatedLink" "https://thisisrelated.org"
```

#### Generate
Once creating and editing are done, you can run `generate` to produce a final `codemeta.json`. Generation accepts an optional `-o | --output` flag that allows for specifying an output file. If this flag is not provided, the output is sent to the console.

```bash
codemetagenerator generate
```

### Path Syntax
The syntax follows the `sjson` ([https://github.com/tidwall/sjson](https://github.com/tidwall/sjson)) [Path Syntax](https://github.com/tidwall/sjson?tab=readme-ov-file#path-syntax).


## CodeMeta
[CodeMeta](https://codemeta.github.io) is a [JSON-LD](https://json-ld.org/) file format used to describe software projects.

### CodeMeta Metadata Usage

A CodeMeta instance file describes the metadata of a software object using [JSON-LD](https://json-ld.org/) notation. An instance file is typically located at the root of your repository and can contain any of the properties described on the [CodeMeta terms page](https://codemeta.github.io/terms/). The instance file is generally named, `codemeta.json`.

JSON-LD uses a [context](https://niem.github.io/json/reference/json-ld/context/) file to associate JSON names with IRIs ([Internationalized Resource Identifier](https://en.wikipedia.org/wiki/Internationalized_Resource_Identifier)). The JSON names then serve as abbreviated, local names for the IRIs that are universally unique identifiers for concepts from widely used schemas such as [schema.org](https://schema.org/). 

The context file [codemeta.jsonld](https://raw.githubusercontent.com/codemeta/codemeta/master/codemeta.jsonld) contains the complete set of JSON properties adopted by the CodeMeta project.

### Example

A full example can be seen in the [ropensci/codemeter](https://github.com/ropensci/codemetar) GitHub project page: [https://github.com/ropensci/codemetar/blob/main/codemeta.json](https://github.com/ropensci/codemetar/blob/main/codemeta.json).

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[Apache License 2.0](https://spdx.org/licenses/Apache-2.0.html)

Copyright 2024 Christopher Coco

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
