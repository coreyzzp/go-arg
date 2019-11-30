package sargs

import (
	"fmt"
	"reflect"
	"strings"
)

type TagType int

const (
	KTagTypeNil = iota
	KTagTypeOption
	KTagTypeArgs
	KTagTypeSubCmd
)

type argTag struct {
	tagType    TagType
	ignore     bool
	isPosOpt   bool
	isRequired bool
	subcmdName string
	shortAlias string
	longAlias  string
	help       string
	defval     string
	envFrom    string
}

const (
	TAG_SYMBOL_ARG        = "arg"
	TAG_SYMBOL_POSITIONAL = "positional"
	TAG_SYMBOL_ENV        = "env"
	TAG_SYMBOL_REQUIRED   = "required"
	TAG_SYMBOL_SUBCOMMAND = "subcommand"
)

// 原来的 https://github.com/alexflint/go-arg 支持以下类型的tag
// 其中arg参数中，以-或者--开头，且不为'-'的字段，会进行忽略
// `arg:"positional,env:envalue,required,subcommand,-,-x,-XX"`
// `help:"help text"`
// 扩展:
// `default:"defaultValue"`

func splitKV(in string, sp1, sp2 string) map[string]string {
	m := make(map[string]string)
	if in == "" {
		return m
	}
	sp := strings.Split(in, sp1)
	for _, spp := range sp {
		spp2 := strings.Split(spp, sp2)
		if len(spp2) == 2 {
			m[spp2[0]] = spp2[1]
		} else if len(spp2) == 1 {
			m[spp2[0]] = ""
		}
	}
	return m
}

func newArgTag(tag string, opt *options) (atag *argTag, err error) {
	var (
		fieldTag   = reflect.StructTag(tag)
		argTextMap = splitKV(fieldTag.Get(TAG_SYMBOL_ARG), ",", ":")
	)

	atag = &argTag{}
	atag.tagType = KTagTypeOption
	atag.defval = fieldTag.Get("default")
	atag.help = fieldTag.Get("help")

	// 没有设置arg参数
	if len(argTextMap) == 0 {
		if opt.tagopt.ignoreBySpecifyArgsEmpty {
			// 设置了 ignoreBySpecifyArgsEmpty 参数，但是arg参数为空，但也认为这个属于option
			atag.ignore = false
		} else {
			// 否则没标注就认为不属于option
			atag.ignore = true
		}
		return
	}

	// 也可以通过 arg:"-" 来忽略
	if _, ok := argTextMap["-"]; ok {
		atag.ignore = true
		return
	}

	for k, v := range argTextMap {
		switch k {
		case TAG_SYMBOL_POSITIONAL:
			atag.isPosOpt = true
			atag.tagType = KTagTypeArgs
		case TAG_SYMBOL_ENV:
			atag.envFrom = v
		case TAG_SYMBOL_REQUIRED:
			atag.isRequired = true
		case TAG_SYMBOL_SUBCOMMAND:
			atag.subcmdName = v
			atag.tagType = KTagTypeSubCmd
		default:
			switch {
			case v != "":
				// 指定的参数不应该有=
				err = fmt.Errorf("unknown arg option:%s=%s", k, v)
				return
			case len(k) < 2:
				// 至少有一个symbol
				err = fmt.Errorf("unknown arg option:%s=%s", k, v)
				return
			case strings.HasPrefix(k, "--"):
				// 只能设置一个long或者short的args
				if atag.longAlias != "" {
					err = fmt.Errorf("multiple long alias setting prev %s, now:%s", atag.longAlias, k)
					return
				}
				atag.longAlias = k
			case strings.HasPrefix(k, "-"):
				if atag.shortAlias != "" {
					err = fmt.Errorf("multiple short alias setting prev %s, now:%s", atag.shortAlias, k)
					return
				}
				atag.shortAlias = k
			default:
				err = fmt.Errorf("unknown arg option:%s", k)
				return
			}
		}
	}

	// TODO: tag合法性判断
	return
}
