// Copyright 2015 Pikkpoiss
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"time"

	"../lib/twodee"
	"github.com/go-gl/mathgl/mgl32"
)

type Mob interface {
	Bored(time.Duration) bool
	SetFrames(f []int)
	Detect(dist float32) bool
	Pos() twodee.Point
	Bounds() twodee.Rectangle
	SearchPattern() []mgl32.Vec2
	Speed() float32
	MoveTo(twodee.Point)
}

type Mobile struct {
	DetectionRadius float32
	BoredThreshold  time.Duration
	speed           float32
	searchPattern   []mgl32.Vec2
}

func (m *Mobile) Bored(d time.Duration) bool {
	return d >= m.BoredThreshold
}

func (m *Mobile) Detect(d float32) bool {
	return d <= m.DetectionRadius
}

func (m *Mobile) SearchPattern() []mgl32.Vec2 {
	return m.searchPattern
}

func (m *Mobile) Speed() float32 {
	return m.speed
}

// MoveMob moves the mob along the given vector, which should be normalized.
func MoveMob(m Mob, v mgl32.Vec2, l *Level) {
	bounds := m.Bounds()
	pos := m.Pos()
	v = l.Collisions.FixMove(mgl32.Vec4{
		bounds.Min.X,
		bounds.Min.Y,
		bounds.Max.X,
		bounds.Max.Y,
	}, v, 0.5, 0.5)
	m.MoveTo(twodee.Pt(pos.X+v[0], pos.Y+v[1]))
}

// MobState is implemented by various states responsible for controlling mobile
// entities.
// ExamineWorld examines the current state of the game and determines which
// state
type MobState interface {
	// ExamineWorld examines the current state of the game and determines
	// which state the mobile should transition to. It's legal to return
	// either the current MobState or nil, indicating that the mobile
	// should transition to a previous state.
	ExamineWorld(Mob, *Level) (newState MobState)
	// Update should be called each frame and may update values in the
	// current state or call functions on the mob.
	Update(Mob, time.Duration, *Level)
	// Enter should be called when entering this MobState.
	Enter(Mob)
	// Enter should be called when exiting this MobState.
	Exit(Mob)
}

// VegState encapsulates the state of being a vegetable.
type VegState struct{}

// ExamineWorld always returns a new SearchState.
func (v *VegState) ExamineWorld(m Mob, l *Level) MobState {
	return &SearchState{m.SearchPattern(), 0}
}

func (v *VegState) Update(m Mob, d time.Duration, l *Level) {}

func (v *VegState) Enter(m Mob) {}

func (v *VegState) Exit(m Mob) {}

// SearchState is the state during which a mobile is aimlessly wandering,
// hoping to chance across the player.
type SearchState struct {
	Pattern        []mgl32.Vec2
	targetPointIdx int
}

// ExamineWorld returns HuntState if the player is seen, otherwise the mob
// continues wandering.
func (s *SearchState) ExamineWorld(m Mob, l *Level) MobState {
	if playerSeen(m, l) {
		return &HuntState{}
	}
	return s
}

// TODO: Not entirely sure that movement updates should be done here instead
// of in ExamineWorld...
func (s *SearchState) Update(m Mob, d time.Duration, l *Level) {
	if len(s.Pattern) == 0 {
		// Do nothing right now with no search pattern.
		return
	}
	tv := s.Pattern[s.targetPointIdx]
	mv := mgl32.Vec2{m.Pos().X, m.Pos().Y}
	if tv.Sub(mv).Len() < 2 { // CLOSE ENOUGH!
		s.targetPointIdx = (s.targetPointIdx + 1) % len(s.Pattern)
		tv = s.Pattern[s.targetPointIdx]
	}
	MoveMob(m, tv.Normalize().Mul(m.Speed()), l)
}

func (s *SearchState) Enter(m Mob) {
	fmt.Println("In the search state!")
	// TODO: set some hunting animation.
}

func (s *SearchState) Exit(m Mob) {
	fmt.Println("Leaving the search state!")
	// TODO: maybe something should happen when we start hunting?
}

// HuntState is the state during which a mobile is actively hunting the player.
type HuntState struct {
	durSinceLastContact time.Duration
}

// ExamineWorld returns the current state if the player is currently seen or
// the mob is not yet tired of chasing. Otherwise, it returns nil.
func (h *HuntState) ExamineWorld(m Mob, l *Level) MobState {
	if playerSeen(m, l) {
		h.durSinceLastContact = time.Duration(0)
		return h
	}
	if !m.Bored(h.durSinceLastContact) {
		return h
	}
	return nil
}

// Update resets the player's hiding timer if the player is seen, otherwise it
// increments.
func (h *HuntState) Update(m Mob, d time.Duration, l *Level) {
	h.durSinceLastContact += d
}

func (h *HuntState) Enter(m Mob) {
	fmt.Println("Entering the hunt state")
}

func (h *HuntState) Exit(m Mob) {
	fmt.Println("Exiting the hunt state")
}

// playerSeen returns true if the player is currently visible to the mob and
// within its detection radius.
func playerSeen(m Mob, l *Level) bool {
	c := l.Collisions
	mpv := mgl32.Vec2{m.Pos().X, m.Pos().Y}
	ppv := mgl32.Vec2{l.Player.Pos().X, l.Player.Pos().Y}
	if c.CanSee(mpv, ppv, 0.5, 0.5) && m.Detect(mpv.Sub(ppv).Len()) {
		return true
	}
	return false
}
