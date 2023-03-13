# logseq2obsidian 

This is a tiny program for converting pages in Logseq to a vault of Obsidian. Please use it at your own risk.

## Usage

Download the source code with golang environment, and build the program

```shell
go build
```

Input in the prompt 

```shell
./logseq2obsidian {logseq-page-dir} -o {target-dir}
```

## Features 

- Automatically arrange Logseq namespace to file hierarchy 
- Automatic linebreaks for blocks in Logseq