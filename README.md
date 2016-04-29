# cloudflare-ddns
A simple ddns client for cloudflare

To use it run

    java -jar ddns.jar config.json
    
Example config-file:

    {
        "apiKey":"7f173c21601a601498726cc6bd66645c088a0",
        "apiEmail":"mail@example.org",
        "domain":"example.org",
        "v4Host":"home",
        "ipv6":{
            "first":"AA:BB:CC:BB:EE:AB",
            "second":"BA:44:AA:CC:AE:AB"
        }
    }
