This directory contains Configuration Files for rosetta-cli tool.
Use the to run rosetta-fantom tests:

```
cd rosetta-cli-conf/mainnet
rosetta-cli check:data --configuration-file config.json --asserter-configuration-file asserted_options.json
```

```
cd rosetta-cli-conf/testnet
rosetta-cli check:construction --configuration-file config.json --asserter-configuration-file asserted_options.json
```

More info about rosetta-cli tool and its configuration files:
* [How to Write a Configuration File for rosetta-cli Testing](https://www.rosetta-api.org/docs/rosetta_configuration_file.html)
