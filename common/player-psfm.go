//go:build poolsuite

package common

import(
	"github.com/mrusme/go-poolsuite"
)

type Player struct {
  Poolsuite                  *poolsuite.Poolsuite
  poolsuiteLoaded            bool
  poolsuitePlaying           bool
}

func NewPlayer() (*Player) {
  player := new(Player)

  player.Poolsuite = poolsuite.NewPoolsuite()
  player.poolsuiteLoaded = false
  player.poolsuitePlaying = false

  return player
}

func (p *Player) poolsuitePlay() {
  p.Poolsuite.Play(
    p.Poolsuite.GetRandomTrackFromPlaylist(
      p.Poolsuite.GetRandomPlaylist(),
    ),
    func() { p.poolsuitePlay() },
  )
}

func (p *Player) Play() {
  if p.poolsuiteLoaded == false {
    p.poolsuiteLoaded = true
    p.Poolsuite.Load()
  }

  if p.poolsuitePlaying == false {
    p.poolsuitePlay()
    p.poolsuitePlaying = true
  } else {
    p.Poolsuite.PauseResume()
    p.poolsuitePlaying = false
  }
}
