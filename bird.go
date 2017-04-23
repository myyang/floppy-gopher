package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

const (
	gravity   = 0.1
	jumpSpeed = 5
)

type bird struct {
	mu sync.RWMutex

	time     int
	textures []*sdl.Texture

	w, h, x, y int32
	speed      float64
	dead       bool
}

func newBird(r *sdl.Renderer) (*bird, error) {

	var textures []*sdl.Texture
	for i := 1; i < 5; i++ {
		texture, err := img.LoadTexture(r, fmt.Sprintf("img/frame-%d.png", i))
		if err != nil {
			return nil, fmt.Errorf("Error to load bird. Error: %v", err)
		}
		textures = append(textures, texture)
	}
	return &bird{textures: textures, y: 300, h: 43, w: 50, x: 50}, nil
}

func (b *bird) update() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity
}

func (b *bird) paint(r *sdl.Renderer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	birdRect := &sdl.Rect{X: b.x, Y: 600 - b.y - b.h/2, W: b.w, H: b.h}
	if err := r.Copy(b.textures[b.time/10%len(b.textures)], nil, birdRect); err != nil {
		return fmt.Errorf("Could not copy bird to render. Error: %v", err)
	}
	return nil
}

func (b *bird) jump() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.speed = -jumpSpeed
}

func (b *bird) destory() {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, v := range b.textures {
		v.Destroy()
	}
}

func (b *bird) isDead() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	return b.dead
}

func (b *bird) restart() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) touch(p *pipe) {
	b.mu.Lock()
	defer b.mu.Unlock()
	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.x > b.x+b.w || p.x+p.w < b.x {
		return
	}
	if !p.inverted && p.h < b.y-b.h/2 {
		return
	} else if p.inverted && (600-p.h) > b.y+b.h/2 {
		return
	}
	b.dead = true
}
