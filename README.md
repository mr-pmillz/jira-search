# Jira-Search

Table of Contents
=================

* [Jira-Search](#jira-search)
   * [About](#about)
   * [Docs](#docs)
   * [Installation](#installation)
   * [Usage](#usage)
      * [Example Usage](#example-usage)

## About

this is a releatively simple tool for quick searching jira issues via JQL queries and displaying the data in a nice easy to read table

## Docs

 * [jira-search](docs/jira-search.md)                         - Search Jira Issues via jql queries and others from the cli.
 * [jira-search search](docs/jira-search_search.md)           - Search Jira for issues matching text

## Installation

Download the binary from releases or you can install the program directly via 

```shell
go install github.com/mr-pmillz/jira-search@latest
```

create your config.yaml file

```shell
cp config.yaml.dist config.yaml
```

## Usage

```shell
Search Jira Issues via jql queries and others from
the cli.

Usage:
  jira-search [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  search      Search Jira for issues matching text

Flags:
      --config string   config file (default is $HOME/config.yaml)
  -h, --help            help for jira-search

Use "jira-search [command] --help" for more information about a command.
```

### Example Usage

Search for all issues assigned to you and in various status types

```shell
jira-search --config /path/to/your/config.yaml search --jql-raw-search 'assignee = "Babu Bott" AND status in ("Ready for Work", "In Progress")'
```

