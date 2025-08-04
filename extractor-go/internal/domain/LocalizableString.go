package domain

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/guilherme-gm/ro-vis/extractor/internal/utils/base64Utils"
)

type StringKind string

const (
	StringKindStatic    StringKind = "static"
	StringKindLocalized StringKind = "localized"
)

type LocalizableString struct {
	Kind  StringKind
	Value string
}

func NewLocalizableStringFromNavi(name string) LocalizableString {
	if strings.HasPrefix(name, "\x1C") {
		if !strings.HasSuffix(name, "\x1C") {
			fmt.Println("Invalid name1:", name)
			panic("Invalid name1")
		}

		value, err := base64Utils.DecodeBase64ToUInt64(name[1 : len(name)-1])
		if err != nil {
			fmt.Println("Error decoding name1:", err)
			panic(err)
		}

		return LocalizableString{
			Kind:  StringKindLocalized,
			Value: strconv.FormatUint(value, 10),
		}
	} else {
		return LocalizableString{
			Kind:  StringKindStatic,
			Value: name,
		}
	}
}
