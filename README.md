# Cloak/Server Client Configuration

This repository provides an example configuration for the Cloak client/server with placeholders where sensitive credentials or specific values need to be replaced. Make sure to replace all placeholders with your actual credentials or required details before using the configuration.

## Configuration Example

Below is a sample configuration for the Cloak client:

```json
{
  "Transport": "CDN",
  "ProxyMethod": "shadowsocks",
  "EncryptionMethod": "plain",
  "UID": "<your-UID-here>",
  "PublicKey": "<your-public-key-here>",
  "ServerName": "<your-server-name-here>",
  "NumConn": 8,
  "BrowserSig": "chrome",
  "StreamTimeout": 300,
  "RemoteHost": "<your-remote-host-here>",
  "RemotePort": "<your-remote-port-here>",
  "CDNWsUrlPath": "<your-cdn-ws-url-path-here>",
  "CDNOriginHost": "<your-cdn-origin-host-here>"
}
```

**ServerName** and **CDNWsUrlPath** have the same value, in particular, the domain name.

Below is a sample configuration for the Cloak server:

```json
{
  "ProxyBook": {
    "shadowsocks": [
      "tcp",
      "127.0.0.1:<--keys-port of Outline>"
    ]
      },
  "BindAddr": [
    "127.0.0.1:<cloak-server port>"
    ],
  "BypassUID": [
    "<user-UID>"
  ],
  "AdminUID": "<admin-UID>",
  "RedirAddr": "<domain-name>",
  "PrivateKey": "<private-key>",
  "StreamTimeout": 300
}
```

Installation script: **install_script.sh**. It works in an interactive mode without having to enter params via cli keys.


