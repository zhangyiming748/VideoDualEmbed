package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zhangyiming748/FastTranslate"
	"github.com/zhangyiming748/FastWhisper"
	"github.com/zhangyiming748/FastYtdlp"
	"github.com/zhangyiming748/archive"
	"github.com/zhangyiming748/finder"
	"github.com/zhangyiming748/lumberjack"
)

var (
	// Download 命令参数
	link   string
	proxy  string
	cookie string

	// Whisper 命令参数
	level    string
	location string
	language string
	root     string
	format   string

	// Trans 命令参数
	transRoot  string
	transProxy string

	//merge 命令参数
	dir string
)

func main() {
	setLog()
	var rootCmd = &cobra.Command{
		Use:   "yt-whisper-bilingual",
		Short: "YouTube双语字幕生成工具",
		Long:  `一个自动化工具，用于从YouTube视频生成双语SRT字幕文件`,
	}

	// Download 子命令
	var downloadCmd = &cobra.Command{
		Use:   "download",
		Short: "下载YouTube视频",
		Long:  `使用ytdlp下载指定URL的视频`,
		Run: func(cmd *cobra.Command, args []string) {
			download(link, proxy, cookie)
		},
	}

	downloadCmd.Flags().StringVar(&link, "link", "", "需要使用ytdlp下载的全部连接的列表文件路径 (必填)")
	downloadCmd.Flags().StringVar(&proxy, "proxy", "", "下载过程中使用的代理")
	downloadCmd.Flags().StringVar(&cookie, "cookie", "", "下载过程中需要提供的cookie的文件路径")
	downloadCmd.MarkFlagRequired("link")

	// Whisper 子命令
	var whisperCmd = &cobra.Command{
		Use:   "whisper",
		Short: "使用Whisper生成字幕",
		Long:  `使用OpenAI Whisper将视频音频转录为原始语言字幕`,
		Run: func(cmd *cobra.Command, args []string) {
			whisper(level, location, language, root, format)
		},
	}

	whisperCmd.Flags().StringVar(&level, "level", "medium.en", "模型等级 (默认: medium)")
	whisperCmd.Flags().StringVar(&location, "location", "/data/models", "模型位置 (默认: /data/models)")
	whisperCmd.Flags().StringVar(&language, "language", "English", "视频语言 (默认: English)")
	whisperCmd.Flags().StringVar(&root, "root", "/data", "视频所在文件夹 (默认: /data)")
	whisperCmd.Flags().StringVar(&format, "format", "srt", "输出字幕的格式srt或all (默认: srt)")

	// Trans 子命令
	var transCmd = &cobra.Command{
		Use:   "trans",
		Short: "翻译字幕文件",
		Long:  `使用translate-shell将原始字幕翻译为目标语言`,
		Run: func(cmd *cobra.Command, args []string) {
			trans(transRoot, transProxy)
		},
	}

	transCmd.Flags().StringVar(&transRoot, "root", "/data", "原始字幕文件所在的位置 (默认: /data)")
	transCmd.Flags().StringVar(&transProxy, "proxy", "", "翻译过程中使用的代理")

	var mergeCmd = &cobra.Command{
		Use:   "merge",
		Short: "内嵌字幕",
		Long:  `使用ffmpeg将合并翻译后的双语字幕内嵌到视频文件中`,
		Run: func(cmd *cobra.Command, args []string) {
			merge(dir)
		},
	}

	mergeCmd.Flags().StringVar(&dir, "root", "/data", "原始字幕文件所在的位置 (默认: /data)")
	// 添加子命令到根命令
	rootCmd.AddCommand(mergeCmd)
	rootCmd.AddCommand(downloadCmd)
	rootCmd.AddCommand(whisperCmd)
	rootCmd.AddCommand(transCmd)

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("命令运行出现致命错误:%v\n", err)
	}
}

func download(link, proxy, cookie string) {
	fmt.Printf("开始下载视频...\n")
	fmt.Printf("链接文件: %s\n", link)
	fmt.Printf("代理设置: %s\n", proxy)
	fmt.Printf("Cookie文件: %s\n", cookie)
	FastYtdlp.Download(link, proxy, cookie)
}

func whisper(ModelType, ModelDir, Language, VideoRoot, SubtitleFormat string) {
	fmt.Printf("开始生成字幕...\n")
	fmt.Printf("模型等级: %s\n", ModelType)
	fmt.Printf("模型位置: %s\n", ModelDir)
	fmt.Printf("视频语言: %s\n", Language)
	fmt.Printf("视频目录: %s\n", VideoRoot)
	fmt.Printf("字幕格式: %s\n", SubtitleFormat)

	videos := finder.FindAllVideos(VideoRoot)
	for _, video := range videos {
		fc := new(FastWhisper.WhisperConfig)
		fc.ModelType = ModelType
		fc.ModelDir = ModelDir
		fc.Language = Language
		fc.VideoRoot = video
		fc.Format = SubtitleFormat
		// 调用获取字幕的方法
		FastWhisper.GetSubtitle(fc)
	}
}

func trans(SrtRoot, proxy string) {
	fmt.Printf("开始翻译字幕...\n")
	fmt.Printf("字幕目录: %s\n", SrtRoot)
	fmt.Printf("代理设置: %s\n", proxy)

	srts := finder.FindAllFiles(SrtRoot)
	for _, srt := range srts {
		if strings.HasSuffix(srt, "origin.srt") {
			continue
		} else if strings.HasSuffix(srt, ".srt") {
			FastTranslate.TranslateSrt(srt, proxy)
		}
	}
}

func merge(root string) {
	log.Printf("当前查找的目录是%v\t不包含子文件夹", root)
	videos := finder.FindAllVideosInRoot(root)
	for _, video := range videos {
		log.Printf("找到的视频文件:%v\n", video)
		srt := strings.Replace(video, filepath.Ext(video), ".srt", 1)
		if exist(srt) {
			log.Printf("视频文件:%v存在对应的字幕文件:%v\n", video, srt)
			if err := archive.MergeMp4WithSameNameSrt(video, srt); err != nil {
				log.Printf("此次转换出现错误:%v\n", err)
				continue
			} else {
				os.Remove(video)
				os.Remove(srt)
			}
		} else {
			log.Printf("视频文件:%v不存在对应的字幕文件:%v\n", video, srt)
			continue
		}
	}
}
func setLog() {
	// 设置全局时区为Asia/Shanghai
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("无法加载时区 Asia/Shanghai: %v", err)
	} else {
		time.Local = location
	}
	// 创建一个用于写入文件的Logger实例
	fileLogger := &lumberjack.Logger{
		Filename:   "ywt.log",
		MaxSize:    1, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
	}
	err = fileLogger.Rotate()
	if err != nil {
		log.Println("转换新日志文件失败", err)
	}
	consoleLogger := log.New(os.Stdout, "CONSOLE: ", log.LstdFlags)
	log.SetOutput(io.MultiWriter(fileLogger, consoleLogger.Writer()))
	log.SetFlags(log.Ltime | log.Lshortfile)
}

// 判断给出的绝对路径是否是一个存在的文件
func exist(fp string) bool {
	info, err := os.Stat(fp)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
