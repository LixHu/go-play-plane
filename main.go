package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color" // 注册PNG格式支持
	"log"
	"math"
	"math/rand"
	"time"

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
	enemyImage  *ebiten.Image
	playerImage *ebiten.Image
)

func init() {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

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
	chineseFontData, err := chineseFontFS.ReadFile("resources/fonts/SourceHanSansCN-Regular.ttf")
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

	// 加载图片资源
	enemyImage = ebiten.NewImage(32, 32)
	enemyImage.Fill(color.RGBA{255, 0, 0, 255}) // 临时使用红色方块代替敌机

	playerImage = ebiten.NewImage(32, 32)
	playerImage.Fill(color.RGBA{0, 255, 0, 255}) // 临时使用绿色方块代替玩家

	// 加载敌机图片
	enemyImgData, err := imagesFS.ReadFile("resources/images/enemy.png")
	if err != nil {
		log.Printf("无法加载敌机图片，使用默认方块: %v", err)
	} else {
		// 从嵌入的资源创建图片
		enemyImg, _, err := image.Decode(bytes.NewReader(enemyImgData))
		if err != nil {
			log.Printf("无法解码敌机图片，使用默认方块: %v", err)
		} else {
			// 创建指定大小的图片缓冲区
			enemyImage = ebiten.NewImage(32, 32)
			// 设置缩放选项
			options := &ebiten.DrawImageOptions{}
			// 计算缩放比例
			scaleX := float64(32) / float64(enemyImg.Bounds().Dx())
			scaleY := float64(32) / float64(enemyImg.Bounds().Dy())
			options.GeoM.Scale(scaleX, scaleY)
			// 将原始图片绘制到指定大小的缓冲区
			enemyImage.DrawImage(ebiten.NewImageFromImage(enemyImg), options)
		}
	}

	// 加载玩家图片
	playerImgData, err := imagesFS.ReadFile("resources/images/player.png")
	if err != nil {
		log.Printf("无法加载玩家图片，使用默认方块: %v", err)
	} else {
		// 从嵌入的资源创建图片
		playerImg, _, err := image.Decode(bytes.NewReader(playerImgData))
		if err != nil {
			log.Printf("无法解码玩家图片，使用默认方块: %v", err)
		} else {
			// 创建指定大小的图片缓冲区
			playerImage = ebiten.NewImage(32, 32)
			// 设置缩放选项
			options := &ebiten.DrawImageOptions{}
			// 计算缩放比例
			scaleX := float64(32) / float64(playerImg.Bounds().Dx())
			scaleY := float64(32) / float64(playerImg.Bounds().Dy())
			options.GeoM.Scale(scaleX, scaleY)
			// 将原始图片绘制到指定大小的缓冲区
			playerImage.DrawImage(ebiten.NewImageFromImage(playerImg), options)
		}
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
	// 启动动画相关字段
	animTimer      int          // 动画计时器
	titleScale     float64      // 标题缩放
	titleRotation  float64      // 标题旋转角度
	titleAlpha     float64      // 标题透明度
	titleY         float64      // 标题Y坐标
	menuItemsAlpha float64      // 菜单项透明度
	starPositions  [][2]float64 // 背景星星位置
}

// Update 处理游戏逻辑更新
func (g *Game) Update() error {
	// 在菜单模式下处理模式选择和动画效果
	if g.gameMode == ModeMenu {
		// 更新动画计时器
		g.animTimer++

		// 标题动画效果
		if g.titleScale < 1.0 {
			g.titleScale += 0.05
			if g.titleScale > 1.0 {
				g.titleScale = 1.0
			}
		}

		// 标题旋转效果
		if g.animTimer < 30 {
			g.titleRotation = math.Sin(float64(g.animTimer)/10) * 0.1
		} else {
			g.titleRotation = math.Sin(float64(g.animTimer)/100) * 0.03
		}

		// 标题透明度渐变
		if g.titleAlpha < 1.0 {
			g.titleAlpha += 0.03
			if g.titleAlpha > 1.0 {
				g.titleAlpha = 1.0
			}
		}

		// 菜单项目透明度渐变，在标题出现后
		if g.titleAlpha >= 0.8 && g.menuItemsAlpha < 1.0 {
			g.menuItemsAlpha += 0.03
			if g.menuItemsAlpha > 1.0 {
				g.menuItemsAlpha = 1.0
			}
		}

		// 按1选择关卡模式或鼠标点击
		if (ebiten.IsKeyPressed(ebiten.Key1) && g.menuItemsAlpha >= 0.9) || (ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.checkMouseInArea(screenWidth/2-150, screenHeight/2+20, 320, 40)) {
			g.gameMode = ModeLevels
			g.currentLevel = 1
			g.targetScore = g.currentLevel * 1000 // 每关目标分数为关卡数 * 1000
			g.enemyManager.SetLevel(g.currentLevel)
			return nil
		}
		// 按2选择无尽模式或鼠标点击
		if (ebiten.IsKeyPressed(ebiten.Key2) && g.menuItemsAlpha >= 0.9) || (ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && g.checkMouseInArea(screenWidth/2-150, screenHeight/2+70, 320, 40)) {
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
				enemy.health -= 1      // 减少敌机血量
				if enemy.health <= 0 { // 只有当血量为0时才销毁敌机
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

// checkMouseInArea 检测鼠标是否在指定区域内
func (g *Game) checkMouseInArea(x, y, width, height float64) bool {
	mx, my := ebiten.CursorPosition()
	return float64(mx) >= x && float64(mx) <= x+width && float64(my) >= y && float64(my) <= y+height
}

// Draw 处理游戏画面渲染
func (g *Game) Draw(screen *ebiten.Image) {
	if g.gameMode == ModeMenu {
		// 绘制渐变背景
		gradientTop := color.RGBA{10, 10, 50, 255}
		gradientBottom := color.RGBA{30, 30, 80, 255}

		// 绘制背景渐变色
		for y := 0; y < screenHeight; y++ {
			// 计算该行的渐变颜色
			ratio := float64(y) / float64(screenHeight)
			r := uint8(float64(gradientTop.R) + ratio*float64(gradientBottom.R-gradientTop.R))
			g := uint8(float64(gradientTop.G) + ratio*float64(gradientBottom.G-gradientTop.G))
			b := uint8(float64(gradientTop.B) + ratio*float64(gradientBottom.B-gradientTop.B))
			ebitenutil.DrawRect(screen, 0, float64(y), float64(screenWidth), 1, color.RGBA{r, g, b, 255})
		}

		// 绘制动态星星背景
		for i, star := range g.starPositions {
			// 星星大小和亮度随时间变化
			starSize := 1.0 + math.Sin(float64(g.animTimer)/20.0+float64(i))*0.5
			starBrightness := 150.0 + math.Sin(float64(g.animTimer)/15.0+float64(i))*50.0

			// 让星星闪烁
			starColor := color.RGBA{255, 255, 255, uint8(starBrightness)}
			ebitenutil.DrawRect(screen, star[0], star[1], starSize, starSize, starColor)
		}

		// 绘制标题
		titleMsg := "飞机大战"
		// 计算文本宽度以便居中显示
		titleWidth := len([]rune(titleMsg)) * 36 // 估计每个汉字宽度约36像素
		//titleX := screenWidth/2 - titleWidth/2
		titleOptions := &ebiten.DrawImageOptions{}
		titleWidth = len([]rune(titleMsg)) * 36 // 估计每个汉字宽度

		titleImgWidth := titleWidth + 60 // 两侧各留30像素空间

		// 设置标题的变换
		titleOptions.GeoM.Translate(-float64(titleImgWidth/2), -30)   // 调整旋转中心点到图像中心
		titleOptions.GeoM.Rotate(g.titleRotation)                     // 应用旋转
		titleOptions.GeoM.Scale(g.titleScale, g.titleScale)           // 应用缩放
		titleOptions.GeoM.Translate(float64(screenWidth/2), g.titleY) // 移动到屏幕中心
		titleOptions.ColorM.Scale(1, 1, 1, g.titleAlpha)              // 应用透明度
		titleOptions.ColorM.Scale(1, 1, 1, g.titleAlpha)              // 应用透明度

		// 创建一个临时图像来渲染标题
		titleImg := ebiten.NewImage(titleImgWidth, 60)

		// 绘制闪亮的标题背景
		backgroundWidth := float64(titleImgWidth - 20) // 两侧各留10像素
		ebitenutil.DrawRect(titleImg, 10, 10, backgroundWidth, 50, color.RGBA{0, 0, 150, 200})

		// 添加标题光晕效果
		glowSize := 5.0 + math.Sin(float64(g.animTimer)/10.0)*2.0
		ebitenutil.DrawRect(titleImg, 10-glowSize, 10-glowSize, backgroundWidth+glowSize*2, 50+glowSize*2, color.RGBA{100, 100, 255, 50})

		// 直接使用已有的中文字体渲染标题
		// 计算文本在图像中的X位置，使其居中
		textX := titleImgWidth/2 - titleWidth/2
		// 绘制阴影效果
		text.Draw(titleImg, titleMsg, chineseFont, textX+2, 47, color.RGBA{0, 0, 80, 255})
		// 主文字
		text.Draw(titleImg, titleMsg, chineseFont, textX, 45, color.RGBA{255, 255, 255, 255})
		// 添加发光效果
		text.Draw(titleImg, titleMsg, chineseFont, textX-1, 44, color.RGBA{100, 200, 255, 150})

		// 将标题绘制到屏幕
		screen.DrawImage(titleImg, titleOptions)

		// 使用透明度来控制菜单项的显示
		menuAlpha := uint8(g.menuItemsAlpha * 255)

		// 绘制模式选择说明
		modeTitle := "- 选择游戏模式 -"
		modeTitleX := screenWidth/2 - 100
		modeTitleY := int(g.titleY) + 120

		// 模式选择背景带有呼吸效果
		pulseEffect := 0.7 + math.Sin(float64(g.animTimer)/20.0)*0.3
		bgWidth := 240.0 * pulseEffect
		ebitenutil.DrawRect(screen, float64(modeTitleX-20)-(bgWidth-200)/2, float64(modeTitleY-25), bgWidth, 35,
			color.RGBA{0, 50, 150, menuAlpha})
		text.Draw(screen, modeTitle, chineseFont, modeTitleX, modeTitleY,
			color.RGBA{220, 220, 255, menuAlpha})

		// 绘制模式选项
		mode1 := "[1] 关卡模式"
		mode1Desc := "逐级挑战，难度递增"
		mode2 := "[2] 无尽模式"
		mode2Desc := "无限挑战，直到失败"

		mode1X := screenWidth/2 - 150
		mode1Y := modeTitleY + 60
		mode2X := screenWidth/2 - 150
		mode2Y := modeTitleY + 110

		// 高亮效果（鼠标或键盘悬停时）
		mx, my := ebiten.CursorPosition()
		mode1Highlight := (mx >= mode1X-10 && mx <= mode1X+310 && my >= mode1Y-25 && my <= mode1Y+15)
		mode2Highlight := (mx >= mode2X-10 && mx <= mode2X+310 && my >= mode2Y-25 && my <= mode2Y+15)

		// 模式选项背景
		mode1BgColor := color.RGBA{0, 0, 100, menuAlpha}
		mode2BgColor := color.RGBA{0, 0, 100, menuAlpha}

		// 为选项添加悬停高亮效果
		if mode1Highlight {
			mode1BgColor = color.RGBA{50, 50, 150, menuAlpha}
		}
		if mode2Highlight {
			mode2BgColor = color.RGBA{50, 50, 150, menuAlpha}
		}

		ebitenutil.DrawRect(screen, float64(mode1X-10), float64(mode1Y-25), 320, 40, mode1BgColor)
		ebitenutil.DrawRect(screen, float64(mode2X-10), float64(mode2Y-25), 320, 40, mode2BgColor)

		// 绘制模式选项文字
		text.Draw(screen, mode1, chineseFont, mode1X, mode1Y, color.RGBA{255, 255, 0, menuAlpha})
		text.Draw(screen, mode1Desc, chineseFont, mode1X+140, mode1Y, color.RGBA{255, 255, 255, menuAlpha})
		text.Draw(screen, mode2, chineseFont, mode2X, mode2Y, color.RGBA{255, 255, 0, menuAlpha})
		text.Draw(screen, mode2Desc, chineseFont, mode2X+140, mode2Y, color.RGBA{255, 255, 255, menuAlpha})

		// 绘制飞机小图标在当前选择的模式旁边
		if mode1Highlight || mode2Highlight {
			planeOptions := &ebiten.DrawImageOptions{}
			// 添加小飞机图标的飘动效果
			planeY := 0.0
			if mode1Highlight {
				planeY = float64(mode1Y - 15)
			} else {
				planeY = float64(mode2Y - 15)
			}
			// 让飞机小图标左右摆动
			planeX := float64(mode1X-40) + math.Sin(float64(g.animTimer)/10.0)*5.0
			planeOptions.GeoM.Scale(0.6, 0.6) // 缩小图标
			planeOptions.GeoM.Translate(planeX, planeY)
			screen.DrawImage(playerImage, planeOptions)
		}

		// 绘制操作提示
		hint := "按对应数字键或点击选择模式"
		hintX := screenWidth/2 - 140
		hintY := mode2Y + 70

		// 提示背景带有呼吸效果
		hintPulse := 0.8 + math.Sin(float64(g.animTimer+30)/20.0)*0.2
		hintBgWidth := 280.0 * hintPulse
		ebitenutil.DrawRect(screen, float64(hintX-20)-(hintBgWidth-280)/2, float64(hintY-25), hintBgWidth, 35,
			color.RGBA{0, 0, 100, uint8(float64(menuAlpha) * 0.7)})
		text.Draw(screen, hint, chineseFont, hintX, hintY,
			color.RGBA{180, 180, 255, menuAlpha})

		// 添加版本信息和作者信息
		versionText := "版本: v1.0.0"
		authorText := "© 2025 飞机大战开发团队"
		text.Draw(screen, versionText, chineseFont, 10, screenHeight-30,
			color.RGBA{200, 200, 200, menuAlpha})
		text.Draw(screen, authorText, chineseFont, screenWidth-240, screenHeight-30,
			color.RGBA{200, 200, 200, menuAlpha})
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

	// 创建随机星星背景
	starPositions := make([][2]float64, 100)
	for i := 0; i < 100; i++ {
		starPositions[i] = [2]float64{float64(rand.Intn(screenWidth)), float64(rand.Intn(screenHeight))}
	}

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
		// 初始化动画参数
		animTimer:      0,
		titleScale:     0.1,
		titleRotation:  0.0,
		titleAlpha:     0.0,
		titleY:         float64(screenHeight) / 4,
		menuItemsAlpha: 0.0,
		starPositions:  starPositions,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
