# corsproxy

Debug utility for CORS-enabled servers

## Usage

**Install:**

```shell
go get -u github.com/goware/corsproxy
```

**Terminal:**

```shell
corsproxy -source=https://remotehost.com/ -listen=1337
```

You can now send requests to http://localhost:1337 which will proxy to https://remotehost.com/
and circumvent CORS. This can be useful when developing against a remote service which uses CORS.


## LICENSE

MIT
