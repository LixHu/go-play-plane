package main

import (
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Bullet 表示玩家发射的子弹
type Bullet struct {
	x float64
	y float64
	speed float64
	width int
	height int
	active bool
}

// NewBullet 创建一个新的子弹
func NewBullet(x, y float64) *Bullet {
	return &Bullet{
		x: x,
		y: y,
		speed: 8,
		width: 4,
		height: 10,
		active: true,
	}
}

// Update 更新子弹的状态
func (b *Bullet) Update() {
	// 向上移动
	b.y -= b.speed

	// 如果飞出屏幕外，标记为非活动状态
	if b.y < -float64(b.height) {
		b.active = false
	}
}

// Draw 绘制子弹
func (b *Bullet) Draw(screen *ebiten.Image) {
	// 使用黄色矩形代表子弹
	ebitenutil.DrawRect(screen, b.x, b.y, float64(b.width), float64(b.height), color.RGBA{255, 255, 0, 255})
}

// CheckCollision 检查子弹是否与敌机发生碰撞
func (b *Bullet) CheckCollision(enemy *Enemy) bool {
	if !b.active || !enemy.active {
		return false
	}

	// 简单的矩形碰撞检测
	return b.x < enemy.x+float64(enemy.width) &&
		b.x+float64(b.width) > enemy.x &&
		b.y < enemy.y+float64(enemy.height) &&
		b.y+float64(b.height) > enemy.y
}

// BulletManager 管理所有子弹
type BulletManager struct {
	bullets []*Bullet
	shootTimer int
	shootInterval int
}

// NewBulletManager 创建一个新的子弹管理器
func NewBulletManager() *BulletManager {
	return &BulletManager{
		bullets: make([]*Bullet, 0),
		shootTimer: 0,
		shootInterval: 10, // 每10帧可以发射一颗子弹
	}
}

// Update 更新所有子弹的状态
func (bm *BulletManager) Update(player *Player) {
	// 更新现有子弹
	for i := len(bm.bullets) - 1; i >= 0; i-- {
		bm.bullets[i].Update()
		// 移除非活动子弹
		if !bm.bullets[i].active {
			bm.bullets = append(bm.bullets[:i], bm.bullets[i+1:]...)
		}
	}

	// 发射新子弹
	bm.shootTimer++
	if ebiten.IsKeyPressed(ebiten.KeySpace) && bm.shootTimer >= bm.shootInterval {
		// 从玩家飞机的中心位置发射子弹
		bulletX := player.x + float64(player.width)/2 - 2
		bulletY := player.y
		bm.bullets = append(bm.bullets, NewBullet(bulletX, bulletY))
		bm.shootTimer = 0
	}
}

// Draw 绘制所有子弹
func (bm *BulletManager) Draw(screen *ebiten.Image) {
	for _, bullet := range bm.bullets {
		bullet.Draw(screen)
	}
}