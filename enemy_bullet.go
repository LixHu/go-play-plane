package main

import (
	"image/color"
	"math"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math/rand"
)

// EnemyBullet 表示敌机发射的子弹
type EnemyBullet struct {
	x float64
	y float64
	speedX float64
	speedY float64
	width int
	height int
	active bool
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
		x: x,
		y: y,
		speedX: baseSpeed * math.Sin(radian),
		speedY: baseSpeed * math.Cos(radian),
		width: 4,
		height: 4, // 修改为正方形，便于旋转
		active: true,
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

// Draw 绘制敌机子弹
func (b *EnemyBullet) Draw(screen *ebiten.Image) {
	// 使用红色矩形代表敌机子弹
	ebitenutil.DrawRect(screen, b.x, b.y, float64(b.width), float64(b.height), color.RGBA{255, 0, 0, 255})
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
	// 更新现有子弹
	for i := len(bm.bullets) - 1; i >= 0; i-- {
		bm.bullets[i].Update()
		// 移除非活动子弹
		if !bm.bullets[i].active {
			bm.bullets = append(bm.bullets[:i], bm.bullets[i+1:]...)
		}
	}

	// 随机让敌机发射子弹
	for _, enemy := range enemies {
		if enemy.active && rand.Float64() < 0.01 { // 1%的概率发射子弹
			// 从敌机的中心位置发射子弹
			bulletX := enemy.x + float64(enemy.width)/2 - 2
			bulletY := enemy.y + float64(enemy.height)
			bm.bullets = append(bm.bullets, NewEnemyBullet(bulletX, bulletY))
		}
	}
}

// Draw 绘制所有敌机子弹
func (bm *EnemyBulletManager) Draw(screen *ebiten.Image) {
	for _, bullet := range bm.bullets {
		bullet.Draw(screen)
	}
}