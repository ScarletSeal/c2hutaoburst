package travelerdendro

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames [][]int

const skillHitmark = 28
const cdStart = 25

func init() {
	skillFrames = make([][]int, 2)

	// Male
	skillFrames[0] = frames.InitAbilSlice(37) // E -> N1
	skillFrames[0][action.ActionDash] = 29    // E -> D
	skillFrames[0][action.ActionJump] = 29    // E -> J
	skillFrames[0][action.ActionSwap] = 36    // E -> Swap

	// Female
	skillFrames[1] = frames.InitAbilSlice(37) // E -> N1/Q
	skillFrames[1][action.ActionDash] = 28    // E -> D
	skillFrames[1][action.ActionJump] = 28    // E -> J
	skillFrames[1][action.ActionSwap] = 35    // E -> Swap
}

func (c *char) Skill(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex: c.Index,
		Abil:       "Razorgrass Blade",
		AttackTag:  combat.AttackTagElementalArt,
		ICDTag:     combat.ICDTagNone,
		ICDGroup:   combat.ICDGroupDefault,
		StrikeType: combat.StrikeTypeDefault,
		Element:    attributes.Dendro,
		Durability: 25,
		Mult:       skill[c.TalentLvlSkill()],
	}

	var count float64 = 2
	if c.Core.Rand.Float64() < 0.50 {
		count = 3
	}

	var skillCB func(a combat.AttackCB)
	if c.Base.Cons >= 1 {
		c.skillC1 = true
		skillCB = c.c1cb()
	}
	c.Core.QueueAttack(
		ai,
		combat.NewCircleHitOnTargetFanAngle(c.Core.Combat.Player(), combat.Point{Y: -0.3}, 6.5, 130),
		skillHitmark,
		skillHitmark,
		skillCB,
	)

	c.Core.QueueParticle(c.Base.Key.String(), count, attributes.Dendro, skillHitmark+c.ParticleDelay)
	c.SetCDWithDelay(action.ActionSkill, 8*60, cdStart)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames[c.gender]),
		AnimationLength: skillFrames[c.gender][action.InvalidAction],
		CanQueueAfter:   skillFrames[c.gender][action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
