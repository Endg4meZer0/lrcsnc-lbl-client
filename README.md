# lrcsnc-lbl-client
A letter-by-letter client for [lrcsnc](https://github.com/Endg4meZer0/lrcsnc) with ability to show some song info letter-by-letter as well.

https://github.com/user-attachments/assets/2cf266b2-e7bb-4dc8-ada8-78d367a25e74

# Build

Clone the repository and build it using Go v1.23 or newer.

```
git clone https://github.com/Endg4meZer0/lrcsnc-lbl-client.git
go build -o lrcsnc-lbl-client
```

# Usage

The client will check for the presence of `config.toml` in the following path: `$XDG_USER_CONFIG/lrcsnc-lbl-client/config.toml`. The default path for most users would look like `$HOME/.config/lrcsnc-lbl-client/config.toml`.

An example of config file can be found in `config_example.toml` with all the describing comments.

Before the start of the client, the server (so, [lrcsnc](https://github.com/Endg4meZer0/lrcsnc)) should be already up and ready to accept connections. You might want to add `sleep 1` to the client's command if you notice connection issues on start-up, though personally I've never encountered problems on my machines.

For now there's no reconnection mechanism, so if the connection drops for some reason, you'll need to restart the client after restarting the server. This may be changed in later versions.

<hr>

That should be all. Have fun!
