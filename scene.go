package tetra3d

// Scene represents a world of sorts, and can contain a variety of Meshes and Nodes, which organize the scene into a
// graph of parents and children. Models (visual instances of Meshes), Cameras, and "empty" NodeBases all are kinds of Nodes.
type Scene struct {
	Name    string   // The name of the Scene. Set automatically to the scene name in your 3D modeler if the DAE file exports it.
	library *Library // The library from which this Scene was created. If the Scene was instantiated through code, this will be nil.
	// Root indicates the root node for the scene hierarchy. For visual Models to be displayed, they must be added to the
	// scene graph by simply adding them into the tree via parenting anywhere under the Root. For them to be removed from rendering,
	// they simply need to be removed from the tree.
	// See this page for more information on how a scene graph works: https://webglfundamentals.org/webgl/lessons/webgl-scene-graph.html
	Root  INode
	World *World
	props *Properties

	updateAutobatch     bool
	autobatchDynamicMap map[*Material]*Model
	autobatchStaticMap  map[*Material]*Model
}

// NewScene creates a new Scene by the name given.
func NewScene(name string) *Scene {

	scene := &Scene{
		Name:                name,
		Root:                NewNode("Root"),
		World:               NewWorld("World"),
		props:               NewProperties(),
		autobatchDynamicMap: map[*Material]*Model{},
		autobatchStaticMap:  map[*Material]*Model{},
	}

	scene.Root.(*Node).scene = scene

	return scene
}

// Clone clones the Scene, returning a copy. Models and Meshes are shared between them.
func (scene *Scene) Clone() *Scene {

	newScene := NewScene(scene.Name)
	newScene.library = scene.library

	// newScene.Models = append(newScene.Models, scene.Models...)
	newScene.Root = scene.Root.Clone()
	newScene.Root.(*Node).scene = newScene

	newScene.World = scene.World // Here, we simply reference the same world; we don't clone it, since a single world can be shared across multiple Scenes
	newScene.props = scene.props.Clone()

	newScene.updateAutobatch = true

	// Update sectors after cloning the scene
	models := newScene.Root.SearchTree().bySectors().Models()

	for _, n := range models {
		n.sector.Neighbors.Clear()
	}

	for _, n := range models {
		n.sector.UpdateNeighbors(models...)
	}

	return newScene

}

// Library returns the Library from which this Scene was loaded. If it was created through code and not associated with a Library, this function will return nil.
func (scene *Scene) Library() *Library {
	return scene.library
}

func (scene *Scene) Properties() *Properties {
	return scene.props
}

var autobatchBlankMat = NewMaterial("autobatch null material")

func (scene *Scene) HandleAutobatch() {

	if scene.updateAutobatch {

		for _, node := range scene.Root.SearchTree().INodes() {

			if model, ok := node.(*Model); ok {

				if !model.autoBatched {

					mat := autobatchBlankMat

					if mats := model.Mesh.Materials(); len(mats) > 0 {
						mat = mats[0]
					}

					if model.AutoBatchMode == AutoBatchDynamic {

						if _, exists := scene.autobatchDynamicMap[mat]; !exists {
							mesh := NewMesh("auto dynamic batch")
							mesh.AddMeshPart(mat)
							m := NewModel(mesh, "auto dynamic batch")
							m.FrustumCulling = false
							m.dynamicBatcher = true
							scene.autobatchDynamicMap[mat] = m
							scene.Root.AddChildren(m)
						}
						scene.autobatchDynamicMap[mat].DynamicBatchAdd(scene.autobatchDynamicMap[mat].Mesh.MeshParts[0], model)

					} else if model.AutoBatchMode == AutoBatchStatic {

						if _, exists := scene.autobatchStaticMap[mat]; !exists {
							m := NewModel(NewMesh("auto static merge"), "auto static merge")
							scene.autobatchStaticMap[mat] = m
							scene.Root.AddChildren(scene.autobatchStaticMap[mat])
						}
						scene.autobatchStaticMap[mat].StaticMerge(model)

					}

					model.autoBatched = true

				}

			}

		}

		for _, dyn := range scene.autobatchDynamicMap {

			for _, models := range dyn.DynamicBatchModels {

				modelList := append(make([]*Model, 0, len(models)), models...)

				for _, model := range modelList {

					if model.Root() == nil {
						dyn.DynamicBatchRemove(model)
					}

				}

			}

		}

		scene.updateAutobatch = false

	}

}
