//go:build !poolsuite

package common

type Player struct {
}

func NewPlayer() (*Player) {
  player := new(Player)

  return player
}

func (p *Player) Play() {
  return
}

