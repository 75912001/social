package util

import (
	"bytes"
	cryptorand "crypto/rand"
	"fmt"
	"github.com/pkg/errors"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	libconstant "social/pkg/lib/constant"
	liberror "social/pkg/lib/error"
	"strconv"
	"strings"
	"unsafe"
)

// GetCurrentPath 获取当前程序 所在 路径
//
//	[×] 支持 win 下 link/快捷方式
//	[✔] 支持 linux ln -s
//	on linux host:ln -s /home/xxx/battle001/battle-linux.exe.001 /home/xxx/zone2/battle001/btl.exe
//	执行btl.exe 输出:service current path:/home/xxx/zone2/battle001
func GetCurrentPath() (currentPath string, err error) {
	if currentPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
		return "", errors.WithMessage(err, GetCodeLocation(1).String())
	}
	return currentPath, nil
}

// IsLittleEndian 判断 是小端
func IsLittleEndian() bool {
	var value int32 = 1 // 占4byte 转换成16进制 0x00 00 00 01
	// 大端(16进制)：00 00 00 01
	// 小端(16进制)：01 00 00 00
	return *(*byte)(unsafe.Pointer(&value)) == 1
}

// If 三目运算符
// NOTE 传递的实参,会在调用时计算
func If(condition bool, trueVal interface{}, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}

// SplitString2Uint32 拆分字符串, 返回 uint32 类型的 slice
func SplitString2Uint32(s string, sep string) (u32Slice []uint32, err error) {
	if 0 == len(s) {
		return u32Slice, nil
	}

	slice := strings.Split(s, sep)
	var u64 uint64
	for _, v := range slice {
		if u64, err = strconv.ParseUint(v, 10, 32); err != nil {
			return u32Slice, errors.WithMessage(err, GetCodeLocation(1).String())
		}
		u32Slice = append(u32Slice, uint32(u64))
	}
	return u32Slice, nil
}

// SplitString2Map 拆分字符串, 返回key为uint32类型、val为int64类型的map
func SplitString2Map(s string, sep1 string, sep2 string) (map[uint32]int64, error) {
	slice := strings.Split(s, sep1)
	m := make(map[uint32]int64)
	var err error
	for _, v := range slice {
		if 0 == len(v) {
			continue
		}
		sliceAttr := strings.Split(v, sep2)
		if len(sliceAttr) != 2 {
			return nil, errors.WithMessage(liberror.Param, GetCodeLocation(1).String())
		}
		var idUint64 uint64
		var valInt64 int64
		if idUint64, err = strconv.ParseUint(sliceAttr[0], 10, 32); err != nil {
			return nil, errors.WithMessage(err, GetCodeLocation(1).String())
		}
		if valInt64, err = strconv.ParseInt(sliceAttr[1], 10, 32); err != nil {
			return nil, errors.WithMessage(err, GetCodeLocation(1).String())
		}
		m[uint32(idUint64)] = valInt64
	}
	return m, nil
}

// WeightedRandom 从权重中选出序号.[0 ... ]
//
//	NOTE 参数 weights 必须有长度
//	参数:
//		weights:权重
//	返回值:
//		idx:weights 的序号 idx
func WeightedRandom(weights []uint32) (idx int, err error) {
	var sum int64
	for _, v := range weights {
		sum += int64(v)
	}
	if sum == 0 { //weights slice 中 无数据 / 数据都为0
		return 0, errors.WithMessage(liberror.Param, GetCodeLocation(1).String())
	}

	r := rand.Int63n(sum) + 1
	for i, v := range weights {
		if r <= int64(v) {
			return i, nil
		}
		r -= int64(v)
	}
	return 0, errors.WithMessage(liberror.System, GetCodeLocation(1).String())
}

// CodeLocation 代码位置
type CodeLocation struct {
	FileName string //文件名
	FuncName string //函数名
	Line     int    //行数
}

// Error 错误信息
func (p *CodeLocation) Error() string {
	return fmt.Sprintf("file:%v line:%v func:%v", p.FileName, p.Line, p.FuncName)
}

// String 错误信息
func (p *CodeLocation) String() string {
	return p.Error()
}

// GetCodeLocation 获取代码位置
//
//	参数:
//		skip:The argument skip is the number of stack frames to ascend, with 0 identifying the caller of Caller.
func GetCodeLocation(skip int) *CodeLocation {
	c := &CodeLocation{
		FileName: libconstant.Unknown,
		FuncName: libconstant.Unknown,
	}

	pc, fileName, line, ok := runtime.Caller(skip)

	if ok {
		c.FileName = fileName
		c.Line = line
		c.FuncName = runtime.FuncForPC(pc).Name()
	}
	return c
}

// GenRandomString 生成随机字符串
//
//	参数:
//		len:需要生成的长度
func GenRandomString(len uint32) (container string, err error) {
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	bigInt := big.NewInt(int64(bytes.NewBufferString(str).Len()))
	for i := uint32(0); i < len; i++ {
		if randomInt, err := cryptorand.Int(cryptorand.Reader, bigInt); err != nil {
			return "", errors.WithMessage(err, GetCodeLocation(1).String())
		} else {
			container += string(str[randomInt.Int64()])
		}
	}
	return container, nil
}

// RandomInt 生成范围内的随机值
//
//	参数:
//		min:最小值
//		max:最大值
func RandomInt(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

// IsDuplicateUint32 是否有重复uint32
func IsDuplicateUint32(uint32Slice []uint32) bool {
	set := make(map[uint32]struct{})
	for _, v := range uint32Slice {
		if _, ok := set[v]; ok {
			return true
		}
		set[v] = struct{}{}
	}
	return false
}

// GenNORepeatIdx 生成 不重复 序号
//
//	参数:
//		set:已有数据
//		uint32Slice:从该slice中随机一个,与set中不重复
//	返回值:
//		unit32Slice 中的index
func GenNORepeatIdx(set map[uint32]struct{}, uint32Slice []uint32) (int, error) {
	var slice []int
	for k, v := range uint32Slice {
		if _, ok := set[v]; ok {
			continue
		}
		slice = append(slice, k)
	}
	if len(slice) == 0 {
		return 0, errors.WithMessage(liberror.NonExistent, GetCodeLocation(1).String())
	}
	return slice[rand.Intn(len(slice))], nil
}

// GenRandValue 生成 随机值
//
//	参数:
//		except:排除 数据
//		uint32Slice:从该slice中随机一个,与except不重复
//	返回值:
//		unit32Slice 中的值
func GenRandValue(except uint32, uint32Slice []uint32) (uint32, error) {
	var slice []uint32
	for _, v := range uint32Slice {
		if v == except {
			continue
		}
		slice = append(slice, v)
	}
	if len(slice) == 0 {
		return 0, errors.WithMessage(liberror.NonExistent, GetCodeLocation(1).String())
	}
	return slice[rand.Intn(len(slice))], nil
}

// Command 调用 linux 命令
//
//	参数:
//		args:"chmod +x /xx/xx/x.sh"
func Command(args string) (outStr string, errStr string, err error) {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", args)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())
	if err != nil {
		return outStr, errStr, errors.WithMessage(err, GetCodeLocation(1).String())
	}
	return outStr, errStr, nil
}
