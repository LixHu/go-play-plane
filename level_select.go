package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

// LevelSelectMenu 关卡选择菜单
type LevelSelectMenu struct {
	levels           int         // 可选关卡数量
	currentSelection int         // 当前选中的关卡
	animTimer        int         // 动画计时器
	titleScale       float64     // 标题缩放
	titleRotation    float64     // 标题旋转
	levelsAlpha      float64     // 关卡透明度
	ready            bool        // 是否已准备好
	levelInfos       []levelInfo // 关卡信息
}

// levelInfo 关卡信息结构体
type levelInfo struct {
	name        string   // 关卡名称
	description string   // 关卡描述
	bossName    string   // BOSS名称
	bossType    BossType // BOSS类型
	difficulty  int      // 难度（1-5星）
	locked      bool     // 是否锁定
}

// NewLevelSelectMenu 创建一个新的关卡选择菜单
func NewLevelSelectMenu() *LevelSelectMenu {
	// 创建关卡信息
	levelInfos := []levelInfo{
		{
			name:        "第1关：初入战场",
			description: "遭遇第一个BOSS，熟悉控制",
			bossName:    "环形魔王",
			bossType:    BossType1,
			difficulty:  1,
			locked:      false,
		},
		{
			name:        "第2关：交叉火力",
			description: "小心交叉弹幕的包围",
			bossName:    "十字统领",
			bossType:    BossType2,
			difficulty:  2,
			locked:      true,
		},
		{
			name:        "第3关：追踪猎手",
			description: "BOSS会发射追踪弹幕",
			bossName:    "追猎者",
			bossType:    BossType3,
			difficulty:  3,
			locked:      true,
		},
		{
			name:        "第4关：混沌风暴",
			description: "终极BOSS，混合所有弹幕类型",
			bossName:    "混沌大帝",
			bossType:    BossType4,
			difficulty:  5,
			locked:      true,
		},
	}

	return &LevelSelectMenu{
		levels:           len(levelInfos),
		currentSelection: 0,
		animTimer:        0,
		titleScale:       0.1,
		titleRotation:    0.0,
		levelsAlpha:      0.0,
		ready:            false,
		levelInfos:       levelInfos,
	}
}

// Update 更新关卡选择菜单
func (lsm *LevelSelectMenu) Update(game *Game) bool {
	lsm.animTimer++

	// 标题动画效果
	if lsm.titleScale < 1.0 {
		lsm.titleScale += 0.05
		if lsm.titleScale > 1.0 {
			lsm.titleScale = 1.0
		}
	}

	// 标题旋转效果
	if lsm.animTimer < 30 {
		lsm.titleRotation = math.Sin(float64(lsm.animTimer)/10) * 0.1
	} else {
		lsm.titleRotation = math.Sin(float64(lsm.animTimer)/100) * 0.03
	}

	// 关卡透明度渐变
	if lsm.titleScale >= 0.9 && lsm.levelsAlpha < 1.0 {
		lsm.levelsAlpha += 0.03
		if lsm.levelsAlpha > 1.0 {
			lsm.levelsAlpha = 1.0
			lsm.ready = true
		}
	}

	// 只有在准备好后才能选择关卡
	if lsm.ready {
		// 键盘操作
		if ebiten.IsKeyPressed(ebiten.KeyLeft) && !ebiten.IsKeyPressed(ebiten.KeyRight) {
			if lsm.animTimer%10 == 0 { // 降低移动速度
				lsm.currentSelection--
				if lsm.currentSelection < 0 {
					lsm.currentSelection = lsm.levels - 1
				}
			}
		} else if ebiten.IsKeyPressed(ebiten.KeyRight) && !ebiten.IsKeyPressed(ebiten.KeyLeft) {
			if lsm.animTimer%10 == 0 { // 降低移动速度
				lsm.currentSelection++
				if lsm.currentSelection >= lsm.levels {
					lsm.currentSelection = 0
				}
			}
		}

		// 鼠标操作
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()

			// 检查点击的是哪个关卡
			for i := 0; i < lsm.levels; i++ {
				// 计算每行最多显示5个关卡
				row := i / 5
				col := i % 5

				x := 70 + col*100
				y := 150 + row*90

				if float64(mx) >= float64(x) && float64(mx) <= float64(x+70) &&
					float64(my) >= float64(y) && float64(my) <= float64(y+70) {
					lsm.currentSelection = i
					break
				}
			}

			// 检查是否点击了开始按钮
			startX := screenWidth/2 - 100
			startY := screenHeight/2 + 150
			if float64(mx) >= float64(startX) && float64(mx) <= float64(startX+200) &&
				float64(my) >= float64(startY-30) && float64(my) <= float64(startY+10) &&
				!lsm.levelInfos[lsm.currentSelection].locked {
				// 开始所选关卡
				game.currentLevel = lsm.currentSelection + 1
				game.targetScore = game.currentLevel * 1000
				game.enemyManager.SetLevel(game.currentLevel)
				game.bossActive = false // 初始没有BOSS
				game.bossDefeated = false
				game.gameMode = ModePlaying // 切换到游戏模式
				return true
			}

			// 检查是否点击了返回按钮
			backX := 80
			backY := 40
			if float64(mx) >= float64(backX-60) && float64(mx) <= float64(backX+60) &&
				float64(my) >= float64(backY-25) && float64(my) <= float64(backY+15) {
				// 返回主菜单
				game.gameMode = ModeMenu
				return true
			}
		}

		// 空格或回车键确认选择
		if (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter)) &&
			!lsm.levelInfos[lsm.currentSelection].locked {
			// 开始所选关卡
			game.currentLevel = lsm.currentSelection + 1
			game.targetScore = game.currentLevel * 1000
			game.enemyManager.SetLevel(game.currentLevel)
			game.bossActive = false // 初始没有BOSS
			game.bossDefeated = false
			game.gameMode = ModePlaying // 切换到游戏模式
			return true
		}

		// ESC键返回主菜单
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			game.gameMode = ModeMenu
			return true
		}
	}

	// 解锁已完成关卡的下一关
	for i := 0; i < lsm.levels-1; i++ {
		if !lsm.levelInfos[i].locked {
			lsm.levelInfos[i+1].locked = false
		}
	}

	return false
}

// Draw 绘制关卡选择菜单
func (lsm *LevelSelectMenu) Draw(screen *ebiten.Image) {
	// 绘制渐变背景
	gradientTop := color.RGBA{20, 20, 60, 255}
	gradientBottom := color.RGBA{40, 40, 100, 255}

	// 绘制背景渐变色
	for y := 0; y < screenHeight; y++ {
		// 计算该行的渐变颜色
		ratio := float64(y) / float64(screenHeight)
		r := uint8(float64(gradientTop.R) + ratio*float64(gradientBottom.R-gradientTop.R))
		g := uint8(float64(gradientTop.G) + ratio*float64(gradientBottom.G-gradientTop.G))
		b := uint8(float64(gradientTop.B) + ratio*float64(gradientBottom.B-gradientTop.B))
		ebitenutil.DrawRect(screen, 0, float64(y), float64(screenWidth), 1, color.RGBA{r, g, b, 255})
	}

	// 绘制标题
	titleMsg := "选择关卡"
	titleOptions := &ebiten.DrawImageOptions{}
	titleWidth := len([]rune(titleMsg)) * 36 // 估计每个汉字宽度约36像素
	titleImgWidth := titleWidth + 60         // 两侧各留30像素空间

	// 设置标题的变换
	titleOptions.GeoM.Translate(-float64(titleImgWidth/2), -30) // 调整旋转中心点
	titleOptions.GeoM.Rotate(lsm.titleRotation)                 // 应用旋转
	titleOptions.GeoM.Scale(lsm.titleScale, lsm.titleScale)     // 应用缩放
	titleOptions.GeoM.Translate(float64(screenWidth/2), 70)     // 移动到屏幕上方

	// 创建临时图像来渲染标题
	titleImg := ebiten.NewImage(titleImgWidth, 60)

	// 绘制标题背景
	backgroundWidth := float64(titleImgWidth - 20)
	ebitenutil.DrawRect(titleImg, 10, 10, backgroundWidth, 50, color.RGBA{0, 0, 150, 200})

	// 添加标题光晕效果
	glowSize := 5.0 + math.Sin(float64(lsm.animTimer)/10.0)*2.0
	ebitenutil.DrawRect(titleImg, 10-glowSize, 10-glowSize, backgroundWidth+glowSize*2, 50+glowSize*2, color.RGBA{100, 100, 255, 50})

	// 直接使用中文字体渲染标题
	textX := titleImgWidth/2 - titleWidth/2
	// 绘制阴影效果
	text.Draw(titleImg, titleMsg, chineseFont, textX+2, 47, color.RGBA{0, 0, 80, 255})
	// 主文字
	text.Draw(titleImg, titleMsg, chineseFont, textX, 45, color.RGBA{255, 255, 255, 255})
	// 添加发光效果
	text.Draw(titleImg, titleMsg, chineseFont, textX-1, 44, color.RGBA{100, 200, 255, 150})

	// 将标题绘制到屏幕
	screen.DrawImage(titleImg, titleOptions)

	// 使用透明度控制关卡项的显示
	menuAlpha := uint8(lsm.levelsAlpha * 255)

	// 绘制返回按钮
	backX := 80
	backY := 40
	ebitenutil.DrawRect(screen, float64(backX-60), float64(backY-25), 120, 40, color.RGBA{0, 0, 100, menuAlpha})
	text.Draw(screen, "返回", chineseFont, backX-30, backY, color.RGBA{255, 255, 255, menuAlpha})

	// 绘制关卡选项
	for i := 0; i < lsm.levels; i++ {
		levelInfo := lsm.levelInfos[i]

		// 计算每行最多显示5个关卡
		row := i / 5
		col := i % 5

		levelX := 70 + col*100
		levelY := 150 + row*90

		// 判断是否是当前选中的关卡
		isSelected := (i == lsm.currentSelection)

		// 关卡选项背景
		var levelBgColor color.RGBA
		if levelInfo.locked {
			// 锁定的关卡显示灰色
			levelBgColor = color.RGBA{100, 100, 100, menuAlpha}
		} else if isSelected {
			// 选中的关卡显示亮蓝色
			levelBgColor = color.RGBA{0, 100, 200, menuAlpha}

			// 为选中关卡添加呼吸效果
			pulseEffect := 1.0 + math.Sin(float64(lsm.animTimer)/10.0)*0.1
			pulseSize := 70.0 * pulseEffect
			ebitenutil.DrawRect(screen,
				float64(levelX)-(pulseSize-70)/2,
				float64(levelY)-(pulseSize-70)/2,
				pulseSize, pulseSize,
				color.RGBA{0, 150, 255, uint8(menuAlpha / 3)})
		} else {
			// 未选中的关卡显示深蓝色
			levelBgColor = color.RGBA{0, 50, 150, menuAlpha}
		}

		// 绘制关卡图标背景
		ebitenutil.DrawRect(screen, float64(levelX), float64(levelY), 70, 70, levelBgColor)

		// 绘制关卡数字
		levelNumText := fmt.Sprintf("%d", i+1)
		numX := levelX + 30 - len(levelNumText)*5
		text.Draw(screen, levelNumText, chineseFont, numX, levelY+40, color.RGBA{255, 255, 0, menuAlpha})

		// 对于锁定的关卡，绘制锁定图标
		if levelInfo.locked {
			// 绘制锁图标
			lockX := levelX + 20
			lockY := levelY + 20
			// 锁顶部
			ebitenutil.DrawRect(screen, float64(lockX), float64(lockY-15), 30, 20, color.RGBA{200, 200, 200, menuAlpha})
			// 锁身体
			ebitenutil.DrawRect(screen, float64(lockX+5), float64(lockY+5), 20, 25, color.RGBA{200, 200, 200, menuAlpha})
			// 锁孔
			ebitenutil.DrawRect(screen, float64(lockX+12), float64(lockY+15), 6, 10, color.RGBA{50, 50, 50, menuAlpha})
		}
	}

	// 绘制当前选中关卡的信息
	if lsm.currentSelection >= 0 && lsm.currentSelection < lsm.levels {
		levelInfo := lsm.levelInfos[lsm.currentSelection]

		// 信息面板背景
		infoPanelX := screenWidth/2 - 200
		infoPanelY := screenHeight/2 + 70
		ebitenutil.DrawRect(screen, float64(infoPanelX), float64(infoPanelY), 400, 60, color.RGBA{0, 0, 100, menuAlpha})

		// 关卡名称
		text.Draw(screen, levelInfo.name, chineseFont, infoPanelX+10, infoPanelY+25, color.RGBA{255, 255, 0, menuAlpha})

		// 关卡描述
		text.Draw(screen, levelInfo.description, chineseFont, infoPanelX+10, infoPanelY+50, color.RGBA{255, 255, 255, menuAlpha})

		// BOSS信息面板
		bossInfoX := screenWidth/2 - 150
		bossInfoY := screenHeight/2 + 10
		bossTitleText := fmt.Sprintf("BOSS: %s", levelInfo.bossName)
		text.Draw(screen, bossTitleText, chineseFont, bossInfoX, bossInfoY, color.RGBA{255, 50, 50, menuAlpha})

		// 难度星级
		difficultyText := "难度: "
		text.Draw(screen, difficultyText, chineseFont, infoPanelX+270, infoPanelY+25, color.RGBA{255, 255, 255, menuAlpha})

		// 绘制星星
		for i := 0; i < 5; i++ {
			starX := infoPanelX + 330 + i*15
			starY := infoPanelY + 22

			if i < levelInfo.difficulty {
				// 点亮的星星
				starColor := color.RGBA{255, 255, 0, menuAlpha}
				ebitenutil.DrawRect(screen, float64(starX), float64(starY), 10, 10, starColor)
			} else {
				// 未点亮的星星
				starColor := color.RGBA{100, 100, 100, menuAlpha}
				ebitenutil.DrawRect(screen, float64(starX), float64(starY), 10, 10, starColor)
			}
		}
	}

	// 绘制开始按钮
	startX := screenWidth/2 - 100
	startY := screenHeight/2 + 150

	// 判断当前选择的关卡是否解锁
	canStart := !lsm.levelInfos[lsm.currentSelection].locked

	// 根据是否可以开始设置按钮颜色
	startBtnColor := color.RGBA{100, 100, 100, menuAlpha} // 默认灰色（锁定状态）
	if canStart {
		// 可开始时呈现绿色，带有呼吸效果
		pulseValue := 50 + int(math.Sin(float64(lsm.animTimer)/10.0)*20.0)
		startBtnColor = color.RGBA{0, 150 + uint8(pulseValue), 0, menuAlpha}
	}

	ebitenutil.DrawRect(screen, float64(startX), float64(startY-30), 200, 40, startBtnColor)
	startText := "开始游戏"
	if !canStart {
		startText = "关卡锁定"
	}
	startTextX := startX + 100 - len([]rune(startText))*12 // 居中文本
	text.Draw(screen, startText, chineseFont, startTextX, startY, color.RGBA{255, 255, 255, menuAlpha})

	// 操作提示
	hintText := "← → 选择关卡   空格/回车 确认   ESC 返回"
	hintX := screenWidth/2 - len([]rune(hintText))*10
	hintY := screenHeight - 30
	text.Draw(screen, hintText, chineseFont, hintX, hintY, color.RGBA{200, 200, 200, menuAlpha})
}
