package timer

import (
	"container/list"
	"math"
)

// 时间轮数量
const cycleSize int = 9

// 时间轮持续时间
//
//	key:序号[0,...]
//	value:到期时间
var gCycleDuration [cycleSize]int64

func init() {
	for i := 0; i < len(gCycleDuration); i++ {
		gCycleDuration[i] = genDuration(i)
	}
	gCycleDuration[cycleSize-1] = math.MaxInt64
}

type cycle struct {
	data      list.List
	minExpire int64 //最小到期时间
}

func (p *cycle) init() {
	p.minExpire = math.MaxInt64
}

// 生成一个轮的时长
//
//	参数:
//		轮序号
//	返回值:
//		4,8,16,32,64,128,256,512,math.MaxInt64
func genDuration(idx int) int64 {
	return int64(1 << (uint)(idx+2))
}

// 根据 时长 找到时间轮的序号
// (当前为从头依次判断,适用于大多数数据 符合头部条件,若数据均匀分布,则适用于使用二分查找)
//
//	参数:
//		duration:时长
//	返回值:
//		轮序号
func findCycleIdx(duration int64) (idx int) {
	for k, v := range gCycleDuration {
		if duration <= v {
			return k
		} else {
			idx++
		}
	}
	return len(gCycleDuration) - 1
}

// 根据 时长 找到时间轮的序号 二分查找 (递归)
func binarySearchCycleIdxRecursion(low int, high int, duration int64) int {
	mid := low + (high-low)/2
	if duration <= gCycleDuration[mid] {
		if mid == 0 {
			return mid
		}
		if gCycleDuration[mid-1] < duration {
			return mid
		}
		return binarySearchCycleIdxRecursion(low, mid-1, duration)
	} else {
		if duration <= gCycleDuration[mid+1] {
			return mid + 1
		}
		return binarySearchCycleIdxRecursion(mid+1, high, duration)
	}
}

// 根据 时长 找到时间轮的序号 二分查找 (迭代)
func binarySearchCycleIdxIteration(duration int64) int {
	low, high := 0, len(gCycleDuration)-1
	for low <= high {
		mid := low + (high-low)/2
		if low == high {
			if gCycleDuration[mid] < duration {
				return mid + 1
			} else {
				return mid
			}
		}
		if gCycleDuration[mid] == duration {
			return mid
		} else if duration < gCycleDuration[mid] {
			high = mid - 1
		} else if duration > gCycleDuration[mid] {
			low = mid + 1
		}
	}
	return low
}

// 向前查找符合时间差的时间轮序号
//
//	参数:
//		duration: 到期 时长
//		idx: 轮序号 0 < idx
func findPrevCycleIdx(duration int64, idx int) int {
	for {
		if 0 != idx && duration <= gCycleDuration[idx-1] {
			idx--
		} else {
			break
		}
	}
	return idx
}
