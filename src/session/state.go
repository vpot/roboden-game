package session

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gdata"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/langs"
	"github.com/quasilyte/ge/xslices"
	"golang.org/x/image/font"

	"github.com/quasilyte/roboden-game/assets"
	"github.com/quasilyte/roboden-game/gamedata"
	"github.com/quasilyte/roboden-game/gameinput"
	"github.com/quasilyte/roboden-game/gameui/eui"
	"github.com/quasilyte/roboden-game/scenes"
	"github.com/quasilyte/roboden-game/serverapi"
	"github.com/quasilyte/roboden-game/steamsdk"
	"github.com/quasilyte/roboden-game/userdevice"
)

type State struct {
	CPUProfile       string
	CPUProfileWriter io.WriteCloser
	MemProfile       string
	MemProfileWriter io.WriteCloser

	ServerProtocol string
	ServerHost     string
	ServerPath     string

	Device userdevice.Info

	MenuInput *gameinput.Handler

	TouchInput         gameinput.Handler
	CombinedInput      gameinput.Handler
	KeyboardInput      gameinput.Handler
	FirstGamepadInput  gameinput.Handler
	SecondGamepadInput gameinput.Handler

	BoundInputs [2]*gameinput.Handler

	SplashLevelConfig   *gamedata.LevelConfig
	BlitzLevelConfig    *gamedata.LevelConfig
	ClassicLevelConfig  *gamedata.LevelConfig
	ArenaLevelConfig    *gamedata.LevelConfig
	InfArenaLevelConfig *gamedata.LevelConfig
	ReverseLevelConfig  *gamedata.LevelConfig
	TutorialLevelConfig *gamedata.LevelConfig

	Persistent PersistentData

	SceneRegistry scenes.Registry

	Resources Resources

	BackgroundImage *ebiten.Image

	Context *ge.Context

	GameData *gdata.Manager

	SentHighscores bool

	GameCommitHash string

	StdoutLogs []string
}

func (state *State) CheckGameItem(key string) bool {
	if state.GameData == nil {
		return false
	}
	return state.GameData.ItemExists(key)
}

func (state *State) LoadGameItem(key string, dst any) error {
	if state.GameData == nil {
		return nil
	}
	data, err := state.GameData.LoadItem(key)
	if err != nil {
		return err
	}
	if data == nil {
		state.SaveGameItem(key, dst)
		return nil
	}
	return json.Unmarshal(data, dst)
}

func (state *State) SaveGameItem(key string, data any) {
	if state.GameData == nil {
		return
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("can't save game data with key %q: %v", key, err))
	}
	err = state.GameData.SaveItem(key, jsonData)
	if err != nil {
		panic(fmt.Sprintf("can't save game data with key %q: %v", key, err))
	}
}

func (state *State) Logf(format string, args ...any) {
	s := format
	if len(args) != 0 {
		s = fmt.Sprintf(format, args...)
	}

	fmt.Println(s)

	if len(state.StdoutLogs) >= 100 {
		state.StdoutLogs = state.StdoutLogs[:0]
	}
	state.StdoutLogs = append(state.StdoutLogs, s)
}

func (state *State) GetInput(id int) *gameinput.Handler {
	return state.BoundInputs[id]
}

func (state *State) GetGamepadInput(id int) *gameinput.Handler {
	if id == 0 {
		return &state.FirstGamepadInput
	}
	return &state.SecondGamepadInput
}

type PersistentData struct {
	Settings GameSettings

	FirstLaunch bool

	PlayerName string

	NumPendingSubmissions int

	PlayerStats PlayerStats

	CachedBlitzLeaderboard    serverapi.LeaderboardResp
	CachedClassicLeaderboard  serverapi.LeaderboardResp
	CachedArenaLeaderboard    serverapi.LeaderboardResp
	CachedInfArenaLeaderboard serverapi.LeaderboardResp
	CachedReverseLeaderboard  serverapi.LeaderboardResp
}

type PlayerStats struct {
	Achievements []Achievement

	OptionsUnlocked []string
	CoresUnlocked   []string
	TurretsUnlocked []string
	DronesUnlocked  []string
	Tier3DronesSeen []string
	ModesUnlocked   []string

	TutorialCompleted bool

	NumVictories int

	TotalPlayTime time.Duration
	TotalScore    int

	HighestClassicScore           int
	HighestClassicScoreDifficulty int

	HighestBlitzScore           int
	HighestBlitzScoreDifficulty int

	HighestArenaScore           int
	HighestArenaScoreDifficulty int

	HighestInfArenaScore           int
	HighestInfArenaScoreDifficulty int

	HighestReverseScore           int
	HighestReverseScoreDifficulty int
}

type Achievement struct {
	Name  string
	Elite bool
}

type Resources struct {
	UI *eui.Resources

	Font1 font.Face
	Font2 font.Face
	Font3 font.Face
}

type GamepadSettings struct {
	Layout        int
	DeadzoneLevel int
	CursorSpeed   int
}

type GameSettings struct {
	Lang               string
	MusicVolumeLevel   int
	EffectsVolumeLevel int
	ScrollingSpeed     int
	EdgeScrollRange    int
	HintMode           int
	WheelScrollingMode int
	XM                 bool
	ShowFPS            bool
	ShowTimer          bool
	LargeDiodes        bool
	LargerFont         bool
	DebugLogs          bool
	DebugDroneLabels   bool
	Demo               bool
	ScreenButtons      bool
	NoPauseSpeedToggle bool
	IntroDifficulty    int
	IntroSpeed         int
	GamepadSettings    [2]GamepadSettings
	Graphics           GraphicsSettings
	Player1InputMethod int
	Player2InputMethod int
}

type GraphicsSettings struct {
	ShadowsEnabled       bool
	VSyncEnabled         bool
	CameraShakingEnabled bool
	AllShadersEnabled    bool
	ScreenFilter         int
	FullscreenEnabled    bool
	AspectRatio          int
}

const (
	ScreenFilterNone = iota
	ScreenFilterCRT
	ScreenFilterSharpenMinor
	ScreenFilterSharpenMajor
	ScreenFilterHueMinusMinor
	ScreenFilterHueMinusMajor
	ScreenFilterHuePlusMinor
	ScreenFilterHuePlusMajor
)

type SavedReplay struct {
	Date      time.Time
	ResultTag string
	Replay    serverapi.GameReplay
}

func (state *State) GetConfigForMode(m gamedata.Mode) *gamedata.LevelConfig {
	switch m {
	case gamedata.ModeBlitz:
		return state.BlitzLevelConfig
	case gamedata.ModeArena:
		return state.ArenaLevelConfig
	case gamedata.ModeInfArena:
		return state.InfArenaLevelConfig
	case gamedata.ModeClassic:
		return state.ClassicLevelConfig
	case gamedata.ModeReverse:
		return state.ReverseLevelConfig
	case gamedata.ModeTutorial:
		return state.TutorialLevelConfig
	default:
		panic("unexpected game mode")
	}
}

func (state *State) AdjustVolumeLevels() {
	state.Context.Audio.SetGroupVolume(assets.SoundGroupMusic,
		assets.VolumeMultiplier(state.Persistent.Settings.MusicVolumeLevel))
	state.Context.Audio.SetGroupVolume(assets.SoundGroupEffect,
		assets.VolumeMultiplier(state.Persistent.Settings.EffectsVolumeLevel))
}

func (state *State) AdjustTextSize(larger bool) {
	if larger {
		state.Resources.Font1 = assets.Font1_3
	} else {
		state.Resources.Font1 = assets.Font1
	}
	state.Resources.Font2 = assets.Font2
	state.Resources.Font3 = assets.Font3
}

func (state *State) ReloadInputs() {
	state.BoundInputs[0] = state.resolveInputMethod(gameinput.PlayerInputMethod(state.Persistent.Settings.Player1InputMethod))
	state.BoundInputs[1] = state.resolveInputMethod(gameinput.PlayerInputMethod(state.Persistent.Settings.Player2InputMethod))
}

func (state *State) UnlockAchievement(a Achievement) bool {
	stats := &state.Persistent.PlayerStats

	current := xslices.Find(stats.Achievements, func(existing *Achievement) bool {
		return existing.Name == a.Name
	})

	if current != nil {
		if current.Elite {
			return false // Can't be improved
		}
		if !current.Elite && !a.Elite {
			return false // Doesn't improve the rank
		}
		// Upgrade the current achievemnt.
		current.Elite = a.Elite
	} else {
		// It's a new achievement. Add it to the list.
		stats.Achievements = append(stats.Achievements, a)
	}

	if state.Device.Steam.Initialized {
		result := steamsdk.UnlockAchievement(a.Name)
		state.Logf("setting %q steam achievement: %v", a.Name, result)
	}

	return true
}

func (state *State) resolveInputMethod(method gameinput.PlayerInputMethod) *gameinput.Handler {
	switch method {
	case gameinput.InputMethodTouch:
		return &state.TouchInput
	case gameinput.InputMethodCombined:
		return &state.CombinedInput
	case gameinput.InputMethodKeyboard:
		return &state.KeyboardInput
	case gameinput.InputMethodGamepad1:
		return &state.FirstGamepadInput
	case gameinput.InputMethodGamepad2:
		return &state.SecondGamepadInput
	default:
		return &state.CombinedInput
	}
}

func (state *State) ReloadLanguage(ctx *ge.Context) {
	var id resource.RawID
	lang := state.Persistent.Settings.Lang
	switch lang {
	case "en":
		id = assets.RawDictEn
	case "ru":
		id = assets.RawDictRu
	case "ch":
		id = assets.RawDictCh
	default:
		panic(fmt.Sprintf("unsupported lang: %q", lang))
	}
	dict, err := langs.ParseDictionary(lang, 4, ctx.Loader.LoadRaw(id).Data)
	if err != nil {
		panic(err)
	}
	if err := dict.Load("", ctx.Loader.LoadRaw(id+1).Data); err != nil {
		panic(err)
	}
	if err := dict.Load("", ctx.Loader.LoadRaw(id+2).Data); err != nil {
		panic(err)
	}
	if err := dict.Load("", ctx.Loader.LoadRaw(id+3).Data); err != nil {
		panic(err)
	}
	ctx.Dict = dict
}

func (state *State) CacheCommonGlyphs() {
	if !state.Device.IsMobile() {
		return
	}

	alphabet := "0123456789<>.,=+-:()[]&@'\"%!?●◌|\\ ~"
	text.CacheGlyphs(state.Resources.Font1, alphabet)
	text.CacheGlyphs(state.Resources.Font2, alphabet)
	text.CacheGlyphs(state.Resources.Font3, alphabet)
}

func (state *State) CacheGlyphs() {
	if !state.Device.IsMobile() {
		return
	}

	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if state.Persistent.Settings.Lang == "ru" {
		alphabet = "абвгдеёжзийклмнопрстуфхцчшщъыьэюяАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯabxyABXYRL"
	}
	text.CacheGlyphs(state.Resources.Font1, alphabet)
	text.CacheGlyphs(state.Resources.Font2, alphabet)
	text.CacheGlyphs(state.Resources.Font3, alphabet)
}

func (state *State) FindNextReplayIndex() int {
	var minDate time.Time
	minIndex := 0
	for i := 1; i < 10; i++ {
		k := state.ReplayDataKey(i)
		if !state.CheckGameItem(k) {
			return i
		}
		var r SavedReplay
		err := state.LoadGameItem(k, &r)
		if err != nil {
			return i
		}
		if minIndex == 0 || r.Date.Before(minDate) {
			minDate = r.Date
			minIndex = i
		}
	}
	if minIndex != 0 {
		return minIndex
	}
	return 1
}

func (state *State) ReplayDataKey(i int) string {
	return fmt.Sprintf("saved_replay_%d.json", i)
}

func (state *State) SchemaDataKey(m gamedata.Mode, i int) string {
	return fmt.Sprintf("%s_schema_%d.json", m.String(), i)
}
