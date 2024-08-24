<h1 align="center" id="title">Spogo</h1>

<p align="center"><img src="https://img.shields.io/github/go-mod/go-version/dionvu/spogo?style=for-the-badge" alt="shields"><img src="https://img.shields.io/github/commit-activity/m/dionvu/spogo?style=for-the-badge" alt="shields"><img src="https://img.shields.io/github/license/dionvu/spogo?style=for-the-badge" alt="shields"></p>

<p align="center" id="description">Spotify in the convenience of your command-line! Inspired by spotify-tui.</p>


<h2>🚀 Demo</h2>

[temp](temp)

<h2>🛠️ Installation Steps:</h2>

<p>1. Clone the repo.</p>

```
git clone https://github.com/dionvu/spogo && cd spogo
```

<p>2. Build.</p>

```
go build
```

<p>3. Move the binary.</p>

```
sudo mv spogo /usr/local/bin
```

<p>4. Clean up cloned files.</p>

```
cd .. && rm -rf spogo
```

<p>5. Run and follow configuration instructions</p>

```
spogo
```

<h2>⚙️ Configuration</h2>

<p>1. Navigate to <a href="https://developer.spotify.com/dashboard">Spotify Developer Dashboard</a> and click on "Create App".</p>

<p>2. Add a the following URI to "Redirect URIs" section, fill in the rest of the fields and press save.</p>

```
http://localhost:42069/callback
```

<p>3. Copy your Spotify Client ID and Client Secret for later use.</p>

<p>4. Navigate to Spogo configuration directory and open "config.yaml" in your prefered editor.</p>

```
cd ~ && cd .config/spogo
vim config.yaml
```

<p>5. Set your Spotify Client ID and Client Secret.</p>

```yaml
spotify:
  client_id: "YOUR_CLIENT_ID_HERE"
  client_secret: "YOUR_CLIENT_SECRET_HERE"
```

<h2>🎵 Spotify Devices</h2>

<p>Since Spogo is not a Spotify client and relies on the Spotify Web API, an external playback device is required.</p>

<p>It is recommended to use <a href="https://github.com/Spotifyd/spotifyd">Spotifyd</a>, but the offical Spotify client works too. (The offical web client also works but I wouldn't recommend it)</p>

<p>Select your playback device by running: </p>

```
spogo d
```

<p>You will only need to do this once since your device will be cached for your convenience.</p>
<h2>🧐 Commands & Features</h2>

Project's features:

*   Search and play your favorite albums tracks podcasts.
*   Manage multiple playback devices.
*   Control playback with simple commands.


<h2>❌ Uninstallation Steps:</h2>

<p>1. Remove binary.</p>

```
sudo rm /usr/local/bin/spogo
```

<p>2. Remove config files.</p>

```
rm -rf ~/.config/spogo
```

<p>3. Remove cache files.</p>

```
rm -rf ~/.cache/spogo
```
