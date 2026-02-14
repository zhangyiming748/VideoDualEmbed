# VideoDualEmbed

这是一个自动化工具，用于从YouTube视频生成双语SRT字幕文件。基于 Cobra 库重构的命令行工具，支持四种主要功能：

## 编译

```bash
go build -o main main.go
```

## 使用方法

### 1. download - 下载YouTube视频

```bash
# 基本用法（link参数为必填）
./main download --link /path/to/links.txt

# 完整参数示例
./main download --link /home/zen/post.link --proxy http://192.168.110.78:8889 --cookie /home/zen/youtube.cookie
```

**参数说明：**
- `--link`: 需要下载的视频链接列表文件路径（必填）
- `--proxy`: 下载过程中使用的代理服务器
- `--cookie`: 下载过程中需要的cookie文件路径

### 2. whisper - 生成字幕

```bash
# 使用默认参数
./main whisper

# 自定义参数
./main whisper --level medium --location /data/models --language English --root /data --format srt
```

**参数说明：**
- `--level`: 模型等级（默认: medium）
- `--location`: 模型文件位置（默认: /data/models）
- `--language`: 视频语言（默认: English）
- `--root`: 视频文件所在目录（默认: /data）
- `--format`: 输出字幕格式，srt或all（默认: srt）

### 3. trans - 翻译字幕

```bash
# 使用默认参数
./main trans

# 自定义参数
./main trans --root /data --proxy http://192.168.110.78:8889
```

**参数说明：**
- `--root`: 原始字幕文件所在目录（默认: /data）
- `--proxy`: 翻译过程中使用的代理服务器

### 4. merge - 内嵌字幕到视频

```bash
# 使用默认参数
./main merge

# 自定义参数
./main merge --root /data
```

**功能说明：**
- 自动查找指定目录下的视频文件
- 为每个视频文件寻找对应的.srt字幕文件
- 使用FFmpeg将字幕内嵌到视频中（编码为H.265/AAC格式）
- 成功内嵌后自动删除原始视频和字幕文件
- 输出文件命名为原文件名+_subInside.mp4

**参数说明：**
- `--root`: 视频和字幕文件所在目录（默认: /data）

## 查看帮助信息

```bash
# 查看主命令帮助
./main --help

# 查看特定子命令帮助
./main download --help
./main whisper --help
./main trans --help
./main merge --help
```

## 参数顺序说明

使用 Cobra 库后，参数顺序不再重要，以下命令效果相同：

```bash
./main download --link file.txt --proxy proxy_url --cookie cookie.txt
./main download --proxy proxy_url --link file.txt --cookie cookie.txt
./main download --cookie cookie.txt --link file.txt --proxy proxy_url
```

## 错误处理

- 必填参数缺失时会显示明确的错误信息
- 参数类型错误时会有相应提示
- 命令执行过程中的错误会详细记录日志

## 工作流程示例

完整的双语字幕生成流程：
```bash
# 1. 下载视频
./main download --link links.txt

# 2. 生成原始语言字幕
./main whisper --root /data/videos --language English

# 3. 翻译为双语字幕
./main trans --root /data/videos

# 4. 将字幕内嵌到视频中
./main merge --root /data/videos
```