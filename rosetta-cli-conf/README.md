This directory contains Configuration Files for rosetta-cli tool.
Use them to run rosetta-fantom tests:

```
cd rosetta-cli-conf/mainnet
rosetta-cli check:data --configuration-file config.json --asserter-configuration-file asserted_options.json
```

```
cd rosetta-cli-conf/testnet
rosetta-cli check:construction --configuration-file config.json --asserter-configuration-file asserted_options.json
```

For CI check:construction tests, use `testnet/config-prefunded.json` as a template for your config file and
pass there your testing prefunded account private key and address.

More info about rosetta-cli tool and its configuration files:
* [How to Test your Rosetta Implementation](https://www.rosetta-api.org/docs/rosetta_test.html)
* [rosetta-cli configuration file reference](https://www.rosetta-api.org/docs/rosetta_configuration_file.html)
