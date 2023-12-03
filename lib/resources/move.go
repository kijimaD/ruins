package resources

import w "github.com/kijimaD/sokotwo/lib/engine/world"

type MovementType uint8

const (
	MovementUp MovementType = iota
	MovementDown
	MovementLeft
	MovementRight
)

func Move(world w.World, movements ...MovementType) {
	gameResources := world.Resources.Game.(*Game)

	levelWidth := gameResources.Level.Grid.NCols
	levelHeight := gameResources.Level.Grid.NRows

	playerIndex := -1
	for iTile, tile := range gameResources.Level.Grid.Data {
		if tile.Contains(TilePlayer) {
			playerIndex = iTile
			break
		}
	}

	playerTile := &gameResources.Level.Grid.Data[playerIndex]
	playerLine := playerIndex / levelWidth
	playerCol := playerIndex % levelWidth

	for _, movement := range movements {
		var directionLine, directionCol int
		switch movement {
		case MovementUp:
			directionLine, directionCol = -1, 0
		case MovementDown:
			directionLine, directionCol = 1, 0
		case MovementLeft:
			directionLine, directionCol = 0, -1
		case MovementRight:
			directionLine, directionCol = 0, 1
		}

		oneFrontLine := playerLine + directionLine
		oneFrontCol := playerCol + directionCol

		// Check grid edge
		if !(0 <= oneFrontLine && oneFrontLine < levelHeight && 0 <= oneFrontCol && oneFrontCol < levelWidth) {
			return
		}
		oneFrontTile := gameResources.Level.Grid.Get(oneFrontLine, oneFrontCol)

		// No move if a wall is ahead
		if oneFrontTile.Contains(TileWall) {
			return
		}

		// 次の階層へ
		if oneFrontTile.Contains(TileWarpNext) {
			newLevel := gameResources.Level.CurrentNum + 1
			world.Manager.DeleteAllEntities()
			InitLevel(world, newLevel)
			return
		}

		oneFrontTile.Set(TilePlayer)
		playerTile.Remove(TilePlayer)

		playerTile = oneFrontTile
		playerLine += directionLine
		playerCol += directionCol

		gameResources.Level.Movements = append(gameResources.Level.Movements, movement)
		gameResources.Level.Modified = true
	}
}
