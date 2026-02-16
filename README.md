# lrcsnc-lbl-client
A letter-by-letter client for [lrcsnc](https://github.com/Endg4meZer0/lrcsnc)

<!-- Insert a video example here later -->

This is a personal-oriented project: I made it work on my machine and that's my wish done (mostly, at least). I still will respond to issues and requests, but this is kind of a disclaimer for my inaction, I hope you understand.

# Build

Clone the repository and build it using Go v1.23 or newer.

```
git clone https://github.com/Endg4meZer0/lrcsnc-lbl-client.git
go build -o lrcsnc-lbl-client
```

# Usage

The client will check for the presence of `config.toml` in your user config directory. The default path for most users would be `$HOME/.config/lrcsnc-lbl-client/config.toml`; adapt the `$HOME/.config` part to your actual user config directory, if it differs.

An example of config file lies at `config_example.toml` with all the necessary comments.

Before the start of the client, the server (so, lrcsnc) should be already up and ready to accept connections. It usually isn't a problem when you start them at the same time (e.g., on device start-up), but better to keep it in mind.

For now I did not implement a reconnection mechanism, so if the connection drops for some reason, you'll need to restart the app.

<hr>

That should be all. Have fun!