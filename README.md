# fireworx

A Go port of [asrael's fireworks-rs](https://github.com/asrael/fireworks-rs), a multi-stage fireworks particle simulator that runs in the browser via WebAssembly.

## Overview

fireworx simulates fireworks physics and renders them to an RGBA pixel buffer, which a host HTML page draws to a canvas. Fireworks launch automatically about every two seconds and also respond to pointer presses and drags, bursting wherever you click or paint across the sky.

The codebase is organized into three packages:

- **`sim`**: physics simulation: particle lifecycle, multi-stage effects, and a catalog of 13 firework types (Brocade, Chrysanthemum, Comet, Crossette, Dragon's Eggs, Fish, Palm, Peony, Pistil, Ring, Strobe, Tourbillon, Willow)
- **`render`**: rasterizer that draws the simulation state into a flat byte buffer, including a procedurally generated city skyline
- **`vec`**: 3D vector math used by the other packages

A working browser demo lives in [`example/`](example/), including the compiled WASM binary and JS glue code.

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

See [`example/`](example/) for a complete reference setup (`index.html`, `scripts/go.js`, `scripts/fireworx.css`, and a prebuilt `fireworx.wasm`).

## Usage

In your HTML, load the WASM runtime and the compiled binary, then drive the simulation by calling the exported functions each frame:

```html
<script src="wasm_exec.js"></script>
<script>
  const go = new Go();
  WebAssembly.instantiateStreaming(fetch("fireworx.wasm"), go.importObject)
    .then(result => go.run(result.instance));
</script>
```

**Exported JS functions:**

- `fireworksTick(frameBuffer)`: advance the simulation by one frame (1/60 s) and render into `frameBuffer`, a `Uint8Array`/`Uint8ClampedArray` of size `width × height × 4` (RGBA)
- `fireworksPointerDown(x, y)`: begin a launch-on-press-and-drag interaction at canvas coordinates `(x, y)`
- `fireworksPointerMove(x, y)`: update the drag position while the pointer is held down; fireworks launch periodically along the drag path
- `fireworksPointerUp()`: end the drag interaction

The canvas is `960 × 640` pixels by default.

A minimal render loop:

```js
const canvas = document.getElementById("canvas");
canvas.width = 960;
canvas.height = 640;
const ctx = canvas.getContext("2d");
const imageData = ctx.createImageData(canvas.width, canvas.height);
const buf = new Uint8ClampedArray(canvas.width * canvas.height * 4);

function frame() {
  fireworksTick(buf);
  imageData.data.set(buf);
  ctx.putImageData(imageData, 0, 0);
  requestAnimationFrame(frame);
}

function canvasCoords(e) {
  const rect = canvas.getBoundingClientRect();
  return [e.clientX - rect.left, e.clientY - rect.top];
}

canvas.addEventListener("pointerdown", e => {
  const [x, y] = canvasCoords(e);
  fireworksPointerDown(x, y);
});

canvas.addEventListener("pointermove", e => {
  const [x, y] = canvasCoords(e);
  fireworksPointerMove(x, y);
});

window.addEventListener("pointerup", () => fireworksPointerUp());
window.addEventListener("pointercancel", () => fireworksPointerUp());

requestAnimationFrame(frame);
```

## Credits

Based on [fireworks-rs](https://github.com/asrael/fireworks-rs) by [asrael](https://github.com/asrael).

## License

MIT. See [LICENSE](LICENSE) for details.
</content>
