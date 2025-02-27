package rosaria

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var skillFrames []int

const skillHitmark = 24

func init() {
	skillFrames = frames.InitAbilSlice(51)
	skillFrames[action.ActionDash] = 38
	skillFrames[action.ActionJump] = 40
	skillFrames[action.ActionSwap] = 50
}

// Skill attack damage queue generator
// Includes optional argument "nobehind" for whether Rosaria appears behind her opponent or not (for her A1).
// Default behavior is to appear behind enemy - set "nobehind=1" to diasble A1 proc
func (c *char) Skill(p map[string]int) action.ActionInfo {
	// No ICD to the 2 hits
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Ravaging Confession (Hit 1)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSpear,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skill[0][c.TalentLvlSkill()],
		HitlagHaltFrames:   0.06 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: false,
	}

	// We always assume that A1 procs on hit 1 to simplify
	var a1cb combat.AttackCBFunc
	if p["nobehind"] != 1 {
		a1cb = c.a1()
	}
	var c4cb combat.AttackCBFunc
	if c.Base.Cons >= 4 {
		c.c4completed = false
		c4cb = c.c4
	}
	c.Core.QueueAttack(
		ai,
		combat.NewBoxHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: -1}, 2, 4),
		skillHitmark,
		skillHitmark,
		a1cb,
		c4cb,
	)

	// Rosaria E is dynamic, so requires a second snapshot
	//TODO: check snapshot timing here
	ai = combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Ravaging Confession (Hit 2)",
		AttackTag:          combat.AttackTagElementalArt,
		ICDTag:             combat.ICDTagNone,
		ICDGroup:           combat.ICDGroupDefault,
		StrikeType:         combat.StrikeTypeSlash,
		Element:            attributes.Cryo,
		Durability:         25,
		Mult:               skill[1][c.TalentLvlSkill()],
		HitlagHaltFrames:   0.09 * 60,
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
	}
	c.QueueCharTask(func() {
		//second hit is 14 frames after the first (if we exclude hitlag)
		c.Core.QueueAttack(
			ai,
			combat.NewCircleHitOnTarget(c.Core.Combat.Player(), combat.Point{Y: 0.5}, 2.8),
			0,
			0,
		)
		// Particles are emitted after the second hit lands
		c.Core.QueueParticle("rosaria", 3, attributes.Cryo, c.ParticleDelay)
	}, skillHitmark+14)

	c.SetCDWithDelay(action.ActionSkill, 360, 23)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(skillFrames),
		AnimationLength: skillFrames[action.InvalidAction],
		CanQueueAfter:   skillFrames[action.ActionDash], // earliest cancel
		State:           action.SkillState,
	}
}
