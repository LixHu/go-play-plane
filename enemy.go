package main

import (
	"image/color"
	//"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Enemy 表示敌机
type Enemy struct {
	x         float64
	y         float64
	speed     float64
	width     int
	height    int
	active    bool
	health    int // 当前血量
	maxHealth int // 最大血量
}

// NewEnemy 创建一个新的敌机
func NewEnemy() *Enemy {
	// 默认血量为2
	defaultHealth := 2
	return &Enemy{
		x:         float64(rand.Intn(screenWidth - 32)),
		y:         -32,
		speed:     2,
		width:     32,
		height:    32,
		active:    true,
		health:    defaultHealth, // 当前血量
		maxHealth: defaultHealth, // 最大血量与当前血量相同
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
	// 绘制敌机图像
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(e.x, e.y)
	screen.DrawImage(enemyImage, options)

	// 血条宽度与敌机相同
	bloodBarWidth := float64(e.width)
	const bloodBarHeight = 5.0

	// 先绘制整个血条的灰色背景（表示总血量）
	ebitenutil.DrawRect(screen, e.x, e.y-8, bloodBarWidth, bloodBarHeight, color.RGBA{100, 100, 100, 200})

	// 计算当前血量比例
	healthRatio := float64(e.health) / float64(e.maxHealth)
	healthWidth := bloodBarWidth * healthRatio

	// 根据健康比例变化颜色
	// 满血时是绿色(0,255,0)，血量越低越变红
	var healthColor color.RGBA
	if healthRatio == 1.0 {
		// 满血时显示绿色
		healthColor = color.RGBA{0, 255, 0, 255}
	} else {
		// 非满血时从黄色渐变到红色
		greenValue := uint8(200 * healthRatio)
		healthColor = color.RGBA{255, greenValue, 0, 255}
	}

	// 绘制健康部分
	ebitenutil.DrawRect(screen, e.x, e.y-8, healthWidth, bloodBarHeight, healthColor)
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

// Update 更新所有敌机的状态
func (em *EnemyManager) Update() {
	// 更新游戏时间
	em.gameTime++

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
		// 根据时间和关卡调整敌机血量
		baseHealth := 2                     // 基础血量
		levelBonus := em.level - 1          // 关卡加成
		timeBonus := int(em.gameTime / 600) // 时间加成，每10秒
		calculatedHealth := baseHealth + levelBonus + timeBonus
		// 限制最大血量，防止过高
		calculatedHealth = min(calculatedHealth, 20)

		// 设置当前血量和最大血量
		enemy.health = calculatedHealth
		enemy.maxHealth = calculatedHealth

		em.enemies = append(em.enemies, enemy)
		em.spawnTimer = 0
	}
}

// SetLevel 设置当前关卡并调整难度
func (em *EnemyManager) SetLevel(level int) {
	em.level = level
	em.difficulty = 1.0 + float64(level-1)*0.2 // 每关增加0.2的难度系数
	em.spawnInterval = max(20, 60-level*5)     // 每关减少5帧的生成间隔，最小20帧
	em.maxEnemies = min(30, 10+level*2)        // 每关增加2个最大敌机数量，最大30个
}

// Draw 绘制所有敌机
func (em *EnemyManager) Draw(screen *ebiten.Image) {
	for _, enemy := range em.enemies {
		enemy.Draw(screen)
	}
}
