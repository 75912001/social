package handler

import (
	liblog "social/lib/log"
)

func OnEventDefault(v interface{}) error {
	switch t := v.(type) {
	default:
		liblog.GetInstance().Errorf("non-existent event:%v %v", v, t)
	}
	return nil
}
