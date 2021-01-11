# totpgen
Utility to generate TOTP tokens

## Usage

First create the config file `.config/totpgen/config.json` and add entires like
```
[
    {
        "name": "aws",
        "secret": "55XX..."
    },
    {
        "name": "github",
        "secret": "55XX..."
    }
]
```

Then invoke `totpgen <name>`. The tool prints out a TOTP token.

You can also create a symlink to the script and name it after `totpgen-<name>` (e.g. `totpgen-aws`).
Invoking the tool like this will also print the TOTP token for the specified name, no arguments required.

## Why
Main motivation was to use it in [aws-vault](https://github.com/99designs/aws-vault). AWS Vault supports creating TOTP tokens via [pass-otp](https://github.com/tadfisher/pass-otp). This is very nice, but limits you to use pass. I created a `scriptotp` (see [scriptotp.go](https://github.com/floj/aws-vault/blob/master/prompt/scriptotp.go)) prompt provider that is able to call whatever script you want. Just point it to `totpgen-aws` by setting `AWS_VAULT_MFA_SCRIPT=totpgen-aws` and you are good to go.