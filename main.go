//go:build js && wasm

package main

import (
	"fireworx/audio"
	"fireworx/render"
	"fireworx/sim"
	vm "fireworx/vec"
	"math/rand/v2"
	"syscall/js"
)

type Game struct {
	fw        *sim.Fireworks
	renderer  *render.Renderer
	frame     []byte
	autoTick  int
	dragging  bool
	dragX     float32
	paintTick int
}

var game *Game

const paintIntervalTicks = 6

func tick(_ js.Value, args []js.Value) any {
	const dt = 1.0 / 60.0
	game.fw.Update(dt)

	game.autoTick++
	if game.autoTick > 120 {
		game.autoTick = 0
		effect := sim.EffectID(rand.IntN(len(sim.Catalog)))
		x := render.Unproject(float32(rand.IntN(render.Width)), 0).X
		game.fw.LaunchFromGround(effect, x)
	}

	if game.dragging {
		game.paintTick++
		if game.paintTick >= paintIntervalTicks {
			game.paintTick = 0
			pos := render.Unproject(game.dragX, 0)
			effect := sim.EffectID(rand.IntN(len(sim.Catalog)))
			game.fw.LaunchFromGround(effect, pos.X)
		}
	}

	game.renderer.Draw(game.fw, game.frame)
	js.CopyBytesToJS(args[0], game.frame)
	return nil
}

func pointerDown(_ js.Value, args []js.Value) any {
	audio.Init()
	x := float32(args[0].Float())
	game.dragging = true
	game.dragX = x
	game.paintTick = 100 //immediate launch on press
	return nil
}

func pointerMove(_ js.Value, args []js.Value) any {
	if !game.dragging {
		return nil
	}
	game.dragX = float32(args[0].Float())
	return nil
}

func pointerUp(_ js.Value, _ []js.Value) any {
	game.dragging = false
	return nil
}

func toggleAudio(_ js.Value, args []js.Value) any {
	audio.ToggleAudio(args[0].Bool())
	return nil
}

func main() {
	world := sim.DefaultWorld()
	fw := sim.NewFireWorks(sim.Catalog, world, 42)
	fw.OnLaunch = func(_ vm.Vec3, _ sim.EffectID, duration float32) {
		audio.PlayLaunchNoise(duration)
	}
	fw.OnBurst = func(_ vm.Vec3, _ sim.EffectID, count uint16) {
		audio.PlayBoom(audio.IntensityFromCount(count))
	}

	game = &Game{
		fw:       fw,
		renderer: render.New(),
		frame:    make([]byte, render.Width*render.Height*4),
	}

	js.Global().Set("fireworksTick", js.FuncOf(tick))
	js.Global().Set("fireworksPointerDown", js.FuncOf(pointerDown))
	js.Global().Set("fireworksPointerMove", js.FuncOf(pointerMove))
	js.Global().Set("fireworksPointerUp", js.FuncOf(pointerUp))
	js.Global().Set("toggleAudio", js.FuncOf(toggleAudio))
	select {}
}
