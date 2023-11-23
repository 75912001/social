package handler

import xrlog "social/pkg/lib/log"

func OnEventDefault(v interface{}) error {
	switch t := v.(type) {
	default:
		xrlog.GetInstance().Fatalf("non-existent event:%v %v", v, t)
	}
	return nil
}
