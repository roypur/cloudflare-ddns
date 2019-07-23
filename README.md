cloudflare-ddns

A simple ddns client for cloudflare

To build from source requires the go compiler. Precompiled binaries can be found under releases.

The releases are compiled for


linux:

- arm
- amd64
- i386
- mips64
- mips64le

freebsd:

- arm
- i386
- amd64

windows:

- amd64
- i386

darwin:

- amd64
- i386

Build instructions:

```bash
git clone https://github.com/roypur/cloudflare-ddns
cd cloudflare-ddns
./build
```

WARNING: This will compile for all architectures and may take a long time

The binaries will be in the bin folder.

To update your ddns run the compiled binary with path to config-file as argument

Example config-file:
```json
{
    "interval": 300,
    "api-key":"7f173c21601a601498726cc6bd66645c088a0",
    "api-email":"mail@example.org",
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
