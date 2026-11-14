//go:build js && wasm

package audio

import (
	"math"
	"math/rand/v2"
	"syscall/js"
)

var (
	ctx        js.Value
	masterGain js.Value
	ready      bool
	muted      bool
)

func ToggleAudio(m bool) {
	muted = m

	if !ready {
		return
	}

	now := ctx.Get("currentTime").Float()
	gain := masterGain.Get("gain")

	if m {
		gain.Call("linearRampToValueAtTime", 0.0001, now+0.01)
	} else {
		gain.Call("linearRampToValueAtTime", 1, now+0.01)
	}
}

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

	masterGain = ctx.Call("createGain")
	masterGain.Get("gain").Set("value", 1)
	masterGain.Call("connect", ctx.Get("destination"))

	ready = true
}

func PlayBoom(intensity float32) {
	if !ready || muted {
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
	gain.Call("connect", masterGain)
	src.Call("start", now)
}

func PlayLaunchNoise(duration float32) {
	if !ready || muted || duration <= 0 {
		return
	}

	now := ctx.Get("currentTime").Float()
	dur := float64(duration)
	dest := masterGain

	thump := ctx.Call("createOscillator")
	thump.Set("type", "sine")
	tf := thump.Get("frequency")
	tf.Call("setValueAtTime", 90, now)
	tf.Call("exponentialRampToValueAtTime", 40, now+0.12)

	thumpGain := ctx.Call("createGain")
	thg := thumpGain.Get("gain")
	thg.Call("setValueAtTime", 0.4, now)
	thg.Call("exponentialRampToValueAtTime", 0.0001, now+0.15)

	thump.Call("connect", thumpGain)
	thumpGain.Call("connect", dest)
	thump.Call("start", now)
	thump.Call("stop", now+0.17)

	noiseFrames := int(44100 * dur)
	noiseBuf := ctx.Call("createBuffer", 1, noiseFrames, 44100)
	noiseData := noiseBuf.Call("getChannelData", 0)
	for i := 0; i < noiseFrames; i++ {
		noiseData.SetIndex(i, rand.Float64()*2-1)
	}
	noiseSrc := ctx.Call("createBufferSource")
	noiseSrc.Set("buffer", noiseBuf)

	hissFilter := ctx.Call("createBiquadFilter")
	hissFilter.Set("type", "bandpass")
	hf := hissFilter.Get("frequency")
	hf.Call("setValueAtTime", 400, now)
	hf.Call("exponentialRampToValueAtTime", 2500, now+dur*0.85)
	hissFilter.Get("Q").Call("setValueAtTime", 0.8, now)

	hissGain := ctx.Call("createGain")
	hg := hissGain.Get("gain")
	hg.Call("setValueAtTime", 0.0001, now)
	hg.Call("linearRampToValueAtTime", 0.08, now+dur*0.7)
	hg.Call("linearRampToValueAtTime", 0.0001, now+dur)

	noiseSrc.Call("connect", hissFilter)
	hissFilter.Call("connect", hissGain)
	hissGain.Call("connect", dest)
	noiseSrc.Call("start", now)
}

func IntensityFromCount(count uint16) float32 {
	v := float32(count) / 210
	if v > 1 {
		v = 1
	}
	return v
}
