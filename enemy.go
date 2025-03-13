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
	x      float64
	y      float64
	speed  float64
	width  int
	height int
	active bool
}

// NewEnemy 创建一个新的敌机
func NewEnemy() *Enemy {
	return &Enemy{
		x:      float64(rand.Intn(screenWidth - 32)),
		y:      -32,
		speed:  2,
		width:  32,
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
	enemies       []*Enemy
	spawnTimer    int
	spawnInterval int
	difficulty    float64 // 难度系数字段
	gameTime      int     // 游戏时间计数器（以帧为单位）
	maxEnemies    int     // 同时存在的最大敌机数量
	level         int     // 当前关卡
}

// NewEnemyManager 创建一个新的敌机管理器
func NewEnemyManager() *EnemyManager {
	rand.Seed(time.Now().UnixNano())
	return &EnemyManager{
		enemies:       make([]*Enemy, 0),
		spawnTimer:    0,
		spawnInterval: 60,  // 初始生成间隔为60帧
		difficulty:    1.0, // 初始难度系数
		gameTime:      0,   // 初始游戏时间
		maxEnemies:    10,  // 初始最大敌机数量
		level:         1,   // 初始关卡
	}
}

// SetLevel 设置当前关卡并调整难度
func (em *EnemyManager) SetLevel(level int) {
	em.level = level
	em.difficulty = 1.0 + float64(level-1)*0.2 // 每关增加0.2的难度系数
	em.spawnInterval = max(20, 60-level*5)     // 每关减少5帧的生成间隔，最小20帧
	em.maxEnemies = min(30, 10+level*2)        // 每关增加2个最大敌机数量，最大30个
}

// Update 更新所有敌机的状态
func (em *EnemyManager) Update() {
	// 更新游戏时间
	em.gameTime++

	// 每30秒（1800帧）增加一次难度
	if em.gameTime%1800 == 0 {
		// 增加难度系数
		em.difficulty += 0.1
		// 减少生成间隔（最小为20帧）
		if em.spawnInterval > 20 {
			em.spawnInterval -= 5
		}
		// 增加最大敌机数量（最大为30）
		if em.maxEnemies < 30 {
			em.maxEnemies += 2
		}
	}

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
	if em.spawnTimer >= em.spawnInterval && len(em.enemies) < em.maxEnemies {
		enemy := NewEnemy()
		// 根据难度调整敌机速度
		enemy.speed *= em.difficulty
		em.enemies = append(em.enemies, enemy)
		em.spawnTimer = 0
	}
}

// Draw 绘制所有敌机
func (em *EnemyManager) Draw(screen *ebiten.Image) {
	for _, enemy := range em.enemies {
		enemy.Draw(screen)
	}
}
