package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// BossType 表示BOSS类型
type BossType int

const (
	BossType1 BossType = iota // 第一关BOSS：环形弹幕
	BossType2                 // 第二关BOSS：交叉弹幕
	BossType3                 // 第三关BOSS：追踪弹幕
	BossType4                 // 第四关BOSS：混合弹幕
)

// Boss 表示关卡BOSS
type Boss struct {
	x           float64
	y           float64
	speedX      float64
	speedY      float64
	width       int
	height      int
	active      bool
	health      int // 当前血量
	maxHealth   int // 最大血量
	bossType    BossType
	phase       int  // 当前阶段，血量降低时进入下一阶段，难度增加
	animTimer   int  // 动画计时器
	shootTimer  int  // 射击计时器
	patternTime int  // 弹幕模式切换计时器
	enterScene  bool // 是否正在入场
}

// NewBoss 创建一个新的BOSS
func NewBoss(bossType BossType) *Boss {
	// 根据BOSS类型设置不同的初始值
	var health int
	var width, height int

	switch bossType {
	case BossType1:
		health = 200
		width = 80
		height = 80
	case BossType2:
		health = 300
		width = 100
		height = 80
	case BossType3:
		health = 400
		width = 100
		height = 100
	case BossType4:
		health = 500
		width = 120
		height = 100
	}

	return &Boss{
		x:           float64(screenWidth/2 - width/2),
		y:           -float64(height), // 从屏幕上方进入
		speedX:      1.0,
		speedY:      1.0,
		width:       width,
		height:      height,
		active:      true,
		health:      health,
		maxHealth:   health,
		bossType:    bossType,
		phase:       1,
		animTimer:   0,
		shootTimer:  0,
		patternTime: 0,
		enterScene:  true,
	}
}

// Update 更新BOSS状态
func (b *Boss) Update(player *Player, bulletManager *EnemyBulletManager) {
	b.animTimer++
	b.patternTime++

	// 入场动画
	if b.enterScene {
		if b.y < 80 {
			b.y += 2
		} else {
			b.enterScene = false
		}
		return
	}

	// 根据血量更新阶段
	healthPercent := float64(b.health) / float64(b.maxHealth)
	if healthPercent <= 0.75 && b.phase == 1 {
		b.phase = 2
	} else if healthPercent <= 0.5 && b.phase == 2 {
		b.phase = 3
	} else if healthPercent <= 0.25 && b.phase == 3 {
		b.phase = 4
	}

	// BOSS移动模式
	switch b.bossType {
	case BossType1:
		// 第一关BOSS：在屏幕上方左右移动
		b.x += b.speedX
		if b.x <= 0 || b.x+float64(b.width) >= float64(screenWidth) {
			b.speedX = -b.speedX
		}

	case BossType2:
		// 第二关BOSS：正弦移动
		b.x += b.speedX
		if b.x <= 0 || b.x+float64(b.width) >= float64(screenWidth) {
			b.speedX = -b.speedX
		}
		b.y = 80 + math.Sin(float64(b.animTimer)/30.0)*40.0

	case BossType3:
		// 第三关BOSS：追踪玩家
		targetX := player.x + float64(player.width/2) - float64(b.width/2)
		targetX = math.Max(0, math.Min(targetX, float64(screenWidth-b.width)))

		if b.x < targetX {
			b.x += b.speedX
		} else if b.x > targetX {
			b.x -= b.speedX
		}

		// 保持在一定距离内
		b.y = 80 + math.Sin(float64(b.animTimer)/40.0)*30.0

	case BossType4:
		// 第四关BOSS：随机突进模式
		if b.patternTime > 180 { // 每3秒随机改变运动方向
			b.speedX = rand.Float64()*4.0 - 2.0
			b.speedY = rand.Float64()*2.0 - 1.0
			b.patternTime = 0
		}

		b.x += b.speedX
		b.y += b.speedY

		// 边界检查
		if b.x <= 0 || b.x+float64(b.width) >= float64(screenWidth) {
			b.speedX = -b.speedX
		}

		if b.y <= 30 || b.y+float64(b.height) >= float64(screenHeight)/2 {
			b.speedY = -b.speedY
		}
	}

	// 发射子弹
	b.shootTimer++
	shootInterval := 30 // 基础射击间隔
	// 不同阶段降低射击间隔（增加射击频率）
	shootInterval = shootInterval - (b.phase-1)*5
	shootInterval = max(10, shootInterval) // 最小间隔10帧

	if b.shootTimer >= shootInterval {
		// 根据BOSS类型和阶段发射不同弹幕
		switch b.bossType {
		case BossType1:
			b.fireCirclePattern(bulletManager, b.phase)

		case BossType2:
			b.fireCrossPattern(bulletManager, b.phase)

		case BossType3:
			b.fireHomingPattern(bulletManager, player, b.phase)

		case BossType4:
			// 混合弹幕，随机使用其他BOSS的弹幕
			pattern := rand.Intn(3)
			switch pattern {
			case 0:
				b.fireCirclePattern(bulletManager, b.phase)
			case 1:
				b.fireCrossPattern(bulletManager, b.phase)
			case 2:
				b.fireHomingPattern(bulletManager, player, b.phase)
			}
		}

		b.shootTimer = 0
	}
}

// fireCirclePattern 发射环形弹幕
func (b *Boss) fireCirclePattern(bulletManager *EnemyBulletManager, phase int) {
	// 发射点在BOSS中心
	centerX := b.x + float64(b.width)/2
	centerY := b.y + float64(b.height)/2

	// 环形弹幕的子弹数量，随阶段增加
	bulletCount := 8 + (phase-1)*2

	// 计算每个子弹的角度
	for i := 0; i < bulletCount; i++ {
		angle := float64(i) * (360.0 / float64(bulletCount))
		radian := angle * math.Pi / 180.0

		// 创建子弹，设置速度方向
		bullet := NewEnemyBulletCustom(centerX, centerY, radian, 3.0, color.RGBA{255, 50, 50, 255})
		bulletManager.bullets = append(bulletManager.bullets, bullet)
	}

	// 高级阶段增加额外的交错环
	if phase >= 3 {
		// 第二环直接发射，不使用goroutine（简化处理）
		// 偏移角度的第二波
		for i := 0; i < bulletCount; i++ {
			angle := float64(i)*(360.0/float64(bulletCount)) + 180.0/float64(bulletCount)
			radian := angle * math.Pi / 180.0

			// 创建子弹，设置速度方向
			bullet := NewEnemyBulletCustom(centerX, centerY, radian, 3.0, color.RGBA{255, 150, 50, 255})
			bulletManager.bullets = append(bulletManager.bullets, bullet)
		}
	}
}

// fireCrossPattern 发射交叉弹幕
func (b *Boss) fireCrossPattern(bulletManager *EnemyBulletManager, phase int) {
	// 发射点在BOSS中心
	centerX := b.x + float64(b.width)/2
	centerY := b.y + float64(b.height)/2

	// 交叉线数量，随阶段增加
	lineCount := 2 + (phase - 1)
	lineCount = min(lineCount, 6) // 最多6条线

	// 每条线上的子弹数量
	bulletsPerLine := 5

	// 计算每条线的角度
	for i := 0; i < lineCount; i++ {
		angle := float64(i) * (180.0 / float64(lineCount))
		radian := angle * math.Pi / 180.0

		// 在每条线上创建多个子弹
		for j := 0; j < bulletsPerLine; j++ {
			// 计算不同速度，形成一条线
			speed := 2.0 + float64(j)*0.5

			// 创建子弹，设置速度方向
			bullet := NewEnemyBulletCustom(centerX, centerY, radian, speed, color.RGBA{0, 150, 255, 255})
			bulletManager.bullets = append(bulletManager.bullets, bullet)
		}
	}

	// 高级阶段添加旋转效果
	if phase >= 3 {
		rotationOffset := float64(b.animTimer%360) * math.Pi / 180.0

		// 额外发射一组旋转的子弹
		for i := 0; i < lineCount; i++ {
			angle := float64(i)*(180.0/float64(lineCount)) + rotationOffset
			radian := angle * math.Pi / 180.0

			// 创建旋转的子弹
			bullet := NewEnemyBulletCustom(centerX, centerY, radian, 3.0, color.RGBA{200, 100, 255, 255})
			bulletManager.bullets = append(bulletManager.bullets, bullet)
		}
	}
}

// fireHomingPattern 发射追踪弹幕
func (b *Boss) fireHomingPattern(bulletManager *EnemyBulletManager, player *Player, phase int) {
	// 发射点在BOSS中心
	centerX := b.x + float64(b.width)/2
	centerY := b.y + float64(b.height)/2

	// 玩家位置
	playerX := player.x + float64(player.width)/2
	playerY := player.y + float64(player.height)/2

	// 计算到玩家的角度
	dx := playerX - centerX
	dy := playerY - centerY
	angle := math.Atan2(dy, dx)

	// 追踪子弹数量随阶段增加
	homingCount := 1 + (phase - 1)

	// 发射多个追踪子弹
	for i := 0; i < homingCount; i++ {
		// 添加一点角度偏移，使子弹有些散布
		angleOffset := (rand.Float64() - 0.5) * 0.5

		// 创建追踪子弹
		bullet := NewEnemyBulletHoming(centerX, centerY, angle+angleOffset, color.RGBA{255, 255, 100, 255})
		bulletManager.bullets = append(bulletManager.bullets, bullet)
	}

	// 高级阶段添加分散攻击
	if phase >= 2 {
		// 额外发射扇形弹幕
		spreadCount := 3 + (phase-2)*2 // 扇形中的子弹数量
		spreadAngle := math.Pi / 3.0   // 60度扇形

		for i := 0; i < spreadCount; i++ {
			spreadOffset := spreadAngle * (float64(i)/float64(spreadCount-1) - 0.5)
			bullet := NewEnemyBulletCustom(centerX, centerY, angle+spreadOffset, 2.5, color.RGBA{255, 200, 0, 255})
			bulletManager.bullets = append(bulletManager.bullets, bullet)
		}
	}
}

// Draw 绘制BOSS
func (b *Boss) Draw(screen *ebiten.Image) {
	// 绘制BOSS图像
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Scale(float64(b.width)/32.0, float64(b.height)/32.0) // 缩放到指定大小
	options.GeoM.Translate(b.x, b.y)

	// 根据BOSS类型设置不同的颜色
	switch b.bossType {
	case BossType1:
		options.ColorM.Scale(1.0, 0.5, 0.5, 1.0) // 红色调
	case BossType2:
		options.ColorM.Scale(0.5, 0.5, 1.0, 1.0) // 蓝色调
	case BossType3:
		options.ColorM.Scale(0.5, 1.0, 0.5, 1.0) // 绿色调
	case BossType4:
		options.ColorM.Scale(1.0, 0.8, 0.0, 1.0) // 金色调
	}

	// 添加闪烁效果
	if b.animTimer%10 < 5 && b.phase >= 3 {
		options.ColorM.Scale(1.2, 1.2, 1.2, 1.0) // 高阶段时闪烁发亮
	}

	screen.DrawImage(enemyImage, options)

	// 绘制BOSS血条背景
	bloodBarWidth := float64(screenWidth - 100)
	const bloodBarHeight = 15.0
	ebitenutil.DrawRect(screen, 50, 20, bloodBarWidth, bloodBarHeight, color.RGBA{100, 100, 100, 200})

	// 计算当前血量比例
	healthRatio := float64(b.health) / float64(b.maxHealth)
	healthWidth := bloodBarWidth * healthRatio

	// 根据健康比例变化颜色
	var healthColor color.RGBA
	if healthRatio > 0.75 {
		// 高血量时显示绿色
		healthColor = color.RGBA{0, 255, 0, 255}
	} else if healthRatio > 0.5 {
		// 中高血量显示黄色
		healthColor = color.RGBA{255, 255, 0, 255}
	} else if healthRatio > 0.25 {
		// 中低血量显示橙色
		healthColor = color.RGBA{255, 165, 0, 255}
	} else {
		// 低血量显示红色
		healthColor = color.RGBA{255, 0, 0, 255}
	}

	// 绘制健康部分
	ebitenutil.DrawRect(screen, 50, 20, healthWidth, bloodBarHeight, healthColor)

	// 添加血条装饰效果
	for i := 0; i < 10; i++ {
		if float64(i)/10.0 <= healthRatio {
			// 绘制分段标记
			markerX := 50 + (bloodBarWidth/10.0)*float64(i)
			ebitenutil.DrawRect(screen, markerX, 18, 2, bloodBarHeight+4, color.RGBA{255, 255, 255, 200})
		}
	}
}
