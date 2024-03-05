# totpgen
Utility to generate TOTP tokens

## Usage

The secret tokens are saved in the systems keychain (via [99designs/go-keychain](https://github.com/99designs/go-keychain)).

To add a totp TOTP secret, run
```sh
totpgen set "<secret-name>" "totp-secret"
# eg
totpgen set google "SAucYHYJyfma1Fa6uFlBqzUluusgIj1slSwKRoVvhGYZsVCt"
totpgen set aws "44uHJtA8IwpKy9JjaaprSizgZ2TSImDY8iUPvm1qaDHReOTJ"
```

To generate the current OPT code run
```sh
totpgen google
# output:
123456
```

You can also create a symlink to the script and name it after `totpgen-<name>` (e.g. `totpgen-google`).
Invoking the tool like this will also print the TOTP token for the specified name, no arguments required.

```sh
ln -sT totpgen totp-google
totp-google
# output:
123456
```

There are a couple more commands:
```sh
# show the names of saved totp configuration
totpgen list
# output:
google
aws

# to remove a secret use the 'set' command with an empty secret
totpgen set google ""
```

## Installation
Via `go install`:
```sh
go install github.com/floj/totpgen
~/go/bin/totpgen --help
```

Manual
```sh
git clone https://github.com/floj/totpgen.git
cd totpgen
./build.sh
./totpgen --help
```

### MacOS
I don't provide precompiled binaries, because past expirence showed that cross-compiled binaries for Mac do not properly work with the OSX keychain. Thus, if you want to use it on Mac, you need to compile it yourself using one of the above command.

#### Additionally available config options for MacOS

| Environment variable  | Configuration | Example |
|-----------------------|---------------|---------|
| `TOTPGEN_KEYCHAIN_NAME` | Name of the Keychain files used | `TOTPGEN_KEYCHAIN_NAME=totpgen-secrets` |


## Why?
Main motivation was to use it in [aws-vault](https://github.com/99designs/aws-vault). AWS Vault supports creating TOTP tokens via [pass-otp](https://github.com/tadfisher/pass-otp). This is very nice, but limits you to use `pass`. I created a `scriptmfa` prompt provider (see [genericscript.go](https://github.com/floj/aws-vault/blob/master/prompt/genericscript.go)) that is able to call whatever script you want. Just point it to `totpgen-aws` by setting `AWS_VAULT_MFA_SCRIPT=totpgen-aws` and you are good to go.
