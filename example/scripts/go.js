document.addEventListener("DOMContentLoaded", () => {
    const go = new Go();

    const canvas = document.getElementById("canvas");
    const ctx = canvas.getContext("2d");
    canvas.width = 960;
    canvas.height = 640;

    const imageData = ctx.createImageData(960, 640);
    const buf = new Uint8ClampedArray(960 * 640 * 4);

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

    window.addEventListener("pointerup", () => {
        fireworksPointerUp();
    });

    window.addEventListener("pointercancel", () => {
        fireworksPointerUp();
    });

    function loop() {
        fireworksTick(buf);
        imageData.data.set(buf);
        ctx.putImageData(imageData, 0, 0);
        requestAnimationFrame(loop);
    }

    WebAssembly.instantiateStreaming(fetch("scripts/fireworx.wasm"), go.importObject).then(r => {
        void go.run(r.instance);
        setTimeout(() => {
            requestAnimationFrame(loop);
        }, 100);
    });
});