//go:build js && wasm

package main

import (
	"fireworx/render"
	"fireworx/sim"
	"math/rand/v2"
	"syscall/js"
)

type Game struct {
	fw       *sim.Fireworks
	renderer *render.Renderer
	frame    []byte
	autoTick int
}

var game *Game

func tick(this js.Value, args []js.Value) any {
	const dt = 1.0 / 60.0
	game.fw.Update(dt)

	game.autoTick++
	if game.autoTick > 90 {
		game.autoTick = 0
		effect := sim.EffectID(rand.IntN(len(sim.Catalog)))
		x := render.Unproject(float32(rand.IntN(render.Width)), 0).X
		game.fw.LaunchFromGround(effect, x)
	}

	game.renderer.Draw(game.fw, game.frame)
	js.CopyBytesToJS(args[0], game.frame)
	return nil
}

func click(this js.Value, args []js.Value) any {
	x := float32(args[0].Float())
	y := float32(args[1].Float())
	pos := render.Unproject(x, y)
	game.fw.Burst(sim.EffectID(rand.IntN(len(sim.Catalog))), pos)
	return nil
}

func main() {
	world := sim.DefaultWorld()
	game = &Game{
		fw:       sim.NewFireWorks(sim.Catalog, world, 42),
		renderer: render.New(),
		frame:    make([]byte, render.Width*render.Height*4),
	}

	js.Global().Set("fireworksTick", js.FuncOf(tick))
	js.Global().Set("fireworksClick", js.FuncOf(click))
	select {}
}
