SHELL := pwsh.exe

clean:
	@echo "Cleaning up..."
	@pwsh -Command "if(Test-Path -Path ./tmp) { rm -recurse -force ./tmp }"

build:
	@echo "Building web app..."
	@pwsh -Command "Set-Item Env:GOARCH wasm; Set-Item Env:GOOS js; go build -o ./tmp/web/app.wasm ./cmd/webapp"
	@echo "Building server..."
	@pwsh -Command "Set-Item Env:GOARCH amd64; Set-Item Env:GOOS windows; go build -o ./tmp/server.exe ./cmd/webapp/"
	@echo "Copying web files..."
	@pwsh -Command "copy ./web/* ./tmp/web/ -Recurse -Force"

run: build
	@echo "Running server..."
	cd ./tmp && ./server.exe
