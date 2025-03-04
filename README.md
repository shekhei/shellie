# What is this?

This is a little tool that runs as a shell assistant, currently only supports bash.

## How to run?

Currently you need to install

- make
- go
- protoc

```bash
git clone ...
make
# Add this to your bash profile
source [clonefolder]/shellie.bash
# run in the background, or a separate terminal
[clonefolder]/bin/server

# either ctrl+k to autosuggest
```

## What does this send?

Your historical commands, current working directory

You can take a look at the bash script, or the ./pb/service.proto

## What next?

- Additional prompt
- partial commands
- maybe also adding more context for your tools
- Also wondering if some outputs, at least if the return code will be useful for the llm to suggest the next token
- zsh support
- more shell support
- "proper" shell support?

## Why a local service?

We want to reuse the http connection if we are talking to a service.
We also want to be able to perhaps do some caching of the results in the future.
