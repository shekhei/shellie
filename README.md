# What is this?

This is a little tool that runs as a shell assistant, currently only supports bash.

## How to start?

```bash
# 1. Download the right package from github release tab
# 2. Unzip to some folder, recommend $HOME/.shellie/
# 3. run the server, it will also generate a ~/.config/shellie.toml for the first time and exit
# 4. modify ~/.config/shellie.toml according to your own environment
# 5. start server in the background
# 6. Configure your shell, see "Shell configuration" section below
```

## How to use?

Bash - After setting up the shell configurations(see below), hit Ctrl+K for an autosuggestion

## Shell configurations

### Bash settings

Add this to your ~/.bashrc

```
source $HOME/.shellie/shellie.bash # replace the $HOME/.shellie with your unzip path
```

### Zsh (Coming soon)

## Setting up shellie server

### Sample ~/.config/shellie.toml for OpenAI

```toml
# ...
[service]
chat_completion_endpoint = "https://api.openai.com/v1/chat/completions"
api_key = "..."
# o3-mini, o1-mini, o1, gpt-4o, etc, see https://platform.openai.com/docs/guides/text-generation
model = "o1-mini"
```

### Sample ~/.config/shellie.toml for Gemini-Flash

```toml
# ...
[service]
chat_completion_endpoint = "https://generativelanguage.googleapis.com/v1beta/openai/"
api_key = "..."
model = "gemini-2.0-flash"
```

### Setting up with Ollama

Ollama uses openai compatible endpoints, you can just replace the "chat_completion_endpoint" with the ollama endpoint.

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
