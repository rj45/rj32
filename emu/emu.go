package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/rj45/rj32/emu/anim"
	"github.com/rj45/rj32/emu/data"
	"github.com/rj45/rj32/emu/rj32"
	"github.com/rj45/rj32/emu/vdp"
)

const (
	spriteMemSize = 8
	numSprites    = 1 << spriteMemSize
	numSheets     = 16
	sheetWidth    = 1024
	sheetHeight   = 512
	linebufWidth  = 1024
	frameWidth    = sheetWidth
	frameHeight   = sheetHeight
	screenWidth   = 640
	screenHeight  = 360
)

// VideoDisplay implements ebiten.Game interface.
type VideoDisplay struct {
	vdp *vdp.VDP
	cpu *rj32.CPU
	bus rj32.Bus

	speedTime time.Time

	framebuf [screenWidth * screenHeight * 4]byte

	lastPrint  time.Time
	lastUpdate time.Time

	xvel [numSprites]int
	yvel [numSprites]int

	animX anim.Anim
	animY anim.Anim

	presentation bool
	presentY     int
	incTicks     int
	inc          int
	takeY        int
}

// Update calculates what's needed for the next frame
func (g *VideoDisplay) Update() error {
	if g.cpu != nil {
		g.bus = g.cpu.Run(g.bus, 800*449)

		dur := time.Since(g.speedTime)
		if dur > 20*time.Second {
			g.speedTime = time.Now()
			fmt.Printf("%.4f MHz\n", float64(g.cpu.Cycles)/dur.Seconds()/1000000)
			g.cpu.Cycles = 0
		}
	}

	if g.presentation {
		_, yoff := ebiten.Wheel()
		inc := int(math.Round(yoff * 16))
		if inc != 0 {
			g.inc += int(math.Round(yoff * 16))

			g.incTicks += 4
		}

		if g.incTicks > 0 {
			amt := g.inc / g.incTicks
			g.presentY += amt
			g.inc -= amt

			if g.presentY > 0 {
				g.presentY = 0
			}
			g.incTicks--
		}

		y := 0
		for i := 0; i < g.vdp.NumRenderedSprites; i++ {
			g.vdp.SetSpritePos(i, 0, y+g.presentY)
			y += 16 * 8
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyD) {
			g.loadTest()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyT) {
			g.presentY = g.takeY
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyN) {
			g.takeY = g.presentY
		}

		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if g.vdp.NumRenderedSprites == vdp.NumSprites {
			g.vdp.NumRenderedSprites = 5
		} else {
			g.vdp.NumRenderedSprites = vdp.NumSprites
		}
		fmt.Println("sprites:", g.vdp.NumRenderedSprites)
	}

	dur := time.Since(g.lastUpdate)
	if dur > 1*time.Second {
		dur = 16 * time.Millisecond
	}
	g.animX.Advance(dur)
	g.animY.Advance(dur)

	g.vdp.X[4] = int16(g.animX.Value())
	g.vdp.Y[4] = int16(g.animY.Value())

	for i := 5; i < vdp.NumSprites; i++ {
		g.vdp.X[i] += int16(g.xvel[i])
		if g.vdp.X[i] > (vdp.ScreenWidth + 64) {
			g.vdp.X[i] = -64
		}
		if g.vdp.X[i] < -64 {
			g.vdp.X[i] = (vdp.ScreenWidth + 64)
		}

		g.vdp.Y[i] += int16(g.yvel[i])
		if g.vdp.Y[i] > (vdp.ScreenHeight + 64) {
			g.vdp.Y[i] = -64
		}
		if g.vdp.Y[i] < -64 {
			g.vdp.Y[i] = (vdp.ScreenHeight + 64)
		}
	}

	if time.Since(g.lastPrint) > 5*time.Second {
		g.lastPrint = time.Now()
		fmt.Println("FPS:", ebiten.CurrentFPS(), "TPS:", ebiten.CurrentTPS())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.loadVid()
	}

	return nil
}

// Draw sprites on the screen
func (g *VideoDisplay) Draw(screen *ebiten.Image) {
	g.vdp.DrawFrame(g.framebuf[:])

	// graphics library specific
	screen.ReplacePixels(g.framebuf[:])
}

// Layout sets the scaling -- 2:1 in this case
func (g *VideoDisplay) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth / 2, outsideHeight / 2
}

func (g *VideoDisplay) loadTest() {
	g.presentation = false
	v := g.vdp
	v.ClearSprites(0, 256)
	v.ResetMemMap()
	v.NumRenderedSprites = 0

	err := v.LoadPalette("testpal.hex")
	if err != nil {
		panic(err)
	}

	err = v.LoadSheetSets("testmap_%d.hex", "testtiles_%d.hex")
	if err != nil {
		panic(err)
	}

	v.Dims[0] = v.Dims[0].SetWidth(80).SetHeight(8)
	v.SPos[0] = v.SPos[0].SetSheetY(8)
	v.SetSpriteSheetSet(0, 0, 0)
	v.NumRenderedSprites++

	v.Y[1] = (8 * 8)
	v.Dims[1] = v.Dims[1].SetWidth(80).SetHeight(16)
	v.SetSpriteSheetSet(1, 1, 1)
	v.NumRenderedSprites++

	v.Y[2] = ((8 + 16) * 8)
	v.Dims[2] = v.Dims[2].SetWidth(80).SetHeight(16)
	v.SetSpriteSheetSet(2, 2, 2)
	v.NumRenderedSprites++

	v.Y[3] = ((8 + 16 + 16) * 8)
	v.Dims[3] = v.Dims[3].SetWidth(80).SetHeight(8)
	v.SetSpriteSheetSet(3, 3, 3)
	v.NumRenderedSprites++

	for i := 4; i < vdp.NumSprites; i++ {
		v.SetSpriteSheetSet(i, 0, 0)

		if rand.Int31n(50) == 0 {
			// big
			v.SetSpriteDims(i, 8, 8)
			v.SetSpriteSheetPos(i, 0, 0)
		} else if rand.Int31n(5) < 1 {
			// medium
			v.SetSpriteDims(i, 4, 4)
			v.SetSpriteSheetPos(i, 8, 0)
		} else if rand.Int31n(5) < 4 {
			// small
			v.SetSpriteDims(i, 2, 2)
			v.SetSpriteSheetPos(i, 12, 0)
		} else {
			// tiny
			v.SetSpriteDims(i, 1, 1)
			v.SetSpriteSheetPos(i, 14, 0)
		}

		v.Addr[i] = v.Addr[i].SetTransparent(true)
		v.SetSpritePos(i,
			int(rand.Int31n(vdp.ScreenWidth)),
			int(rand.Int31n(vdp.ScreenHeight)))

		for g.xvel[i] == 0 {
			g.xvel[i] = int(rand.Int31n(10) - 5)
		}
		for g.yvel[i] == 0 {
			g.yvel[i] = int(rand.Int31n(10) - 5)
		}
	}

	v.NumRenderedSprites++
	g.animX.Next(100, 0, anim.EaseInOut)
	g.animX.Next(640-100, 2*time.Second, anim.EaseInOut)
	g.animX.Next(640-100, 2*time.Second, anim.EaseInOut)
	g.animX.Next(100, 2*time.Second, anim.EaseInOut)
	g.animX.Next(100, 2*time.Second, anim.EaseInOut)
	g.animX.Loop = true

	g.animY.Next(50, 0, anim.EaseInOut)
	g.animY.Next(50, 2*time.Second, anim.EaseInOut)
	g.animY.Next(360-50, 2*time.Second, anim.EaseInOut)
	g.animY.Next(360-50, 2*time.Second, anim.EaseInOut)
	g.animY.Next(50, 2*time.Second, anim.EaseInOut)
	g.animY.Loop = true
}

func (g *VideoDisplay) loadVid() {
	g.presentation = true
	v := g.vdp
	v.ClearSprites(0, 256)
	v.ResetMemMap()

	err := v.LoadPalette("vidpal.hex")
	if err != nil {
		panic(err)
	}

	err = v.LoadSheetSets("vidmap_%d.hex", "vidtiles_%d.hex")
	if err != nil {
		panic(err)
	}

	v.NumRenderedSprites = 22

	for i := 0; i < 11; i++ {
		v.SetSpriteDims((i*2)+0, 80, 16)
		v.SetSpriteSheetSet((i*2)+0, i, i)

		v.SetSpriteDims((i*2)+1, 80, 16)
		v.SetSpriteSheetSet((i*2)+1, i, i)
		v.SetSpriteSheetPos((i*2)+1, 0, 16)
	}

	g.presentation = true
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var novdp = flag.Bool("novdp", false, "disable vdp")
var run = flag.String("run", "", "run program from hex")
var trace = flag.Bool("trace", false, "trace cpu instructions")
var maxcycles = flag.Int("maxcycles", 0, "max cycles to run (0: infinity)")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var cpu *rj32.CPU
	if *run != "" {
		cpu = &rj32.CPU{}
		cpu.Trace = *trace
		var buf []byte
		var err error
		if *run == "-" {
			buf, err = io.ReadAll(os.Stdin)
		} else {
			buf, err = os.ReadFile(*run)
		}
		if err != nil {
			panic(err)
		}
		data.Load(16, buf, func(a int, val uint64) {
			cpu.Prog[a] = rj32.DecodeInst(uint16(val))
		})
	}

	if *novdp {
		if cpu == nil {
			return
		}
		var bus rj32.Bus
		before := time.Now()
		for !cpu.Halt && !cpu.Error {
			max := 100000
			if *maxcycles != 0 {
				max = *maxcycles - int(cpu.Cycles)
			}
			bus = cpu.Run(bus, max)
			cpu.Cycles++

			if *maxcycles > 0 && cpu.Cycles > uint64(*maxcycles) {
				fmt.Println("Failed to terminate in time")
				os.Exit(1)
			}
			dur := time.Since(before)
			if dur > 20*time.Second {
				before = time.Now()
				fmt.Printf("%.4f MHz\n", float64(cpu.Cycles)/dur.Seconds()/1000000)
				cpu.Cycles = 0
			}
		}
		if cpu.Error {
			fmt.Println("error!")
			os.Exit(1)
		}
		if cpu.Halt {
			fmt.Println("success!")
			os.Exit(0)
		}
	}

	vdp := vdp.NewVDP()

	display := &VideoDisplay{
		vdp: vdp,
		cpu: cpu,
	}

	ebiten.Wheel()

	// display.loadTest()
	display.loadVid()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("rj32 emu")

	// ebiten.SetVsyncEnabled(false)

	if err := ebiten.RunGame(display); err != nil {
		log.Fatal(err)
	}
}
