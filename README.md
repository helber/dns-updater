
DNS-UPDATER
-----------

This is a simple application to change DNS "A" Record on cloudflare.
This application can act as dynamic dns inside your home network.
To run you need export some variables:

`A_HOST`: Is o host name of "A" record
`CF_API_KEY`: cloudflare API key
`CF_API_EMAIL`: cloudflare e-mail user ("login")

You awso need a valid domain that you can handle.

You can create a free account [here-https://www.cloudflare.com/](https://www.cloudflare.com/)

```bash
export A_HOST="www.example.com"
export CF_API_KEY="XXXXXOOOOXXXXOOOO"
export CF_API_EMAIL="jjjj@example.com"
```

