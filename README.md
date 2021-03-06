## cloudflare-ddns

A simple ddns client for cloudflare
To build from source requires the go compiler

Build instructions:

```bash
git clone https://github.com/roypur/cloudflare-ddns
cd cloudflare-ddns
./build.sh
```

How to run:

```bash
bin/linux_amd64 config.json
```

Example config file:

```json
{
    "interval": 300,
    "timeout": 1,
    "token": "ef3a5f32a4ebd99eca390469a68b25a199d6e924",
    "domain":"example.org",
    "ipv4": {
        "imap": {"local": false},
        "smtp": {"local": false}
    },
    "ipv6": {
        "first": {
            "local": true,
            "addr":"AA:BB:CC:BB:EE:AB",
            "prefix-length":48,
            "host-prefix-length":64,
            "prefix-id":"1",
            "ismac":true
        },
        "second": {
            "local": false,
            "addr":"::1",
            "prefix-length":48,
            "host-prefix-length":64,
            "prefix-id":"1",
            "ismac":false
        }
    }
}
```
