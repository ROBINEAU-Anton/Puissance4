package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
)

// Affichage des graphismes à l'écran selon l'état actuel du jeu.
func (g *game) Draw(screen *ebiten.Image) {

	screen.Fill(globalBackgroundColor)

	switch g.gameState {
	/*case connState:
	g.connDraw(screen)*/
	case waitState:
		g.waitDraw(screen)
	case titleState:
		g.titleDraw(screen)
	case colorSelectState:
		g.colorSelectDraw(screen)
	case colorOponentWait:
		g.colorSelectWait(screen)
	case turnState:
		g.turnAttribution(screen)
	case playState:
		g.playDraw(screen)
	case resultState:
		g.resultDraw(screen)
	case resetState:
		g.resetDraw(screen)
	}

}

// Affichage des graphismes d'attente d'un deuxième client
func (g game) waitDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", smallFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", smallFont, 210, 230, globalTextColor)
	if g.stateFrame%60 >= 0 && g.stateFrame%60 < 20 {
		text.Draw(screen, "En attente d'un adversaire.", smallFont, 180, 500, globalTextColor)
	} else if g.stateFrame%60 >= 20 && g.stateFrame%60 < 40 {
		text.Draw(screen, "En attente d'un adversaire..", smallFont, 180, 500, globalTextColor)
	} else {
		text.Draw(screen, "En attente d'un adversaire...", smallFont, 180, 500, globalTextColor)
	}
}

// Affichage des graphismes de l'écran titre.
func (g game) titleDraw(screen *ebiten.Image) {
	text.Draw(screen, "Puissance 4 en réseau", largeFont, 90, 150, globalTextColor)
	text.Draw(screen, "Projet de programmation système", smallFont, 105, 190, globalTextColor)
	text.Draw(screen, "Année 2023-2024", smallFont, 210, 230, globalTextColor)

	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "Appuyez sur entrée", smallFont, 210, 500, globalTextColor)
	}
}

// Affichage des graphismes de l'écran de sélection des couleurs des joueurs.
func (g game) colorSelectDraw(screen *ebiten.Image) {
	text.Draw(screen, "Quelle couleur pour vos pions ?", smallFont, 110, 80, globalTextColor)

	line := 0
	col := 0
	for numColor := 0; numColor < globalNumColor; numColor++ {

		xPos := (globalNumTilesX-globalNumColorCol)/2 + col
		yPos := (globalNumTilesY-globalNumColorLine)/2 + line

		if numColor == g.p1Color {
			vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, color.NRGBA{0, 0, 255, 255}, true)
			if g.p2Color == numColor {
				vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, color.NRGBA{127, 0, 255, 255}, true)
			}
		} else if numColor == g.p2Color {
			if g.colSelectedp2 {
				vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, color.NRGBA{0, 255, 0, 255}, true)
			} else {
				vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2, color.NRGBA{255, 0, 0, 255}, true)
			}
		}

		vector.DrawFilledCircle(screen, float32(globalTileSize/2+xPos*globalTileSize), float32(globalTileSize+globalTileSize/2+yPos*globalTileSize), globalTileSize/2-globalCircleMargin, globalTokenColors[numColor], true)

		col++
		if col >= globalNumColorCol {
			col = 0
			line++
		}
	}
}

// Affichage des graphismes de l'écran d'attente de la selection de la couleur de l'adversaire.
func (g game) colorSelectWait(screen *ebiten.Image) {
	text.Draw(screen, "La partie va bientôt commencer", smallFont, 110, 80, globalTextColor)
	vector.DrawFilledCircle(screen, float32(globalWidth/2-200), float32(globalHeight/2-20), globalTileSize, globalTokenColors[g.p1Color], true)
	text.Draw(screen, "VS", largeFont, globalWidth/2-30, globalHeight/2, globalTextColor)
	vector.DrawFilledCircle(screen, float32(globalWidth/2+200), float32(globalHeight/2-20), globalTileSize, globalTokenColors[g.p2Color], true)
	if g.stateFrame >= globalBlinkDuration/3 {
		text.Draw(screen, "En attente de l'adversaire", smallFont, 190, 600, globalTextColor)
	}
	text.Draw(screen, "Appuyez sur Echap pour changer de couleur", smallFont, 50, 650, globalTextColor)
}

// Affichage des graphismes lors de l'attribution du tour.
func (g game) turnAttribution(screen *ebiten.Image) {
	text.Draw(screen, "La partie va bientôt commencer", smallFont, 110, 80, globalTextColor)
	vector.DrawFilledCircle(screen, float32(globalWidth/2-200), float32(globalHeight/2-20), globalTileSize, globalTokenColors[g.p1Color], true)
	text.Draw(screen, "VS", largeFont, globalWidth/2-30, globalHeight/2, globalTextColor)
	vector.DrawFilledCircle(screen, float32(globalWidth/2+200), float32(globalHeight/2-20), globalTileSize, globalTokenColors[g.p2Color], true)
	text.Draw(screen, "La partie va commencer", smallFont, 190, 600, globalTextColor)
}

// Affichage des graphismes durant le jeu.
func (g game) playDraw(screen *ebiten.Image) {
	g.drawGrid(screen)

	if g.turn == p1Turn {
		vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p1Color], true)
		text.Draw(screen, "J1", smallFont, (globalTileSize/2+g.tokenPosition*globalTileSize)-15, (globalTileSize/2)+10, color.NRGBA{0, 0, 0, 255})
	}
	if g.turn == p2Turn {
		vector.DrawFilledCircle(screen, float32(globalTileSize/2+g.tokenPosition*globalTileSize), float32(globalTileSize/2), globalTileSize/2-globalCircleMargin, globalTokenColors[g.p2Color], true)
		text.Draw(screen, "J2", smallFont, (globalTileSize/2+g.tokenPosition*globalTileSize)-15, (globalTileSize/2)+10, color.NRGBA{0, 0, 0, 255})
	}
}

// Affichage des graphismes à l'écran des résultats.
func (g game) resultDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	message := "Égalité"
	if g.result == p1wins {
		message = "Gagné !"
	} else if g.result == p2wins {
		message = "Perdu…"
	}
	text.Draw(screen, message, smallFont, 300, 400, globalTextColor)
	text.Draw(screen, "Appuyez sur entrer pour recommencer", smallFont, 60, 450, globalTextColor)
}

// Affichage des graphismes en attente du redémarrage de la partie par l'adversaire.
func (g game) resetDraw(screen *ebiten.Image) {
	g.drawGrid(offScreenImage)

	options := &ebiten.DrawImageOptions{}
	options.ColorScale.ScaleAlpha(0.2)
	screen.DrawImage(offScreenImage, options)

	text.Draw(screen, "En attente de l'adversaire", smallFont, 150, 400, globalTextColor)
}

// Affichage de la grille de puissance 4, incluant les pions déjà joués.
func (g game) drawGrid(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 0, globalTileSize, globalTileSize*globalNumTilesX, globalTileSize*globalNumTilesY, globalGridColor, true)

	for x := 0; x < globalNumTilesX; x++ {
		for y := 0; y < globalNumTilesY; y++ {

			var tileColor color.Color
			switch g.grid[x][y] {
			case p1Token:
				tileColor = globalTokenColors[g.p1Color]
			case p2Token:
				tileColor = globalTokenColors[g.p2Color]
			default:
				tileColor = globalBackgroundColor
			}

			vector.DrawFilledCircle(screen, float32(globalTileSize/2+x*globalTileSize), float32(globalTileSize+globalTileSize/2+y*globalTileSize), globalTileSize/2-globalCircleMargin, tileColor, true)
		}
	}
}
