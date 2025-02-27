package main

import (
	"bytes"

	_ "embed"

	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"
	"github.com/solarlune/tetra3d/examples"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

//go:embed heart.png
var heartImg []byte

//go:embed scene.gltf
var scene []byte

type Game struct {
	Scene              *tetra3d.Scene
	Camera             *tetra3d.Camera
	WireframeDrawHeart bool
	HeartSprite        *ebiten.Image

	System examples.BasicSystemHandler
}

func NewGame() *Game {
	game := &Game{}

	game.Init()

	return game
}

// In this example, we will simply create a cube and place it in the scene.

func (g *Game) Init() {

	scene, err := tetra3d.LoadGLTFData(scene, nil)
	if err != nil {
		panic(err)
	}
	g.Scene = scene.FindScene("Scene").Clone()

	g.Camera = g.Scene.Root.Get("Camera").(*tetra3d.Camera)

	reader := bytes.NewReader(heartImg)
	newImg, _, err := ebitenutil.NewImageFromReader(reader)
	if err != nil {
		panic(err)
	}

	g.System = examples.NewBasicSystemHandler(g)
	g.System.UsingBasicFreeCam = false
	g.HeartSprite = newImg

}

func (g *Game) Update() error {

	moveSpd := 0.05
	dx := 0.0
	dz := 0.0

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dz -= moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dz += moveSpd
	}

	g.Scene.Root.Get("Heart").Move(dx, 0, dz)

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.WireframeDrawHeart = !g.WireframeDrawHeart
	}

	return g.System.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Clear the screen with a color.
	screen.Fill(g.Scene.World.ClearColor.ToRGBA64())

	// Clear the Camera.
	g.Camera.Clear()

	// Render the scene.
	g.Camera.RenderScene(g.Scene)

	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(-float64(g.HeartSprite.Bounds().Dx())/2, -float64(g.HeartSprite.Bounds().Dy())/2)
	g.Camera.DrawImageIn3D(
		g.Camera.ColorTexture(),
		tetra3d.SpriteRender3d{
			Image:         g.HeartSprite,
			Options:       opt,
			WorldPosition: g.Scene.Root.Get("Heart").WorldPosition(),
		},
	)

	// Draw depth texture if the debug option is enabled; draw color texture otherwise.
	screen.DrawImage(g.Camera.ColorTexture(), nil)

	if g.System.DrawDebugText {
		g.Camera.DebugDrawText(screen, "This is an example showing\nhow you can render a sprite in 2D, while\nmaintaining its ability to render over or under\nother 3D objects by simply moving\nit through 3D space.\n\nA: Toggle wireframe view of heart position", 0, 200, 1, colors.LightGray())
	}

	if g.WireframeDrawHeart {
		g.Camera.DrawDebugWireframe(screen, g.Scene.Root.Get("Heart"), colors.White())
	}

	g.System.Draw(screen, g.Camera)

}

func (g *Game) Layout(w, h int) (int, int) {
	// Here, we simply return the camera's width and height.
	return g.Camera.Size()
}

func main() {

	ebiten.SetWindowTitle("Tetra3d - 3D Sprite Test")

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
