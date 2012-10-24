package rog

type SceneStack struct {
    top *SceneElement
    size int
}

func NewSceneStack() *SceneStack {
    new_stack := &SceneStack {
        top: nil,
        size: 0,
    }

    return new_stack
}

func (self *SceneStack) Len() int {
    return self.size
}

func (self *SceneStack) Push(scene Scene) {
    self.top = &SceneElement{scene, self.top}
    self.size++ 
}

func (self *SceneStack) Pop() Scene {
    if self.size > 0 {
        scene:= self.top.scene 
        self.top = self.top.next
        self.size--
        return scene
    }

    return nil
}

func (self *SceneStack) Top() Scene {
    if self.size > 0 {
        return self.top.scene
    }
    return nil
}

type SceneElement struct {
    scene Scene
    next *SceneElement
}

type Scene interface {
    HandleKeys()
    Update()
    Render()
}