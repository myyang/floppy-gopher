package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

type pipes struct {
	mu      sync.RWMutex
	texture *sdl.Texture
	speed   int32

	pipes []*pipe
}

func newPipes(r *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(r, "img/pipe.png")
	if err != nil {
		return nil, err
	}

	ps := &pipes{
		texture: texture,
		speed:   3,
	}

	go func() {
		for {
			ps.mu.Lock()
			ps.pipes = append(ps.pipes, newPipe())
			ps.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return ps, nil
}

func (ps *pipes) paint(r *sdl.Renderer) error {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, v := range ps.pipes {
		if err := v.paint(r, ps.texture); err != nil {
			return err
		}
	}
	return nil
}

func (ps *pipes) restart() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.pipes = nil
}

func (ps *pipes) update() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	var rem []*pipe
	for _, v := range ps.pipes {
		v.update(ps.speed)
		if v.x+v.w > 0 {
			rem = append(rem, v)
		}
	}
	ps.pipes = rem
}

func (ps *pipes) destory() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.texture.Destroy()
}

func (ps *pipes) touch(b *bird) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	for _, v := range ps.pipes {
		b.touch(v)
	}
}

type pipe struct {
	mu       sync.RWMutex
	x, h, w  int32
	inverted bool
}

func newPipe() *pipe {

	return &pipe{
		x:        800,
		h:        int32(100 + rand.Intn(300)),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}
}

func (p *pipe) paint(r *sdl.Renderer, texture *sdl.Texture) error {
	p.mu.RLock()
	defer p.mu.RUnlock()

	pipeRect := &sdl.Rect{X: p.x, Y: 600 - p.h, W: p.w, H: p.h}
	flipFlag := sdl.FLIP_NONE
	if p.inverted {
		pipeRect.Y = 0
		flipFlag = sdl.FLIP_VERTICAL

	}
	if err := r.CopyEx(texture, nil, pipeRect, 0, nil, flipFlag); err != nil {
		return fmt.Errorf("Could not copy pipe to render. Error: %v", err)
	}
	return nil
}

func (p *pipe) restart() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x = 400

}

func (p *pipe) update(speed int32) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.x -= speed
}
