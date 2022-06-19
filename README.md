# magnets

For the Ebitenengine game jam https://itch.io/jam/ebiten-game-jam

## building WASM

```shell
GOOS=js GOARCH=wasm go build -o dist/magnets.wasm github.com/jakecoffman/cmd/magnets
cp wasm_exec.js index.html dist
```
