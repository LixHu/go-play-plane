package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// Player 表示玩家控制的飞机
type Player struct {
	x                 float64
	y                 float64
	speed             float64
	width             int
	height            int
	multiShotCount    int // 永久性多弹道数量
	screenShotEnabled bool
	powerUpTimer      int // 用于控制全屏攻击的持续时间
	attackPower       int // 攻击力
}

// NewPlayer 创建一个新的玩家飞机
func NewPlayer() *Player {
	return &Player{
		x:              float64(screenWidth) / 2,
		y:              float64(screenHeight) - 50,
		speed:          4,
		width:          32,
		height:         32,
		multiShotCount: 0,
		attackPower:    1,
	}
}

// EnableMultiShot 启用多弹道能力
func (p *Player) EnableMultiShot() {
	p.multiShotCount++ // 永久增加一个弹道
}

// EnableScreenShot 启用全屏攻击能力
func (p *Player) EnableScreenShot() {
	p.screenShotEnabled = true
	p.powerUpTimer = 180 // 能力持续180帧（约3秒）
}

// EnableAttackBoost 增加攻击力
func (p *Player) EnableAttackBoost() {
	p.attackPower++ // 永久增加一点攻击力
}

// Update 更新玩家飞机的状态
func (p *Player) Update() {
	// 处理键盘输入
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && p.x > 0 {
		p.x -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && p.x < float64(screenWidth-p.width) {
		p.x += p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) && p.y > 0 {
		p.y -= p.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && p.y < float64(screenHeight-p.height) {
		p.y += p.speed
	}

	// 更新全屏攻击状态
	if p.screenShotEnabled {
		p.powerUpTimer--
		if p.powerUpTimer <= 0 {
			p.screenShotEnabled = false
		}
	}
}

// Draw 绘制玩家飞机
func (p *Player) Draw(screen *ebiten.Image) {
	// 临时使用一个简单的蓝色矩形代表玩家飞机
	ebitenutil.DrawRect(screen, p.x, p.y, float64(p.width), float64(p.height), color.RGBA{0, 0, 255, 255})
}
