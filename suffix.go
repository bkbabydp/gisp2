package gisp

import (
	p "github.com/Dwarfartisan/goparsec2"
)

// DotSuffix 表示带 dot 分割的后缀的表达式
func DotSuffix(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		d, err := DotParser(st)
		if err != nil {
			return nil, err
		}
		return dotSuffix(Dot{x, d.(Atom)})(st)
	}
}

func dotSuffix(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		d, err := p.Try(DotParser)(st)
		if err != nil {
			return x, nil
		}
		return dotSuffix(Dot{x, d.(Atom)})(st)
	}
}

// BracketSuffix 表示带 [] 后缀的表达式
func BracketSuffix(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		b, err := p.Try(BracketParser())(st)
		if err != nil {
			return nil, err
		}
		return bracketSuffix(Bracket{x, b.([]interface{})})(st)
	}
}

// BracketSuffixExt 带扩展环境，可以在指定的环境中解释[]中的token
func BracketSuffixExt(env Env) func(interface{}) p.P {
	return func(x interface{}) p.P {
		return func(st p.State) (interface{}, error) {
			b, err := p.Try(BracketParserExt(env))(st)
			if err != nil {
				return nil, err
			}
			return bracketSuffixExt(env)(Bracket{x, b.([]interface{})})(st)
		}
	}
}

func bracketSuffix(x interface{}) p.P {
	return func(st p.State) (interface{}, error) {
		b, err := p.Try(BracketParser())(st)
		if err != nil {
			return x, nil
		}
		return Bracket{x, b.([]interface{})}, nil
	}
}

func bracketSuffixExt(env Env) func(interface{}) p.P {
	return func(x interface{}) p.P {
		return func(st p.State) (interface{}, error) {
			b, err := p.Try(BracketParserExt(env))(st)
			if err != nil {
				return x, nil
			}
			return Bracket{x, b.([]interface{})}, nil
		}
	}
}

// DotSuffixParser 定义 dot 表达式判定
func DotSuffixParser(x interface{}) p.P {
	return p.Choice(p.Try(DotSuffix(x)), p.Return(x))
}

// BracketSuffixParser 定义 bracket 表达式判定
func BracketSuffixParser(x interface{}) p.P {
	return p.Choice(p.Try(BracketSuffix(x)), p.Return(x))
}

// SuffixParser 定义了后缀表达式的通用判定
func SuffixParser(prefix interface{}) p.P {
	suffix := p.Choice(p.Try(DotSuffix(prefix)), BracketSuffix(prefix))
	return func(st p.State) (interface{}, error) {
		s, err := suffix(st)
		if err != nil {
			return prefix, nil
		}
		return SuffixParser(s)(st)
	}
}

// SuffixParserExt 在后缀表达式判定中允许传入环境
func SuffixParserExt(env Env) func(interface{}) p.P {
	return func(prefix interface{}) p.P {
		suffix := p.Choice(p.Try(DotSuffix(prefix)), BracketSuffixExt(env)(prefix))
		return func(st p.State) (interface{}, error) {
			s, err := p.Try(suffix)(st)
			if err != nil {
				return prefix, nil
			}
			return SuffixParserExt(env)(s)(st)
		}
	}
}
