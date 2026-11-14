//go:build js && wasm

package audio

import (
	"math"
	"math/rand/v2"
	"syscall/js"
)

var (
	ctx   js.Value
	ready bool
)

func Init() {
	if ready {
		return
	}

	audioContextClass := js.Global().Get("AudioContext")
	if audioContextClass.IsUndefined() {
		audioContextClass = js.Global().Get("webkitAudioContext") // safari...
	}

	if audioContextClass.IsUndefined() {
		return // no Web Audio available
	}

	ctx = audioContextClass.New()
	ready = true
}

func PlayBoom(intensity float32) {
	if !ready {
		return
	}

	if intensity < 0.05 {
		intensity = 0.05
	}

	if intensity > 1 {
		intensity = 1
	}

	now := ctx.Get("currentTime").Float()
	const sampleRate = 44100
	dur := 0.25 + float64(intensity)*0.35
	frames := int(float64(sampleRate) * dur)

	buf := ctx.Call("createBuffer", 1, frames, sampleRate)
	data := buf.Call("getChannelData", 0)
	for i := 0; i < frames; i++ {
		t := float64(i) / float64(frames)
		envelope := math.Exp(-6 * t)
		data.SetIndex(i, (rand.Float64()*2-1)*envelope*float64(intensity))
	}

	filter := ctx.Call("createBiquadFilter")
	filter.Set("type", "lowpass")
	filter.Get("frequency").Call("setValueAtTime", 1800-1200*float64(intensity), now)

	src := ctx.Call("createBufferSource")
	src.Set("buffer", buf)

	gain := ctx.Call("createGain")
	gain.Get("gain").Call("setValueAtTime", float64(intensity), now)

	src.Call("connect", filter)
	filter.Call("connect", gain)
	gain.Call("connect", ctx.Get("destination"))
	src.Call("start", now)
}

func PlayWhistle(duration float32) {
	if !ready || duration <= 0 {
		return
	}

	now := ctx.Get("currentTime").Float()
	dur := float64(duration)

	osc1 := ctx.Call("createOscillator")
	osc1.Set("type", "sine")
	f1 := osc1.Get("frequency")
	f1.Call("setValueAtTime", 900, now)
	f1.Call("linearRampToValueAtTime", 1500, now+dur*0.85)

	osc2 := ctx.Call("createOscillator")
	osc2.Set("type", "sine")
	f2 := osc2.Get("frequency")
	f2.Call("setValueAtTime", 915, now) // slight dissonance
	f2.Call("linearRampToValueAtTime", 1520, now+dur*0.85)

	wobble := ctx.Call("createOscillator")
	wobble.Set("type", "sine")
	wobble.Get("frequency").Call("setValueAtTime", 9, now)
	wobbleGain := ctx.Call("createGain")
	wobbleGain.Get("gain").Call("setValueAtTime", 15, now) // ± 15Hz flutter
	wobble.Call("connect", wobbleGain)
	wobbleGain.Call("connect", f1)
	wobbleGain.Call("connect", f2)

	toneGain := ctx.Call("createGain")
	tg := toneGain.Get("gain")
	tg.Call("setValueAtTime", 0.0001, now)
	tg.Call("linearRampToValueAtTime", 0.12, now+0.08)
	tg.Call("setValueAtTime", 0.12, now+dur*0.85)
	tg.Call("linearRampToValueAtTime", 0.0001, now+dur)

	noiseFrames := int(44100 * dur)
	noiseBuf := ctx.Call("createBuffer", 1, noiseFrames, 44100)
	noiseData := noiseBuf.Call("getChannelData", 0)
	for i := 0; i < noiseFrames; i++ {
		noiseData.SetIndex(i, rand.Float64()*2-1)
	}
	noiseSrc := ctx.Call("createBufferSource")
	noiseSrc.Set("buffer", noiseBuf)

	noiseFilter := ctx.Call("createBiquadFilter")
	noiseFilter.Set("type", "bandpass")
	nf := noiseFilter.Get("frequency")
	nf.Call("setValueAtTime", 1200, now)
	nf.Call("linearRampToValueAtTime", 2000, now+dur*0.85)
	noiseFilter.Get("Q").Call("setValueAtTime", 1.2, now)

	noiseGain := ctx.Call("createGain")
	ng := noiseGain.Get("gain")
	ng.Call("setValueAtTime", 0.0001, now)
	ng.Call("linearRampToValueAtTime", 0.05, now+0.08)
	ng.Call("setValueAtTime", 0.05, now+dur*0.85)
	ng.Call("linearRampToValueAtTime", 0.0001, now+dur)

	osc1.Call("connect", toneGain)
	osc2.Call("connect", toneGain)
	noiseSrc.Call("connect", noiseFilter)
	noiseFilter.Call("connect", noiseGain)

	toneGain.Call("connect", ctx.Get("destination"))
	noiseGain.Call("connect", ctx.Get("destination"))

	osc1.Call("start", now)
	osc2.Call("start", now)
	wobble.Call("start", now)
	noiseSrc.Call("start", now)

	osc1.Call("stop", now+dur+0.02)
	osc2.Call("stop", now+dur+0.02)
	wobble.Call("stop", now+dur+0.02)
}

func IntensityFromCount(count uint16) float32 {
	v := float32(count) / 210
	if v > 1 {
		v = 1
	}
	return v
}
