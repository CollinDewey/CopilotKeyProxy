# CopilotKeyProxy

Reverse Proxy to translate an API key to the token that Github expects

Get your API key by running getKey.go and following the instructions

```
go run getKey.go
```

Run the server by running proxy.go, optionally with a listening address, and pointing your OpenAI compatible client towards it

```
go run proxy.go
go run proxy.go :12345
```
