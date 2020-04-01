# sdrie
![loc](https://tokei.rs/b1/github/nektro/sdrie)
[![license](https://img.shields.io/github/license/nektro/sdrie.svg)](https://github.com/nektro/sdrie/blob/master/LICENSE)
[![paypal](https://img.shields.io/badge/donate-paypal-blue.svg?logo=paypal)](https://www.paypal.me/nektro)
[![discord](https://img.shields.io/discord/551971034593755159.svg)](https://discord.gg/P6Y4zQC)

An in-process key/value store for data with expiration dates in Go

## Installing
```
$ go get -u github.com/nektro/sdrie
```

## Usage
### `sdrie.New`
- `New() SdrieDataStore`
- `New` returns a new instance of a `SdrieDataStore`.

### `SdrieDataStore.Set`
- `Set(key string, value string, lifespan int64)`
- `Set` adds `value` to the data store associated to `key` and will survive for `lifespan` milliseconds.

### `SdrieDataStore.Get`
- `Get(key string) interface{}`
- `Get` retrieves the value associated to `key`, or `nil` otherwise.

### `SdrieDataStore.Has`
- `Has(key string) bool`
- `Has` returns a `bool` based on whether or not `key` exists in the data store. 

### `SdrieDataStore.Delete`
- `Delete(key string)`
- `Delete` retrieves the value associated to `key` or 'no-op' if key doesn't exist

## Contributing
We take issues all the time right here on GitHub. We use labels extensively to show the progress through the fixing process. Question issues are okay but make sure to close the issue when it's been answered!

[![issues](https://img.shields.io/github/issues/nektro/sdrie.svg)](https://github.com/nektro/sdrie/issues)

When making a pull request, please have it be associated with an issue and make a comment on the issue saying that you're working on it so everyone else knows what's going on :D

[![pulls](https://img.shields.io/github/issues-pr/nektro/sdrie.svg)](https://github.com/nektro/sdrie/pulls)

## Contact
- hello@nektro.net
- Meghan#2032 on discordapp.com
- @nektro on [twitter.com](https://twitter.com/nektro)

## License
MIT
