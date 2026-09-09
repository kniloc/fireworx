# fireworx

A Go port of [asrael's fireworks-rs](https://github.com/asrael/fireworks-rs), a multi-stage fireworks particle simulator that runs in the browser via WebAssembly.

## Overview

fireworx simulates fireworks physics and renders them to an RGBA pixel buffer, which a host HTML page draws to a canvas. Fireworks launch automatically every 90 frames and also respond to clicks, bursting at wherever you click on screen.

The codebase is organized into three packages:

- **`sim`**: physics simulation: particle lifecycle, multi-stage effects, and a catalog of firework types
- **`render`**: rasterizer that draws the simulation state into a flat byte buffer
- **`vec`**: 2D vector math used by the other packages

## Requirements

- Go 1.27 or later (with WebAssembly support)
- A browser that supports WebAssembly

## Building

Compile to WebAssembly using the `js/wasm` target:

```sh
GOOS=js GOARCH=wasm go build -o fireworx.wasm .
```

You'll also need the Go WebAssembly runtime shim in your HTML directory:

```sh
cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" .
```

## Usage

In your HTML, load the WASM runtime and the compiled binary, then drive the simulation by calling the two exported functions each frame:

```html
<script src="wasm_exec.js"></script>
<script>
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("fireworx.wasm"), go.importObject)
    .then(result => go.run(result.instance));
</script>
```

**Exported JS functions:**

- `fireworksTick(frameBuffer)`: advance the simulation by one frame (1/60 s) and render into `frameBuffer`, a `Uint8Array` of size `width × height × 4` (RGBA)
- `fireworksClick(x, y)`: burst a random firework effect at screen coordinates `(x, y)`

A minimal render loop:

```js
const canvas = document.getElementById("canvas");
const ctx = canvas.getContext("2d");
const buf = new Uint8Array(canvas.width * canvas.height * 4);

function frame() {
  fireworksTick(buf);
  const imageData = new ImageData(new Uint8ClampedArray(buf.buffer), canvas.width, canvas.height);
  ctx.putImageData(imageData, 0, 0);
  requestAnimationFrame(frame);
}

canvas.addEventListener("click", e => {
  const rect = canvas.getBoundingClientRect();
  fireworksClick(e.clientX - rect.left, e.clientY - rect.top);
});

requestAnimationFrame(frame);
```

## Credits

Based on [fireworks-rs](https://github.com/asrael/fireworks-rs) by [asrael](https://github.com/asrael).

## License

MIT. See [LICENSE](LICENSE) for details.
