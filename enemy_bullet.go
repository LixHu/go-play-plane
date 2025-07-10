package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
	"math/rand"
)

// EnemyBullet 表示敌机发射的子弹
type EnemyBullet struct {
	x        float64
	y        float64
	speedX   float64
	speedY   float64
	width    int
	height   int
	active   bool
	color    color.RGBA // 子弹颜色
	isHoming bool       // 是否为追踪子弹
}

// NewEnemyBullet 创建一个新的敌机子弹
func NewEnemyBullet(x, y float64) *EnemyBullet {
	// 随机生成发射角度（-60度到60度之间）
	angle := rand.Float64()*120 - 60
	// 将角度转换为弧度
	radian := angle * math.Pi / 180
	// 基础速度
	baseSpeed := 4.0

	return &EnemyBullet{
		x:        x,
		y:        y,
		speedX:   baseSpeed * math.Sin(radian),
		speedY:   baseSpeed * math.Cos(radian),
		width:    4,
		height:   4, // 修改为正方形，便于旋转
		active:   true,
		color:    color.RGBA{255, 0, 0, 255}, // 默认红色
		isHoming: false,
	}
}

// NewEnemyBulletCustom 创建一个自定义方向和速度的敌机子弹
func NewEnemyBulletCustom(x, y, angle, speed float64, bulletColor color.RGBA) *EnemyBullet {
	return &EnemyBullet{
		x:        x,
		y:        y,
		speedX:   speed * math.Cos(angle),
		speedY:   speed * math.Sin(angle),
		width:    6,
		height:   6,
		active:   true,
		color:    bulletColor,
		isHoming: false,
	}
}

// NewEnemyBulletHoming 创建一个追踪玩家的敌机子弹
func NewEnemyBulletHoming(x, y, angle float64, bulletColor color.RGBA) *EnemyBullet {
	return &EnemyBullet{
		x:        x,
		y:        y,
		speedX:   3.0 * math.Cos(angle),
		speedY:   3.0 * math.Sin(angle),
		width:    8,
		height:   8,
		active:   true,
		color:    bulletColor,
		isHoming: true,
	}
}

// Update 更新敌机子弹的状态
func (b *EnemyBullet) Update() {
	// 更新位置
	b.x += b.speedX
	b.y += b.speedY

	// 检查墙壁碰撞
	if b.x <= 0 || b.x+float64(b.width) >= float64(screenWidth) {
		// 水平反弹
		b.speedX = -b.speedX
		// 确保子弹不会卡在墙内
		if b.x <= 0 {
			b.x = 0
		} else {
			b.x = float64(screenWidth) - float64(b.width)
		}
	}

	// 如果飞出屏幕底部，标记为非活动状态
	if b.y > float64(screenHeight) {
		b.active = false
	}

	// 如果飞到屏幕顶部，反弹
	if b.y <= 0 {
		b.speedY = -b.speedY
		b.y = 0
	}
}

// UpdateHoming 更新追踪子弹的状态，使其朝向玩家
func (b *EnemyBullet) UpdateHoming(player *Player) {
	// 只有追踪子弹才执行此逻辑
	if b.isHoming {
		// 计算到玩家的方向
		playerCenterX := player.x + float64(player.width)/2
		playerCenterY := player.y + float64(player.height)/2
		bulletCenterX := b.x + float64(b.width)/2
		bulletCenterY := b.y + float64(b.height)/2

		// 计算方向向量
		dx := playerCenterX - bulletCenterX
		dy := playerCenterY - bulletCenterY

		// 计算距离
		distance := math.Sqrt(dx*dx + dy*dy)

		// 防止除以零
		if distance > 0 {
			// 归一化方向向量
			dx /= distance
			dy /= distance

			// 计算当前速度大小
			speed := math.Sqrt(b.speedX*b.speedX + b.speedY*b.speedY)

			// 缓慢转向玩家（增加追踪性）
			turnFactor := 0.1 // 转向因子，越大转向越快
			b.speedX = b.speedX*(1-turnFactor) + dx*speed*turnFactor
			b.speedY = b.speedY*(1-turnFactor) + dy*speed*turnFactor
		}
	}

	// 调用基本的更新方法
	b.Update()
}

// Draw 绘制敌机子弹
func (b *EnemyBullet) Draw(screen *ebiten.Image) {
	// 使用自定义颜色绘制子弹
	bulletColor := b.color
	// 如果颜色为空值，使用默认红色
	if bulletColor.A == 0 {
		bulletColor = color.RGBA{255, 0, 0, 255}
	}

	// 绘制子弹
	ebitenutil.DrawRect(screen, b.x, b.y, float64(b.width), float64(b.height), bulletColor)

	// 如果是追踪子弹，添加发光效果
	if b.isHoming {
		// 绘制外发光
		glowColor := color.RGBA{bulletColor.R, bulletColor.G, bulletColor.B, 100}
		ebitenutil.DrawRect(screen, b.x-2, b.y-2, float64(b.width)+4, float64(b.height)+4, glowColor)
	}
}

// CheckCollision 检查子弹是否与玩家发生碰撞
func (b *EnemyBullet) CheckCollision(player *Player) bool {
	if !b.active {
		return false
	}

	// 简单的矩形碰撞检测
	return b.x < player.x+float64(player.width) &&
		b.x+float64(b.width) > player.x &&
		b.y < player.y+float64(player.height) &&
		b.y+float64(b.height) > player.y
}

// EnemyBulletManager 管理所有敌机子弹
type EnemyBulletManager struct {
	bullets []*EnemyBullet
}

// NewEnemyBulletManager 创建一个新的敌机子弹管理器
func NewEnemyBulletManager() *EnemyBulletManager {
	return &EnemyBulletManager{
		bullets: make([]*EnemyBullet, 0),
	}
}

// Update 更新所有敌机子弹的状态
func (bm *EnemyBulletManager) Update(enemies []*Enemy) {
	// 获取玩家实例（为追踪子弹使用）
	var player *Player
	if game != nil {
		player = game.player
	}

	// 更新现有子弹
	for i := len(bm.bullets) - 1; i >= 0; i-- {
		// 根据子弹类型调用不同的更新方法
		if bm.bullets[i].isHoming && player != nil {
			bm.bullets[i].UpdateHoming(player)
		} else {
			bm.bullets[i].Update()
		}

		// 移除非活动子弹
		if !bm.bullets[i].active {
			bm.bullets = append(bm.bullets[:i], bm.bullets[i+1:]...)
		}
	}

	// 随机让敌机发射子弹
	if enemies != nil {
		for _, enemy := range enemies {
			if enemy.active && rand.Float64() < 0.01 { // 1%的概率发射子弹
				// 从敌机的中心位置发射子弹
				bulletX := enemy.x + float64(enemy.width)/2 - 2
				bulletY := enemy.y + float64(enemy.height)
				bm.bullets = append(bm.bullets, NewEnemyBullet(bulletX, bulletY))
			}
		}
	}
}

// Draw 绘制所有敌机子弹
func (bm *EnemyBulletManager) Draw(screen *ebiten.Image) {
	for _, bullet := range bm.bullets {
		bullet.Draw(screen)
	}
}
