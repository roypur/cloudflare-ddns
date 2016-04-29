# cloudflare-ddns
A simple ddns client for cloudflare

Build instructions

    git clone https://github.com/roypur/cloudflare-ddns
    cd cloudflare-ddns
    mvn package

The packaged file will be under target/ddns-version.jar

To use it run

    java -jar ddns.jar config.json
    
Example config-file:

    {
        "apiKey":"7f173c21601a601498726cc6bd66645c088a0",
        "apiEmail":"mail@example.org",
        "domain":"example.org",
        "v4Host":"home",
        "validateCerts":true,
        "ipv6":{
            "first":"AA:BB:CC:BB:EE:AB",
            "second":"BA:44:AA:CC:AE:AB"
        }
    }
