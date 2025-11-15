package delayParser

import (
	"fmt"
	"os"
	"regexp"

	"github.com/guilherme-gm/ro-vis/extractor/internal/domain"
	"github.com/guilherme-gm/ro-vis/extractor/internal/ro/rostructs"
)

type SkillDelayParser interface {
	IsUpdateInRange(update *domain.Update) bool
	HasFiles(update *domain.Update) bool
	GetRelevantFiles() []*regexp.Regexp
	Parse(basePath string, change *domain.UpdateChange, skillIdTable map[string]int) []rostructs.SkillDelayV1
}

var SkillDelayListV1 = "data/luafiles514/lua files/skillinfoz/skilldelaylist.lub"
var SkillDelayListV1Regex = regexp.MustCompile(`(?i)^` + SkillDelayListV1 + `$`)

// SKFLAG_ variables are not present in lua, but probably provided by the RagExe,
// we are giving made up values to them so we can track them
// this list must be updated manually if new values are found
var SkFlags = map[string]int{
	"SKFLAG_NOREDUCT":                       1,
	"SKFLAG_DISABLE_FIXEDCASTING_REDUCTION": 2,
}

func validateFileConstats(filePath string, knownFlags map[string]int) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`SKFLAG_[A-Z0-9_]+`)
	matches := re.FindAllString(string(data), -1)
	if len(matches) == 0 {
		fmt.Println(" No falgs found")
		return nil
	}

	seen := make(map[string]struct{}, len(matches))
	for _, m := range matches {
		seen[m] = struct{}{}
	}

	var unknownFlags []string
	for flag := range seen {
		if _, ok := knownFlags[flag]; !ok {
			fmt.Println("Unknown flag:", flag)
			unknownFlags = append(unknownFlags, flag)
		}
	}

	if len(unknownFlags) > 0 {
		fmt.Println("has unknown")
		return fmt.Errorf("unknown flags: %v", unknownFlags)
	}

	return nil
}
