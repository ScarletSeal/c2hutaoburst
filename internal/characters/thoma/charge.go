package thoma

import (
	"github.com/genshinsim/gcsim/internal/frames"
	"github.com/genshinsim/gcsim/pkg/core/action"
	"github.com/genshinsim/gcsim/pkg/core/attributes"
	"github.com/genshinsim/gcsim/pkg/core/combat"
)

var chargeFrames []int

const chargeHitmark = 27

func init() {
	chargeFrames = frames.InitAbilSlice(64)
	chargeFrames[action.ActionDash] = 34
	chargeFrames[action.ActionJump] = 34
	chargeFrames[action.ActionSwap] = 63
}

// Charge attack damage queue generator
// Very standard - consistent with other characters like Xiangling
func (c *char) ChargeAttack(p map[string]int) action.ActionInfo {
	ai := combat.AttackInfo{
		ActorIndex:         c.Index,
		Abil:               "Charge",
		AttackTag:          combat.AttackTagExtra,
		ICDTag:             combat.ICDTagExtraAttack,
		ICDGroup:           combat.ICDGroupPoleExtraAttack,
		StrikeType:         combat.StrikeTypeSpear,
		Element:            attributes.Physical,
		Durability:         25,
		Mult:               charged[c.TalentLvlAttack()],
		HitlagFactor:       0.01,
		CanBeDefenseHalted: true,
		IsDeployable:       true,
	}

	c.Core.QueueAttack(
		ai,
		combat.NewCircleHit(
			c.Core.Combat.Player(),
			c.Core.Combat.PrimaryTarget(),
			nil,
			0.8,
		),
		chargeHitmark,
		chargeHitmark,
	)

	return action.ActionInfo{
		Frames:          frames.NewAbilFunc(chargeFrames),
		AnimationLength: chargeFrames[action.InvalidAction],
		CanQueueAfter:   chargeHitmark,
		State:           action.ChargeAttackState,
	}
}
