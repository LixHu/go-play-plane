package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth  = 640
	screenHeight = 480
	gameTitle    = "打飞机游戏"
)

// GameMode 游戏模式
type GameMode int

const (
	ModeMenu    GameMode = iota // 菜单模式
	ModeLevels                  // 关卡模式
	ModeEndless                 // 无尽模式
)

var (
	gameFont    font.Face
	chineseFont font.Face
)

func init() {
	// 加载英文字体
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	gameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 加载中文字体
	chineseFontData, err := os.ReadFile("resources/fonts/SourceHanSansCN-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	chineseTT, err := opentype.Parse(chineseFontData)
	if err != nil {
		log.Fatal(err)
	}

	chineseFont, err = opentype.NewFace(chineseTT, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// Game 结构体用于保存游戏状态
type Game struct {
	player             *Player
	enemyManager       *EnemyManager
	bulletManager      *BulletManager
	enemyBulletManager *EnemyBulletManager
	powerUpManager     *PowerUpManager
	score              int
	isGameOver         bool
	gameMode           GameMode // 当前游戏模式
	currentLevel       int      // 当前关卡（仅用于关卡模式）
	difficulty         float64  // 游戏难度系数
}

// Update 处理游戏逻辑更新
func (g *Game) Update() error {
	// 在菜单模式下处理模式选择
	if g.gameMode == ModeMenu {
		// 按1选择关卡模式
		if ebiten.IsKeyPressed(ebiten.Key1) {
			g.gameMode = ModeLevels
			g.currentLevel = 1
			g.difficulty = 1.0
			return nil
		}
		// 按2选择无尽模式
		if ebiten.IsKeyPressed(ebiten.Key2) {
			g.gameMode = ModeEndless
			g.difficulty = 1.0
			return nil
		}
		return nil
	}

	// 如果游戏已结束，只处理重新开始的输入
	if g.isGameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			// 重置游戏状态
			g.player = NewPlayer()
			g.enemyManager = NewEnemyManager()
			g.bulletManager = NewBulletManager()
			g.enemyBulletManager = NewEnemyBulletManager()
			g.powerUpManager = NewPowerUpManager()
			g.score = 0
			g.isGameOver = false
		}
		return nil
	}

	// 更新玩家状态
	g.player.Update()

	// 更新敌机状态
	g.enemyManager.Update()

	// 更新子弹状态
	g.bulletManager.Update(g.player)

	// 更新道具状态
	g.powerUpManager.Update()

	// 检测子弹与敌机的碰撞
	for _, bullet := range g.bulletManager.bullets {
		for _, enemy := range g.enemyManager.enemies {
			if bullet.CheckCollision(enemy) {
				bullet.active = false
				enemy.active = false
				g.score += 100
				// 在敌机被击毁的位置生成道具
				g.powerUpManager.SpawnPowerUp(enemy.x, enemy.y)
			}
		}
	}

	// 检测玩家与道具的碰撞
	for _, powerUp := range g.powerUpManager.powerUps {
		if powerUp.active && g.checkPlayerPowerUpCollision(powerUp) {
			powerUp.active = false
			// 根据道具类型给予玩家相应的能力
			switch powerUp.pType {
			case MultiShot:
				g.player.EnableMultiShot()
			case ScreenShot:
				g.player.EnableScreenShot()
			}
		}
	}

	// 检测玩家与敌机的碰撞
	for _, enemy := range g.enemyManager.enemies {
		if enemy.active && g.checkPlayerCollision(enemy) {
			g.isGameOver = true
			break
		}
	}

	// 更新敌机子弹状态
	g.enemyBulletManager.Update(g.enemyManager.enemies)

	// 检测敌机子弹与玩家的碰撞
	for _, bullet := range g.enemyBulletManager.bullets {
		if bullet.CheckCollision(g.player) {
			g.isGameOver = true
			break
		}
	}

	return nil
}

// checkPlayerCollision 检测玩家与敌机的碰撞
func (g *Game) checkPlayerCollision(enemy *Enemy) bool {
	return g.player.x < enemy.x+float64(enemy.width) &&
		g.player.x+float64(g.player.width) > enemy.x &&
		g.player.y < enemy.y+float64(enemy.height) &&
		g.player.y+float64(g.player.height) > enemy.y
}

// checkPlayerPowerUpCollision 检测玩家与道具的碰撞
func (g *Game) checkPlayerPowerUpCollision(powerUp *PowerUp) bool {
	return g.player.x < powerUp.x+float64(powerUp.width) &&
		g.player.x+float64(g.player.width) > powerUp.x &&
		g.player.y < powerUp.y+float64(powerUp.height) &&
		g.player.y+float64(g.player.height) > powerUp.y
}

// Draw 处理游戏画面渲染
func (g *Game) Draw(screen *ebiten.Image) {
	// 在菜单模式下显示模式选择界面
	if g.gameMode == ModeMenu {
		// 绘制标题
		titleMsg := "飞机大战"
		titleX := screenWidth/2 - 100
		titleY := screenHeight/3 - 30
		// ebitenutil.DrawRect(screen, float64(titleX-10), float64(titleY-25), 220, 40, color.RGBA{0, 0, 100, 100})
		text.Draw(screen, titleMsg, chineseFont, titleX, titleY, color.White)

		// 绘制模式选择说明
		modeTitle := "- 游戏模式选择 -"
		modeTitleX := screenWidth/2 - 80
		modeTitleY := screenHeight/2 - 50
		text.Draw(screen, modeTitle, chineseFont, modeTitleX, modeTitleY, color.White)

		// 绘制模式选项
		mode1 := "[1] 关卡模式"
		mode1Desc := "逐级挑战，难度递增"
		mode2 := "[2] 无尽模式"
		mode2Desc := "无限挑战，直到失败"

		mode1X := screenWidth/2 - 150
		mode1Y := screenHeight / 2
		mode2X := screenWidth/2 - 150
		mode2Y := screenHeight/2 + 30

		// 使用更大的字体大小
		text.Draw(screen, mode1, chineseFont, mode1X, mode1Y, color.White)
		text.Draw(screen, mode1Desc, chineseFont, mode1X+140, mode1Y, color.White)
		text.Draw(screen, mode2, chineseFont, mode2X, mode2Y, color.White)
		text.Draw(screen, mode2Desc, chineseFont, mode2X+140, mode2Y, color.White)

		// 绘制操作提示，使用更明显的颜色和位置
		hint := "按对应数字键选择模式"
		hintX := screenWidth/2 - 100
		hintY := screenHeight * 2 / 3
		// 绘制提示文字的背景
		ebitenutil.DrawRect(screen, float64(hintX-10), float64(hintY-15), 220, 30, color.RGBA{0, 0, 100, 100})
		text.Draw(screen, hint, chineseFont, hintX, hintY, color.White)
		return
	}

	// 绘制玩家
	g.player.Draw(screen)

	// 绘制敌机
	g.enemyManager.Draw(screen)

	// 绘制子弹
	g.bulletManager.Draw(screen)

	// 绘制敌机子弹
	g.enemyBulletManager.Draw(screen)

	// 绘制道具
	g.powerUpManager.Draw(screen)

	// 绘制分数
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))

	// 如果游戏结束，显示游戏结束信息
	if g.isGameOver {
		msg := fmt.Sprintf("游戏结束！最终得分：%d\n按空格键重新开始", g.score)
		ebitenutil.DebugPrintAt(screen, msg, screenWidth/2-150, screenHeight/2)
	}
}

// Layout 返回游戏窗口的大小
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(gameTitle)

	game := &Game{
		player:             NewPlayer(),
		enemyManager:       NewEnemyManager(),
		bulletManager:      NewBulletManager(),
		enemyBulletManager: NewEnemyBulletManager(),
		powerUpManager:     NewPowerUpManager(),
		score:              0,
		isGameOver:         false,
		gameMode:           ModeMenu,
		currentLevel:       1,
		difficulty:         1.0,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
