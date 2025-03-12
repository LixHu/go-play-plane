package main

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// 全局变量，用于跟踪玩家得分
var score int

// PowerUpType 定义道具类型
type PowerUpType int

const (
	MultiShot PowerUpType = iota // 多弹道
	ScreenShot                    // 全屏攻击
)

// PowerUp 表示道具
type PowerUp struct {
	x      float64
	y      float64
	speed  float64
	width  int
	height int
	active bool
	pType  PowerUpType
}

// NewPowerUp 创建一个新的道具
func NewPowerUp(x, y float64, pType PowerUpType) *PowerUp {
	return &PowerUp{
		x:      x,
		y:      y,
		speed:  1.5,
		width:  20,
		height: 20,
		active: true,
		pType:  pType,
	}
}

// Update 更新道具的状态
func (p *PowerUp) Update() {
	// 道具向下移动
	p.y += p.speed

	// 如果飞出屏幕外，标记为非活动状态
	if p.y > float64(screenHeight) {
		p.active = false
	}
}

// Draw 绘制道具
func (p *PowerUp) Draw(screen *ebiten.Image) {
	// 根据道具类型使用不同的颜色
	var clr color.Color
	switch p.pType {
	case MultiShot:
		clr = color.RGBA{0, 255, 0, 255} // 绿色表示多弹道
	case ScreenShot:
		clr = color.RGBA{0, 0, 255, 255} // 蓝色表示全屏攻击
	}

	// 绘制道具
	ebitenutil.DrawRect(screen, p.x, p.y, float64(p.width), float64(p.height), clr)
}

// PowerUpManager 管理所有道具
type PowerUpManager struct {
	powerUps []*PowerUp
}

// NewPowerUpManager 创建一个新的道具管理器
func NewPowerUpManager() *PowerUpManager {
	return &PowerUpManager{
		powerUps: make([]*PowerUp, 0),
	}
}

// SpawnPowerUp 生成道具
func (pm *PowerUpManager) SpawnPowerUp(x, y float64) {
	// 基础掉落概率为35%
	baseProb := 0.35
	// 根据玩家得分增加掉落概率，每1000分增加5%的掉落概率，最高不超过60%
	scoreBonus := math.Min(float64(score)/1000.0*0.05, 0.25)
	if rand.Float64() < baseProb+scoreBonus {
		// 随机选择道具类型，降低全屏攻击道具的掉落概率
		if rand.Float64() < 0.3 { // 30%概率掉落全屏攻击道具
			pm.powerUps = append(pm.powerUps, NewPowerUp(x, y, ScreenShot))
		} else {
			pm.powerUps = append(pm.powerUps, NewPowerUp(x, y, MultiShot))
		}
	}
}

// Update 更新所有道具的状态
func (pm *PowerUpManager) Update() {
	// 更新现有道具
	for i := len(pm.powerUps) - 1; i >= 0; i-- {
		pm.powerUps[i].Update()
		// 移除非活动道具
		if !pm.powerUps[i].active {
			pm.powerUps = append(pm.powerUps[:i], pm.powerUps[i+1:]...)
		}
	}
}

// Draw 绘制所有道具
func (pm *PowerUpManager) Draw(screen *ebiten.Image) {
	for _, powerUp := range pm.powerUps {
		powerUp.Draw(screen)
	}
}