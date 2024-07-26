package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"strconv"
	"strings"
)

// Mise à jour de l'état du jeu en fonction des entrées au clavier.
func (g *game) Update() error {
	g.stateFrame++
	g.stateFrame = g.stateFrame % globalBlinkDuration

	switch g.gameState {
	case waitState:
		var read string
		select {
		case read = <-g.reader:
			if read != "" {
				log.Println(read)
				g.gameState++
			}
		default:
		}
	case titleState:
		if g.titleUpdate() {
			g.gameState++
		}
	case colorSelectState:
		var read string
		select {
		case read = <-g.reader:
			read = strings.ReplaceAll(read, "\n", "")
			if read == "comeback" {
				g.colSelectedp2 = false
			}
		default:
		}
		if g.colorSelectUpdate(read) {
			g.gameState++
			g.writer <- "selected"
		}
	case colorOponentWait:
		var read string
		if g.colSelectedp2 {
			g.writer <- "play"
			g.gameState++
		} else {
			select {
			case read = <-g.reader:
				log.Println(read)
				read = strings.ReplaceAll(read, "\n", "")
				if read == "selected" {
					g.colSelectedp2 = true
				} else if col, err := strconv.Atoi(read); err == nil {
					g.p2Color = col
				}
			default:

				if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
					g.gameState--
					g.writer <- "comeback"
				}
			}
		}

	case turnState:
		var read string
		select {
		case read = <-g.reader:
			log.Println(read)
			read = strings.ReplaceAll(read, "\n", "")
			if read == "1" {
				g.turn = p2Turn
				g.gameState++
			} else if read == "0" {
				g.turn = p1Turn
				g.gameState++
			}
		default:
		}

	case playState:
		var read string
		select {
		case read = <-g.reader:
		default:
		}
		g.tokenPosUpdate(read)
		var lastXPositionPlayed int = 0
		var lastYPositionPlayed int = 0
		if g.turn == p1Turn {
			lastXPositionPlayed, lastYPositionPlayed = g.p1Update()
		} else {
			lastXPositionPlayed, lastYPositionPlayed = g.p2Update(read)
		}
		if lastXPositionPlayed >= 0 {
			finished, result := g.checkGameEnd(lastXPositionPlayed, lastYPositionPlayed)
			if finished {
				g.result = result
				g.gameState++
			}
		}
	case resultState:
		if g.resultUpdate() {
			g.gameState++
			g.writer <- "reset"
		}
	case resetState:
		var read string
		select {
		case read = <-g.reader:
		default:
		}
		read = strings.ReplaceAll(read, "\n", "")
		if read == "reset" {
			g.reset()
			g.gameState = playState
		}
	}

	return nil
}

// Mise à jour de l'état du jeu à l'écran titre.
func (g *game) titleUpdate() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de l'état du jeu lors de la sélection des couleurs.
func (g *game) colorSelectUpdate(colp2 string) bool {

	col := g.p1Color % globalNumColorCol
	line := g.p1Color / globalNumColorLine

	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		col = (col + 1) % globalNumColorCol
		g.p1Color = line*globalNumColorLine + col
		g.writer <- strconv.Itoa(g.p1Color)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		col = (col - 1 + globalNumColorCol) % globalNumColorCol
		g.p1Color = line*globalNumColorLine + col
		g.writer <- strconv.Itoa(g.p1Color)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		line = (line + 1) % globalNumColorLine
		g.p1Color = line*globalNumColorLine + col
		g.writer <- strconv.Itoa(g.p1Color)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		line = (line - 1 + globalNumColorLine) % globalNumColorLine
		g.p1Color = line*globalNumColorLine + col
		g.writer <- strconv.Itoa(g.p1Color)
	}

	if colp2 != "" && colp2 != "comeback" {
		colp2 = strings.ReplaceAll(colp2, "\n", "")
		if colp2 == "selected" {
			g.colSelectedp2 = true
		} else {
			g.p2Color, _ = strconv.Atoi(colp2)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if g.p1Color == g.p2Color {
			if !g.colSelectedp2 {
				g.colSelectedp1 = true
				return true
			} else {
				return false
			}
		} else {
			g.colSelectedp1 = true
			return true
		}
	}
	return false
}

func (g *game) colorOponent(read string) bool {
	if read != "" {
		read = strings.ReplaceAll(read, "\n", "")
		readConv, _ := strconv.Atoi(read)
		g.p2Color = readConv
		return true
	}
	return false
}

// Gestion de la position du prochain pion à jouer par le joueur 1.
func (g *game) tokenPosUpdate(pos string) {
	if g.turn == p1Turn {
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			g.tokenPosition = (g.tokenPosition - 1 + globalNumTilesX) % globalNumTilesX
			g.writer <- fmt.Sprint(g.tokenPosition)
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			g.tokenPosition = (g.tokenPosition + 1) % globalNumTilesX
			g.writer <- strconv.Itoa(g.tokenPosition)
		}
	}
	if g.turn == p2Turn && pos != "" && len(pos) == 2 {
		pos = strings.ReplaceAll(pos, "\n", "")
		var err error
		g.tokenPosition, err = strconv.Atoi(pos)
		if err != nil {
			log.Println("Error when converting : ", err)
		}
	}
}

// Gestion du moment où le prochain pion est joué par le joueur 1.
func (g *game) p1Update() (int, int) {
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if updated, yPos := g.updateGrid(p1Token, g.tokenPosition); updated {
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
			g.writer <- "update"
			g.turn = p2Turn
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Gestion de la position du prochain pion joué par le joueur 2 et
// du moment où ce pion est joué.
func (g *game) p2Update(pass string) (int, int) {
	/*position := rand.Intn(globalNumTilesX)
	updated, yPos := g.updateGrid(p2Token, position)
	for ; !updated; updated, yPos = g.updateGrid(p2Token, position) {
		position = (position + 1) % globalNumTilesX
	}*/
	lastXPositionPlayed := -1
	lastYPositionPlayed := -1
	pass = strings.ReplaceAll(pass, "\n", "")
	if pass == "update" {
		if updated, yPos := g.updateGrid(p2Token, g.tokenPosition); updated {
			lastXPositionPlayed = g.tokenPosition
			lastYPositionPlayed = yPos
			g.turn = p1Turn
		}
	}
	return lastXPositionPlayed, lastYPositionPlayed
}

// Mise à jour de l'état du jeu à l'écran des résultats.
func (g game) resultUpdate() bool {
	return inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

// Mise à jour de la grille de jeu lorsqu'un pion est inséré dans la
// colonne de coordonnée (x) position.
func (g *game) updateGrid(token, position int) (updated bool, yPos int) {
	for y := globalNumTilesY - 1; y >= 0; y-- {
		if g.grid[position][y] == noToken {
			updated = true
			yPos = y
			g.grid[position][y] = token
			return
		}
	}
	return
}

// Vérification de la fin du jeu : est-ce que le dernier joueur qui
// a placé un pion gagne ? est-ce que la grille est remplie sans gagnant
// (égalité) ? ou est-ce que le jeu doit continuer ?
func (g game) checkGameEnd(xPos, yPos int) (finished bool, result int) {

	tokenType := g.grid[xPos][yPos]

	// horizontal
	count := 0
	for x := xPos; x < globalNumTilesX && g.grid[x][yPos] == tokenType; x++ {
		count++
	}
	for x := xPos - 1; x >= 0 && g.grid[x][yPos] == tokenType; x-- {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// vertical
	count = 0
	for y := yPos; y < globalNumTilesY && g.grid[xPos][y] == tokenType; y++ {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut gauche/bas droit
	count = 0
	for x, y := xPos, yPos; x < globalNumTilesX && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x+1, y+1 {
		count++
	}

	for x, y := xPos-1, yPos-1; x >= 0 && y >= 0 && g.grid[x][y] == tokenType; x, y = x-1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// diag haut droit/bas gauche
	count = 0
	for x, y := xPos, yPos; x >= 0 && y < globalNumTilesY && g.grid[x][y] == tokenType; x, y = x-1, y+1 {
		count++
	}

	for x, y := xPos+1, yPos-1; x < globalNumTilesX && y >= 0 && g.grid[x][y] == tokenType; x, y = x+1, y-1 {
		count++
	}

	if count >= 4 {
		if tokenType == p1Token {
			return true, p1wins
		}
		return true, p2wins
	}

	// egalité ?
	if yPos == 0 {
		for x := 0; x < globalNumTilesX; x++ {
			if g.grid[x][0] == noToken {
				return
			}
		}
		return true, equality
	}

	return
}
