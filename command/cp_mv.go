package baidupcscmd

import (
	"fmt"
	"github.com/iikira/BaiduPCS-Go/baidupcs"
	"path"
)

// RunCopy 执行 批量拷贝文件/目录
func RunCopy(paths ...string) {
	runCpMvOp("copy", paths...)
}

// RunMove 执行 批量 重命名/移动 文件/目录
func RunMove(paths ...string) {
	runCpMvOp("move", paths...)
}

func runCpMvOp(op string, paths ...string) {
	err := cpmvPathValid(paths...)
	if err != nil {
		fmt.Println(err)
		return
	}

	paths = getAllPaths(paths...)
	froms, to := cpmvParsePath(paths...)

	toInfo, err := info.FilesDirectoriesMeta(to)
	if err != nil {
		if len(froms) != 1 {
			fmt.Println(err)
			return
		}

		if op == "copy" { // 拷贝
			err = info.Copy(baidupcs.CpMvJSON{
				From: froms[0],
				To:   to,
			})
			if err != nil {
				fmt.Println(err)
				fmt.Println("文件/目录拷贝失败: ")
				fmt.Printf("%s <-> %s\n", froms[0], to)
				return
			}
			fmt.Println("文件/目录拷贝成功: ")
			fmt.Printf("%s <-> %s\n", froms[0], to)
		} else { // 重命名
			err = info.Rename(froms[0], to)
			if err != nil {
				fmt.Println(err)
				fmt.Println("重命名失败: ")
				fmt.Printf("%s -> %s\n", froms[0], to)
				return
			}
			fmt.Println("重命名成功: ")
			fmt.Printf("%s -> %s\n", froms[0], to)
		}
		return
	}

	if !toInfo.Isdir {
		fmt.Printf("目标 %s 不是一个目录, 操作失败\n", toInfo.Path)
		return
	}

	cj := new(baidupcs.CpMvJSONList)
	cj.List = make([]baidupcs.CpMvJSON, len(froms))
	for k := range froms {
		cj.List[k].From = froms[k]
		cj.List[k].To = to + "/" + path.Base(froms[k])
	}

	switch op {
	case "copy":
		err = info.Copy(cj.List...)
		if err != nil {
			fmt.Println(err)
			fmt.Println("操作失败, 以下文件/目录拷贝失败: ")
			fmt.Println(cj)
			return
		}
		fmt.Println("操作成功, 以下文件/目录拷贝成功: ")
		fmt.Println(cj)
	case "move":
		err = info.Move(cj.List...)
		if err != nil {
			fmt.Println(err)
			fmt.Println("操作失败, 以下文件/目录移动失败: ")
			fmt.Println(cj)
			return
		}
		fmt.Println("操作成功, 以下文件/目录移动成功: ")
		fmt.Println(cj)
	default:
		panic("Unknown op:" + op)
	}
	return
}

// cpmvPathValid 检查路径的有效性
func cpmvPathValid(paths ...string) (err error) {
	if len(paths) <= 1 {
		return fmt.Errorf("参数不完整")
	}

	return nil
}

// cpmvParsePath 解析路径
func cpmvParsePath(paths ...string) (froms []string, to string) {
	if len(paths) == 0 {
		return nil, ""
	}
	froms = paths[:len(paths)-1]
	to = paths[len(paths)-1]
	return
}
