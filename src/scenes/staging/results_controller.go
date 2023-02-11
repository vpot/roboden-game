package staging

import (
	"fmt"
	"time"

	"github.com/ebitenui/ebitenui/widget"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/roboden-game/assets"
	"github.com/quasilyte/roboden-game/controls"
	"github.com/quasilyte/roboden-game/gameui/eui"
	"github.com/quasilyte/roboden-game/session"
)

type resultsController struct {
	state *session.State

	scene          *ge.Scene
	backController ge.SceneController

	results battleResults
}

type battleResults struct {
	Victory         bool
	TimePlayed      time.Duration
	SurvivingDrones int

	ResourcesGathered float64
	DronesProduced    int
	CreepsDefeated    int
}

func newResultsController(state *session.State, backController ge.SceneController, results battleResults) *resultsController {
	return &resultsController{
		state:          state,
		backController: backController,
		results:        results,
	}
}

func (c *resultsController) Init(scene *ge.Scene) {
	c.scene = scene
	c.initUI()
}

func (c *resultsController) Update(delta float64) {
	if c.state.MainInput.ActionIsJustPressed(controls.ActionBack) {
		c.back()
		return
	}
}

func (c *resultsController) initUI() {
	uiResources := eui.LoadResources(c.scene.Context().Loader)

	root := eui.NewAnchorContainer()
	rowContainer := eui.NewRowLayoutContainer()
	root.AddChild(rowContainer)

	smallFont := c.scene.Context().Loader.LoadFont(assets.FontSmall).Face

	titleString := "Victory!"
	if !c.results.Victory {
		titleString = "Defeat"
	}
	titleLabel := eui.NewLabel(uiResources, titleString, smallFont)
	rowContainer.AddChild(titleLabel)

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}))

	var timeString string
	timePlayed := c.results.TimePlayed
	if timePlayed.Minutes() < 1 {
		timeString = fmt.Sprintf("%ds", int(timePlayed.Seconds()))
	} else if timePlayed.Minutes() < 60 {
		timeString = fmt.Sprintf("%dm %ds", int(timePlayed.Minutes()), int(timePlayed.Seconds()))
	} else {
		timeString = fmt.Sprintf("%dh %dm %ds", int(timePlayed.Hours()), int(timePlayed.Minutes()), int(timePlayed.Seconds()))
	}

	lines := []string{
		fmt.Sprintf("Time Played: %v", timeString),
		fmt.Sprintf("Resources Gathered: %v", int(c.results.ResourcesGathered)),
		fmt.Sprintf("Drone Survivors: %v", c.results.SurvivingDrones),
		fmt.Sprintf("Drones Total: %v", c.results.DronesProduced),
		fmt.Sprintf("Creeps Defeated: %v", c.results.CreepsDefeated),
	}

	for _, l := range lines {
		label := eui.NewLabel(uiResources, l, smallFont)
		rowContainer.AddChild(label)
	}

	rowContainer.AddChild(eui.NewSeparator(widget.RowLayoutData{Stretch: true}))

	rowContainer.AddChild(eui.NewButton(uiResources, c.scene, "To The Menu", func() {
		c.back()
	}))

	uiObject := eui.NewSceneObject(root)
	c.scene.AddGraphics(uiObject)
	c.scene.AddObject(uiObject)
}

func (c *resultsController) back() {
	c.scene.Context().ChangeScene(c.backController)
}
