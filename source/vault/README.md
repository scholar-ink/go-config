# Vault Source

The vault source reads config from vault secrets

## Vault Format

The vault source expects keys under the default prefix `/secret/micro/config`

Values are expected to be json

```
// set database
vault kv put secret/micro/config/database value='{"address": "10.0.0.1", "port": 3306}'
// set cache
vault kv put secret/micro/config/cache value='{"address": "10.0.0.2", "port": 6379}'
```

Keys are split on `/` so access becomes

```
conf.Get("secret", "micro", "config", "database")
```

## New Source

Specify source with data

```go
vaultSource := consul.NewSource(
	// optionally specify vault address; default to localhost:8500
	vault.WithAddress("10.0.0.10:8500"),
	// optionally specify prefix; defaults to secret/micro/config
	vault.WithPrefix("/my/prefix"),
	// optionally strip the provided prefix from the keys, defaults to false
	vault.StripPrefix(true),
)
```

## Load Source

Load the source into config

```go
// Create new config
conf := config.NewConfig()

// Load file source
conf.Load(vaultSource)
```
