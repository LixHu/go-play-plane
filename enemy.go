package main

import (
	"image/color"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Enemy 表示敌机
type Enemy struct {
	x float64
	y float64
	speed float64
	width int
	height int
	active bool
}

// NewEnemy 创建一个新的敌机
func NewEnemy() *Enemy {
	return &Enemy{
		x: float64(rand.Intn(screenWidth - 32)),
		y: -32,
		speed: 2,
		width: 32,
		height: 32,
		active: true,
	}
}

// Update 更新敌机的状态
func (e *Enemy) Update() {
	// 向下移动
	e.y += e.speed

	// 如果飞出屏幕外，标记为非活动状态
	if e.y > float64(screenHeight) {
		e.active = false
	}
}

// Draw 绘制敌机
func (e *Enemy) Draw(screen *ebiten.Image) {
	// 临时使用一个简单的红色矩形代表敌机
	ebitenutil.DrawRect(screen, e.x, e.y, float64(e.width), float64(e.height), color.RGBA{255, 0, 0, 255})
}

// EnemyManager 管理所有敌机
type EnemyManager struct {
	enemies []*Enemy
	spawnTimer int
	spawnInterval int
}

// NewEnemyManager 创建一个新的敌机管理器
func NewEnemyManager() *EnemyManager {
	rand.Seed(time.Now().UnixNano())
	return &EnemyManager{
		enemies: make([]*Enemy, 0),
		spawnTimer: 0,
		spawnInterval: 60, // 每60帧生成一个新敌机
	}
}

// Update 更新所有敌机的状态
func (em *EnemyManager) Update() {
	// 更新现有敌机
	for i := len(em.enemies) - 1; i >= 0; i-- {
		em.enemies[i].Update()
		// 移除非活动敌机
		if !em.enemies[i].active {
			em.enemies = append(em.enemies[:i], em.enemies[i+1:]...)
		}
	}

	// 生成新敌机
	em.spawnTimer++
	if em.spawnTimer >= em.spawnInterval {
		em.enemies = append(em.enemies, NewEnemy())
		em.spawnTimer = 0
	}
}

// Draw 绘制所有敌机
func (em *EnemyManager) Draw(screen *ebiten.Image) {
	for _, enemy := range em.enemies {
		enemy.Draw(screen)
	}
}