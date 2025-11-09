package domain_test

import (
	"testing"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/stretchr/testify/assert"
)

func buildDummySkill() domain.Skill {
	return domain.Skill{
		SkillID:       1,
		Constant:      domain.NewNullableString("Constant"),
		Name:          domain.NewNullableString("Skill 1"),
		FileVersion:   1,
		Description:   domain.NewNullableString("Description 1"),
		MaxLevel:      domain.NewNullableInt32(3),
		SpAmount:      []int32{1, 2, 3},
		SeparateLevel: domain.NewNullableBool(true),
		AttackRange:   []int32{1, 2, 3},
		NeedSkillList: []domain.NeedSkillEntry{
			{SkillID: 2, Level: 1},
		},
	}
}

func TestSkill_Equals(t *testing.T) {
	fromSkill := buildDummySkill()
	toSkill := buildDummySkill()

	assert.True(t, fromSkill.Equals(toSkill))

	// Test description
	toSkill = buildDummySkill()
	toSkill.Description = domain.NewNullableString("Description 2")
	assert.False(t, fromSkill.Equals(toSkill))

	// Test SpAmount
	toSkill = buildDummySkill()
	toSkill.SpAmount = []int32{1, 2, 4}
	assert.False(t, fromSkill.Equals(toSkill))

	// Test SeparateLevel
	toSkill = buildDummySkill()
	toSkill.SeparateLevel = domain.NewNullableBool(false)
	assert.False(t, fromSkill.Equals(toSkill))

	// Test AttackRange
	toSkill = buildDummySkill()
	toSkill.AttackRange = []int32{1, 2, 4}
	assert.False(t, fromSkill.Equals(toSkill))

	// Test NeedSkillList
	toSkill = buildDummySkill()
	toSkill.NeedSkillList = []domain.NeedSkillEntry{
		{SkillID: 2, Level: 2},
	}
	assert.False(t, fromSkill.Equals(toSkill))
}
