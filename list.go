package gisp

import (
	"fmt"
	"reflect"
	"strings"
)

// List 实现基本的 List 类型
type List []interface{}

func (list List) String() string {
	frags := []string{}
	for _, item := range list {
		frags = append(frags, fmt.Sprintf("%v", item))
	}
	body := strings.Join(frags, " ")
	return fmt.Sprintf("(%s)", body)
}

// Eval 实现 Lisp.Eval 方法
func (list List) Eval(env Env) (interface{}, error) {
	l := len(list)
	if l == 0 {
		return nil, nil
	}
	var lisp interface{}
	switch fun := list[0].(type) {
	case Atom:
		var ok bool
		if lisp, ok = env.Lookup(fun.Name); !ok {
			return nil, fmt.Errorf("any callable named %s not found", fun.Name)
		}
	case List:
		var err error
		lisp, err = fun.Eval(env)
		if err != nil {
			return nil, err
		}
	case Func:
		lisp = fun
	case Expr:
		lisp = fun
	}
	switch item := lisp.(type) {
	case Expr:
		return item(env)(list[1:]...)
	case Task:
		return item.Eval(env)
	case Go:
		return item.Eval(env)
	case Lambda:
		lisp, err := item.Task(list[1:]...)
		if err != nil {
			return nil, err
		}
		return lisp.Eval(env)
	case Func:
		lisp, err := item.Task(list[1:]...)
		if err != nil {
			return nil, err
		}
		return lisp.Eval(env)
	case Let:
		return item.Eval(env)
	default:
		return nil, fmt.Errorf("%v:%v is't callable", list[0], reflect.TypeOf(list[0]))
	}

}
