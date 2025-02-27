package main

import (
	"image/color"
	"math/rand"

	_ "embed"

	"github.com/solarlune/tetra3d"
	"github.com/solarlune/tetra3d/colors"
	"github.com/solarlune/tetra3d/examples"

	"github.com/hajimehoshi/ebiten/v2"
)

// Shared cube mesh
var cubeMesh = tetra3d.NewCubeMesh()

type CubeElement struct {
	Model     *tetra3d.Model
	Root      *tetra3d.Node
	Target    *tetra3d.GridPoint
	Navigator *tetra3d.Navigator
}

func newCubeElement(root tetra3d.INode) *CubeElement {

	grid := root.Get("Network").(*tetra3d.Grid)

	element := &CubeElement{
		Model:     tetra3d.NewModel(cubeMesh, "Cube Element"),
		Root:      root.(*tetra3d.Node),
		Navigator: tetra3d.NewNavigator(nil),
	}

	element.Navigator.FinishMode = tetra3d.FinishModeStop
	element.Model.Color.Set(0.8+rand.Float32()*0.2, rand.Float32()*0.5, rand.Float32()*0.5, 1)
	element.Model.SetWorldScale(0.1, 0.1, 0.1)
	element.Model.SetWorldPositionVec(grid.RandomPoint().WorldPosition())

	element.ChooseNewTarget()
	root.AddChildren(element.Model)
	return element
}

func (cube *CubeElement) Update() {

	if cube.Navigator.HasPath() {
		cube.Navigator.AdvanceDistance(0.05)
		cube.Model.SetWorldPositionVec(cube.Navigator.WorldPosition())
		if cube.Navigator.Finished() {
			cube.Navigator.SetPath(nil)
		}
	} else {
		cube.ChooseNewTarget()
	}

}

func (cube *CubeElement) ChooseNewTarget() {
	grid := cube.Root.Get("Network").(*tetra3d.Grid)
	closest := grid.NearestGridPoint(cube.Model.WorldPosition())
	cube.Navigator.SetPath(closest.PathTo(grid.RandomPoint()))
}

type Game struct {
	Scene  *tetra3d.Scene
	Cubes  []*CubeElement
	Camera examples.BasicFreeCam
	System examples.BasicSystemHandler
}

//go:embed grids.gltf
var grids []byte

func NewGame() *Game {
	game := &Game{}

	game.Init()

	return game
}

func (g *Game) Init() {

	data, err := tetra3d.LoadGLTFData(grids, nil)
	if err != nil {
		panic(err)
	}

	g.Scene = data.Scenes[0]

	g.Scene.World.LightingOn = false

	for i := 0; i < 40; i++ {
		newCube := newCubeElement(g.Scene.Root)
		g.Cubes = append(g.Cubes, newCube)
	}

	g.Camera = examples.NewBasicFreeCam(g.Scene)
	g.System = examples.NewBasicSystemHandler(g)

}

func (g *Game) Update() error {

	for _, cube := range g.Cubes {
		cube.Update()
	}

	g.Camera.Update()

	return g.System.Update()

}

func (g *Game) Draw(screen *ebiten.Image) {

	// Clear, but with a color
	screen.Fill(color.RGBA{60, 70, 80, 255})

	// Clear the Camera
	g.Camera.Clear()

	// Render the logo first
	g.Camera.RenderScene(g.Scene)

	// We rescale the depth or color textures here just in case we render at a different resolution than the window's; this isn't necessary,
	// we could just draw the images straight.
	screen.DrawImage(g.Camera.ColorTexture(), nil)

	g.System.Draw(screen, g.Camera.Camera)

	if g.System.DrawDebugText {
		txt := `This example shows how Grids work.
Grids are networks of vertices,
connected by their edges. Navigators can
navigate from point to point on Grids.

Currently, navigation is calculated using
number of hops, rather than overall
distance to navigate.`
		g.Camera.DebugDrawText(screen, txt, 0, 200, 1, colors.White())
	}

}

func (g *Game) Layout(w, h int) (int, int) {
	return g.Camera.Size()
}

func main() {
	ebiten.SetWindowTitle("Tetra3d - Webs Test")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
