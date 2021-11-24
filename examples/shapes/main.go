package main

import (
	"errors"
	"fmt"
	"image/color"
	"image/png"
	"math"
	"os"
	"runtime/pprof"
	"strings"
	"time"

	_ "embed"
	_ "image/png"

	"github.com/kvartborg/vector"
	"github.com/solarlune/tetra3d"
	"golang.org/x/image/font/basicfont"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

//go:embed shapes.dae
var shapesDAE []byte

type Game struct {
	Width, Height int
	Scene         *tetra3d.Scene

	Camera       *tetra3d.Camera
	CameraTilt   float64
	CameraRotate float64

	Time              float64
	DrawDebugText     bool
	DrawDebugDepth    bool
	PrevMousePosition vector.Vector
}

func NewGame() *Game {
	game := &Game{
		Width:             398,
		Height:            224,
		PrevMousePosition: vector.Vector{},
		DrawDebugText:     true,
	}

	game.Init()

	return game
}

func (g *Game) Init() {
	// Load the DAE file and turn it into a scene. Note that we could also pass options to change how the file
	// is loaded (specifically, which way is up), but we don't have to do that because it will do this by default if
	// nil is passed as the second argument.
	dae, err := tetra3d.LoadDAEData(shapesDAE, nil)
	if err != nil {
		panic(err)
	}

	g.Scene = dae.Scenes[0]

	// Unfortunately, we have to manually load images and apply them to Meshes; the easiest way to do this for now
	// seems to simply be to give Meshes a Material with the name of their image file.
	for _, m := range dae.Meshes {
		if strings.HasSuffix(m.MaterialName, ".png") {
			// Put images on meshes that have material names that end with ".png"
			m.Image, _, _ = ebitenutil.NewImageFromFile(m.MaterialName)
		}
	}

	g.Camera = tetra3d.NewCamera(g.Width, g.Height)
	g.Camera.SetLocalPosition(vector.Vector{0, 0, 5})
	g.Camera.Far = 20
	// g.Camera.RenderDepth = false // You can turn off depth rendering if your computer doesn't do well with shaders or rendering to offscreen buffers,
	// but this will turn off inter-object depth sorting. Instead, Tetra's Camera will render objects in order of distance to camera.

	ebiten.SetCursorMode(ebiten.CursorModeCaptured)

	fmt.Println(g.Scene.Root.TreeToString())

}

func (g *Game) Update() error {
	var err error

	moveSpd := 0.1

	g.Time += 1.0 / 60

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		err = errors.New("quit")
	}

	// We use Camera.Rotation.Forward().Invert() because the camera looks down -Z (so its forward vector is inverted)
	forward := g.Camera.LocalRotation().Forward().Invert()
	right := g.Camera.LocalRotation().Right()

	pos := g.Camera.LocalPosition()

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		pos = pos.Add(forward.Scale(moveSpd))
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		pos = pos.Add(right.Scale(moveSpd))
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		pos = pos.Add(forward.Scale(-moveSpd))
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		pos = pos.Add(right.Scale(-moveSpd))
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		pos[1] += moveSpd
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		pos[1] -= moveSpd
	}

	g.Camera.SetLocalPosition(pos)

	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	// Rotate and tilt the camera according to mouse movements
	mx, my := ebiten.CursorPosition()

	mv := vector.Vector{float64(mx), float64(my)}

	diff := mv.Sub(g.PrevMousePosition)

	g.CameraTilt -= diff[1] * 0.005
	g.CameraRotate -= diff[0] * 0.005

	g.CameraTilt = math.Max(math.Min(g.CameraTilt, math.Pi/2-0.1), -math.Pi/2+0.1)

	tilt := tetra3d.NewMatrix4Rotate(1, 0, 0, g.CameraTilt)
	rotate := tetra3d.NewMatrix4Rotate(0, 1, 0, g.CameraRotate)

	// Order of this is important - tilt * rotate works, rotate * tilt does not, lol
	g.Camera.SetLocalRotation(tilt.Mult(rotate))

	g.PrevMousePosition = mv.Clone()

	// Fog controls
	if ebiten.IsKeyPressed(ebiten.Key1) {
		g.Scene.FogColor.SetRGB(1, 0, 0)
		g.Scene.FogMode = tetra3d.FogAdd
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		g.Scene.FogColor.SetRGB(0, 0, 0)
		g.Scene.FogMode = tetra3d.FogMultiply
	} else if ebiten.IsKeyPressed(ebiten.Key3) {
		g.Scene.FogColor.SetRGB(0, 0, 0)
		g.Scene.FogMode = tetra3d.FogOverwrite
	} else if ebiten.IsKeyPressed(ebiten.Key4) {
		g.Scene.FogMode = tetra3d.FogOff
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		f, err := os.Create("screenshot" + time.Now().Format("2006-01-02 15:04:05") + ".png")
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		png.Encode(f, g.Camera.ColorTexture)
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		g.Init()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.StartProfiling()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.DrawDebugText = !g.DrawDebugText
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		g.Camera.DebugDrawWireframe = !g.Camera.DebugDrawWireframe
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF3) {
		g.Camera.DebugDrawNormals = !g.Camera.DebugDrawNormals
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF5) {
		g.DrawDebugDepth = !g.DrawDebugDepth
	}

	return err
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear, but with a color
	// screen.Fill(color.RGBA{20, 25, 30, 255})
	screen.Fill(color.Black)

	g.Camera.Clear()

	g.Camera.RenderNodes(g.Scene, g.Scene.Root)

	opt := &ebiten.DrawImageOptions{}
	w, h := g.Camera.ColorTexture.Size()

	// We rescale the depth or color textures just in case we render at a different resolution than the window's.
	opt.GeoM.Scale(float64(g.Width)/float64(w), float64(g.Height)/float64(h))
	if g.DrawDebugDepth {
		screen.DrawImage(g.Camera.DepthTexture, opt)
	} else {
		screen.DrawImage(g.Camera.ColorTexture, opt)
	}

	if g.DrawDebugText {
		g.Camera.DrawDebugText(screen, 1)
		txt := "F1 to toggle this text\nWASD: Move, Mouse: Look\n1, 2, 3, 4: Change fog\nF1, F2, F3, F5: Debug views\nF4: Toggle fullscreen\nESC: Quit"
		text.Draw(screen, txt, basicfont.Face7x13, 0, 128, color.RGBA{255, 0, 0, 255})
	}
}

func (g *Game) StartProfiling() {
	outFile, err := os.Create("./cpu.pprof")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Beginning CPU profiling...")
	pprof.StartCPUProfile(outFile)
	go func() {
		time.Sleep(2 * time.Second)
		pprof.StopCPUProfile()
		fmt.Println("CPU profiling finished.")
	}()
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.Width, g.Height
}

func main() {
	ebiten.SetWindowTitle("Tetra3d Test - Shapes")
	ebiten.SetWindowResizable(true)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
