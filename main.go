package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sync/atomic"
	"time"

	"github.com/codegangsta/cli"
)

var (
	// CleanCount 清理文件数量
	CleanCount int64
	// StartTime 开始时间
	StartTime time.Time
)

func main() {
	StartTime = time.Now()
	app := cli.NewApp()
	app.Name = "fileclean"
	app.Version = "0.1.0"
	app.Usage = "文件清理程序"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "dir, d",
			Value: ".",
			Usage: "文件目录",
		},
		cli.BoolFlag{
			Name:  "recur, r",
			Usage: "递归处理",
		},
		cli.StringSliceFlag{
			Name:  "name",
			Usage: "文件名",
		},
		cli.StringSliceFlag{
			Name:  "reg",
			Usage: "正则过滤的文件名",
		},
		cli.BoolFlag{
			Name:  "exclude, e",
			Usage: "排除当前指定的文件名",
		},
		cli.BoolFlag{
			Name:  "all, a",
			Usage: "清理所有文件",
		},
	}
	app.Action = handleAction
	app.Run(os.Args)
}

// Config 配置参数
type Config struct {
	Dir     string
	Name    []string
	Reg     []string
	Recur   bool
	Exclude bool
	All     bool
}

// handleAction 执行清理函数
func handleAction(ctx *cli.Context) {
	defer func() {
		if err := recover(); err != nil {
			cli.ShowAppHelp(ctx)
			fmt.Println(err)
		}
	}()
	config := Config{
		Dir:     ctx.String("dir"),
		Name:    ctx.StringSlice("name"),
		Reg:     ctx.StringSlice("reg"),
		Recur:   ctx.Bool("recur"),
		Exclude: ctx.Bool("exclude"),
		All:     ctx.Bool("all"),
	}
	if !config.All && (len(config.Name) == 0 && len(config.Reg) == 0) {
		panic("未知的文件名！")
	}
	stat, err := os.Lstat(config.Dir)
	if os.IsNotExist(err) || !stat.IsDir() {
		panic("未知的文件目录！")
	}
	root, err := filepath.Abs(config.Dir)
	if err != nil {
		panic(err)
	}
	err = handleRecursive(root, &config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\r [清理完成] 耗时：%dms,清理文件数：%d \t", time.Duration(time.Now().Sub(StartTime).Nanoseconds())/time.Millisecond, CleanCount)
}

// handleRecursive 递归处理文件
func handleRecursive(root string, config *Config) error {
	file, err := os.Open(root)
	if err != nil {
		return err
	}
	defer func() {
		file.Close()
		exist, err := checkExistFile(root)
		if err != nil {
			panic(err)
		}
		if !exist {
			os.Remove(root)
		}
	}()
	files, err := file.Readdir(0)
	if err != nil {
		return nil
	}
	for _, info := range files {
		if info.IsDir() {
			if config.Recur {
				if err := handleRecursive(path.Join(root, info.Name()), config); err != nil {
					return err
				}
			}
			continue
		}
		matched, err := checkMatched(config, info.Name())
		if err != nil {
			return err
		}
		if config.All {
			if !(config.Exclude && (len(config.Name) > 0 || len(config.Reg) > 0) && matched) {
				if err := handleRemove(root, info); err != nil {
					return err
				}
			}
			continue
		}
		if matched {
			if err := handleRemove(root, info); err != nil {
				return err
			}
		}
	}
	return nil
}

// checkMatched 检查文件名匹配
func checkMatched(config *Config, fileName string) (matched bool, err error) {
	for _, name := range config.Name {
		if name == fileName {
			matched = true
			return
		}
	}
	for _, reg := range config.Reg {
		matched, err = regexp.MatchString(reg, fileName)
		if err != nil || matched {
			return
		}
	}
	return
}

// checkExistFile 检查目录是否存在文件
func checkExistFile(name string) (exist bool, err error) {
	file, err := os.Open(name)
	if err != nil {
		return
	}
	defer file.Close()
	names, _ := file.Readdirnames(1)
	if len(names) > 0 {
		exist = true
	}
	return
}

func handleRemove(root string, info os.FileInfo) error {
	if err := os.Remove(path.Join(root, info.Name())); err != nil {
		return err
	}
	atomic.AddInt64(&CleanCount, 1)
	fmt.Printf("\r [...] 用时：%dms,清理文件数：%d \t", time.Duration(time.Now().Sub(StartTime).Nanoseconds())/time.Millisecond, CleanCount)
	return nil
}
