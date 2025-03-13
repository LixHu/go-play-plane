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
	targetScore        int      // 当前关卡目标分数
}

// Update 处理游戏逻辑更新
func (g *Game) Update() error {
	// 在菜单模式下处理模式选择
	if g.gameMode == ModeMenu {
		// 按1选择关卡模式
		if ebiten.IsKeyPressed(ebiten.Key1) {
			g.gameMode = ModeLevels
			g.currentLevel = 1
			g.targetScore = g.currentLevel * 1000 // 每关目标分数为关卡数 * 1000
			g.enemyManager.SetLevel(g.currentLevel)
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

	// 如果游戏已结束，处理重新开始或返回菜单的输入
	if g.isGameOver {
		// 按空格键重新开始当前模式
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			// 重置游戏状态
			g.player = NewPlayer()
			g.enemyManager = NewEnemyManager()
			g.bulletManager = NewBulletManager()
			g.enemyBulletManager = NewEnemyBulletManager()
			g.powerUpManager = NewPowerUpManager()
			g.score = 0
			g.isGameOver = false
			if g.gameMode == ModeLevels {
				g.currentLevel = 1
				g.targetScore = g.currentLevel * 1000
				g.enemyManager.SetLevel(g.currentLevel)
			}
		}
		// 按ESC键返回模式选择界面
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.gameMode = ModeMenu
			g.isGameOver = false
			g.score = 0
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

				// 在关卡模式下检查是否达到目标分数
				if g.gameMode == ModeLevels && g.score >= g.targetScore {
					// 进入下一关
					g.currentLevel++
					g.targetScore = g.currentLevel * 1000
					g.enemyManager.SetLevel(g.currentLevel)
				}
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
	if g.gameMode == ModeMenu {
		// 绘制半透明背景
		ebitenutil.DrawRect(screen, 0, 0, float64(screenWidth), float64(screenHeight), color.RGBA{0, 0, 50, 180})

		// 绘制标题
		titleMsg := "飞机大战"
		titleX := screenWidth/2 - 100
		titleY := screenHeight/3 - 30
		// 标题背景
		ebitenutil.DrawRect(screen, float64(titleX-20), float64(titleY-40), 240, 50, color.RGBA{0, 0, 100, 200})
		text.Draw(screen, titleMsg, chineseFont, titleX, titleY, color.RGBA{255, 255, 255, 255})

		// 绘制模式选择说明
		modeTitle := "- 游戏模式选择 -"
		modeTitleX := screenWidth/2 - 80
		modeTitleY := screenHeight/2 - 50
		// 模式选择背景
		ebitenutil.DrawRect(screen, float64(modeTitleX-20), float64(modeTitleY-25), 200, 35, color.RGBA{0, 0, 100, 150})
		text.Draw(screen, modeTitle, chineseFont, modeTitleX, modeTitleY, color.RGBA{200, 200, 255, 255})

		// 绘制模式选项
		mode1 := "[1] 关卡模式"
		mode1Desc := "逐级挑战，难度递增"
		mode2 := "[2] 无尽模式"
		mode2Desc := "无限挑战，直到失败"

		mode1X := screenWidth/2 - 150
		mode1Y := screenHeight/2 + 20
		mode2X := screenWidth/2 - 150
		mode2Y := screenHeight/2 + 70

		// 模式选项背景
		ebitenutil.DrawRect(screen, float64(mode1X-10), float64(mode1Y-25), 320, 40, color.RGBA{0, 0, 100, 100})
		ebitenutil.DrawRect(screen, float64(mode2X-10), float64(mode2Y-25), 320, 40, color.RGBA{0, 0, 100, 100})

		// 绘制模式选项文字
		text.Draw(screen, mode1, chineseFont, mode1X, mode1Y, color.RGBA{255, 255, 0, 255})
		text.Draw(screen, mode1Desc, chineseFont, mode1X+140, mode1Y, color.White)
		text.Draw(screen, mode2, chineseFont, mode2X, mode2Y, color.RGBA{255, 255, 0, 255})
		text.Draw(screen, mode2Desc, chineseFont, mode2X+140, mode2Y, color.White)

		// 绘制操作提示
		hint := "按对应数字键选择模式"
		hintX := screenWidth/2 - 100
		hintY := screenHeight*3/4 + 30
		// 提示背景
		ebitenutil.DrawRect(screen, float64(hintX-20), float64(hintY-25), 240, 35, color.RGBA{0, 0, 100, 150})
		text.Draw(screen, hint, chineseFont, hintX, hintY, color.RGBA{200, 200, 255, 255})
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
	scoreText := fmt.Sprintf("得分: %d", g.score)
	scoreX := 20
	scoreY := 30
	// 分数背景
	ebitenutil.DrawRect(screen, float64(scoreX-10), float64(scoreY-25), 150, 35, color.RGBA{0, 0, 100, 150})
	text.Draw(screen, scoreText, chineseFont, scoreX, scoreY, color.RGBA{255, 255, 0, 255})

	// 在关卡模式下显示当前关卡和目标分数
	if g.gameMode == ModeLevels {
		levelText := fmt.Sprintf("当前关卡: %d", g.currentLevel)
		levelX := 20
		levelY := 70
		// 关卡背景
		ebitenutil.DrawRect(screen, float64(levelX-10), float64(levelY-25), 150, 35, color.RGBA{0, 0, 100, 150})
		text.Draw(screen, levelText, chineseFont, levelX, levelY, color.RGBA{255, 255, 0, 255})

		targetText := fmt.Sprintf("目标分数: %d", g.targetScore)
		targetX := 20
		targetY := 110
		// 目标分数背景
		ebitenutil.DrawRect(screen, float64(targetX-10), float64(targetY-25), 150, 35, color.RGBA{0, 0, 100, 150})
		text.Draw(screen, targetText, chineseFont, targetX, targetY, color.RGBA{255, 255, 0, 255})
	}

	// 如果游戏结束，显示游戏结束信息
	if g.isGameOver {
		// 绘制半透明背景
		ebitenutil.DrawRect(screen, 0, 0, float64(screenWidth), float64(screenHeight), color.RGBA{0, 0, 0, 180})

		// 绘制游戏结束标题
		gameOverMsg := "游戏结束！"
		gameOverX := screenWidth/2 - 80
		gameOverY := screenHeight/2 - 50
		// 标题背景
		ebitenutil.DrawRect(screen, float64(gameOverX-20), float64(gameOverY-35), 200, 45, color.RGBA{0, 0, 100, 200})
		text.Draw(screen, gameOverMsg, chineseFont, gameOverX, gameOverY, color.RGBA{255, 50, 50, 255})

		// 绘制最终得分
		scoreMsg := fmt.Sprintf("最终得分：%d", g.score)
		scoreX := screenWidth/2 - 80
		scoreY := screenHeight / 2
		// 得分背景
		ebitenutil.DrawRect(screen, float64(scoreX-20), float64(scoreY-25), 200, 35, color.RGBA{0, 0, 100, 150})
		text.Draw(screen, scoreMsg, chineseFont, scoreX, scoreY, color.RGBA{255, 255, 0, 255})

		// 绘制操作提示
		restartMsg := "按空格键重新开始当前模式"
		restartX := screenWidth/2 - 150
		restartY := screenHeight/2 + 50
		menuMsg := "按ESC键返回模式选择"
		menuX := screenWidth/2 - 100
		menuY := screenHeight/2 + 90

		// 提示背景
		ebitenutil.DrawRect(screen, float64(restartX-10), float64(restartY-25), 320, 35, color.RGBA{0, 0, 100, 150})
		ebitenutil.DrawRect(screen, float64(menuX-10), float64(menuY-25), 220, 35, color.RGBA{0, 0, 100, 150})

		text.Draw(screen, restartMsg, chineseFont, restartX, restartY, color.RGBA{200, 200, 255, 255})
		text.Draw(screen, menuMsg, chineseFont, menuX, menuY, color.RGBA{200, 200, 255, 255})
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
