# Chessbot Playground Server

**Chessbot Playground** is an interface developed to allow the graphical testing of chessbots including a server, API, web UI and illustrative bot.
This repository is the main entry point to the program and is responsible for running a server with an API as well as serving the web UI.

The other repositories can be found here:
[Chessbot Playground Bot](https://github.com/matetirpak/chessbot-playground-bot)
[Chessbot Playground WebUI](https://github.com/matetirpak/chessbot-playground-webui)

---

### Getting Started

#### Ubuntu

Make sure [go](https://go.dev/doc/install) is installed and copy the following line into ~/.bashrc or ~/.profile to make it available. 1.23.3 was used for testing.

```shell
export PATH=$PATH:/usr/local/go/bin:$HOME/go/bin
```

Clone the repository and install necessary dependencies: 

```bash
git clone https://github.com/matetirpak/chessbot-playground-server.git
cd chessbot-playground-server
go mod tidy
go install github.com/swaggo/swag/cmd/swag@latest
```

Use go generate to both fetch the web UI and generate the API documentation:

```bash
cd tools
go generate
```

To apply customizations to the [Chessbot Playground WebUI](https://github.com/matetirpak/chessbot-playground-webui) replace web/webui/src with the modified src/ directory.

Run the server:

```bash
go run cmd/server/server.go
```
or
```bash
go build -o bin/ ./cmd/server/...
./bin/server
```

Once started, it displays the ports it connects to.
API requests use port 8080, whilst the web UI connects to 8081 and can be opened by entering localhost:8081/ into the browser.

---

### API Documentation

The documentation of the API is located at docs/. Once the server runs, it provides a graphical interface at [localhost:8080/documentation/](http://localhost:8080/documentation/).
It is auto generated from code and can be recompiled with:

```bash
cd tools
go generate
```

Instructions to create an own bot as well as an example can be found at [Chessbot Playground Bot](https://github.com/matetirpak/chessbot-playground-bot).
