# README

## What is q
q - is the simple cli tool build to interact with open ai like model from the terminal.

## How to configure
Create a config file in `.config/q/config.json`

```json
{
    "baseUrl": "http://192.168.1.126:1234/v1",
    "model":   "qwen/qwen3-vl-30b"
}
```

## How to use
### 1. Pass question as first argument
`q "How to list all docker containers?"`

### 2. Pass with -r
`q -r "How to list all docker containers?"`

### 3. Pass question as args (does not supports special symbols)
`q How to list all docker containers`

