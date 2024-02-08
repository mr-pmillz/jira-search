# Jira-Search

[![Go Report Card](https://goreportcard.com/badge/github.com/mr-pmillz/jira-search)](https://goreportcard.com/report/github.com/mr-pmillz/jira-search)
![GitHub all releases](https://img.shields.io/github/downloads/mr-pmillz/jira-search/total?style=social)
![GitHub repo size](https://img.shields.io/github/repo-size/mr-pmillz/jira-search?style=plastic)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mr-pmillz/jira-search?style=plastic)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/mr-pmillz/jira-search?style=plastic)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/mr-pmillz/jira-search?style=plastic)
[![Twitter](https://img.shields.io/twitter/url?style=social&url=https%3A%2F%2Fgithub.com%2Fmr-pmillz%2Fjira-search)](https://twitter.com/intent/tweet?text=Wow:&url=https%3A%2F%2Fgithub.com%2Fmr-pmillz%2Fjira-search)
[![CI](https://github.com/mr-pmillz/jira-search/actions/workflows/ci.yml/badge.svg)](https://github.com/mr-pmillz/jira-search/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/mr-pmillz/jira-search/branch/master/graph/badge.svg?token=L7GRPOPHCL)](https://codecov.io/gh/mr-pmillz/jira-search)

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

create yourself an api key in your jira account profile

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

Make an alias

```shell
alias mywork="jira-search --config /path/to/your/config.yaml search --jql-raw-search 'assignee = \"Babu Bott\" AND status in (\"Ready for Work\", \"In Progress\")'"
```

