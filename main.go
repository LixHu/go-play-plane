package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	gameTitle    = "打飞机游戏"
)

// Game 结构体用于保存游戏状态
type Game struct {
	player *Player
	enemyManager *EnemyManager
	bulletManager *BulletManager
	enemyBulletManager *EnemyBulletManager
	score int
	isGameOver bool
}

// Update 处理游戏逻辑更新
func (g *Game) Update() error {
	// 如果游戏已结束，只处理重新开始的输入
	if g.isGameOver {
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			// 重置游戏状态
			g.player = NewPlayer()
			g.enemyManager = NewEnemyManager()
			g.bulletManager = NewBulletManager()
			g.enemyBulletManager = NewEnemyBulletManager()
			g.score = 0
			g.isGameOver = false
		}
		return nil
	}

	// 更新玩家状态
	g.player.Update()

	// 更新敌机状态
	g.enemyManager.Update()

	// 更新子弹状态
	g.bulletManager.Update(g.player)

	// 检测子弹与敌机的碰撞
	for _, bullet := range g.bulletManager.bullets {
		for _, enemy := range g.enemyManager.enemies {
			if bullet.CheckCollision(enemy) {
				bullet.active = false
				enemy.active = false
				g.score += 100
			}
		}
	}

	// 检测玩家与敌机的碰撞
	for _, enemy := range g.enemyManager.enemies {
		if enemy.active && g.checkPlayerCollision(enemy) {
			g.isGameOver = true
			break
		}
	}

	// 更新敌机子弹状态
	g.enemyBulletManager.Update(g.enemyManager.enemies)

	// 检测敌机子弹与玩家的碰撞
	for _, bullet := range g.enemyBulletManager.bullets {
		if bullet.CheckCollision(g.player) {
			g.isGameOver = true
			break
		}
	}

	return nil
}

// checkPlayerCollision 检测玩家与敌机的碰撞
func (g *Game) checkPlayerCollision(enemy *Enemy) bool {
	return g.player.x < enemy.x+float64(enemy.width) &&
		g.player.x+float64(g.player.width) > enemy.x &&
		g.player.y < enemy.y+float64(enemy.height) &&
		g.player.y+float64(g.player.height) > enemy.y
}

// Draw 处理游戏画面渲染
func (g *Game) Draw(screen *ebiten.Image) {
	// 绘制玩家
	g.player.Draw(screen)

	// 绘制敌机
	g.enemyManager.Draw(screen)

	// 绘制玩家子弹
	g.bulletManager.Draw(screen)

	// 绘制敌机子弹
	g.enemyBulletManager.Draw(screen)

	// 绘制分数
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d", g.score))

	// 如果游戏结束，显示游戏结束信息
	if g.isGameOver {
		msg := fmt.Sprintf("游戏结束！最终得分：%d\n按空格键重新开始", g.score)
		ebitenutil.DebugPrintAt(screen, msg, screenWidth/2-150, screenHeight/2)
	}
}

// Layout 返回游戏窗口的大小
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(gameTitle)

	game := &Game{
		player: NewPlayer(),
		enemyManager: NewEnemyManager(),
		bulletManager: NewBulletManager(),
		enemyBulletManager: NewEnemyBulletManager(),
		score: 0,
		isGameOver: false,
	}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}