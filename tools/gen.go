//go:generate swag init --generalInfo server.go --dir ../cmd/server/,../internal/api --parseDependency --output ../docs
//go:generate sh -c "wget -O ../web/webui/main.zip https://github.com/matetirpak/chessbot-playground-webui/archive/refs/heads/main.zip && unzip -o ../web/webui/main.zip -d ../web/webui/ && rm -f ../web/webui/main.zip && cp -r ../web/webui/chessbot-playground-webui-main/src ../web/webui/ && rm -rf ../web/webui/chessbot-playground-webui-main"

package tools
