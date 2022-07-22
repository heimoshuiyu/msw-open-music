# MSW Open Music Project

[![CI](https://github.com/heimoshuiyu/msw-open-music/actions/workflows/build.yml/badge.svg)](https://github.com/heimoshuiyu/msw-open-music/actions/workflows/build.yml)

> 找一首歌最好的方法是：打开一个超长的歌单，然后随机播放，直到你找到为止。

一个 💪 轻量级 ⚡️ 高性能 🖥️ 跨平台的 个人音乐串流平台。管理你现有的音乐文件并在其他设备上播放。

前端网页应用基于 `react.js` 和 `water.css` 构建。后端服务器程序使用 `golang` 和 `sqlite` 构建。

## 介绍

截图

![demo1](demo1.jpg)

### 功能特点

- 🔎 索引现有的音乐文件，并记录文件名和文件夹元信息

- 📕 使用 文件夹 📁 标签 🏷️ 评论 💬 来管理你的音乐。

- 🌐 提供一个轻量高效的网页前端并支持多种语言。

- 👥 支持多用户。

- 🔥 调用 `ffmpeg` 配合可自定义的预设配置来转码你的音乐。

- 🔗 分享音乐链接给好友！

### 如果你遇到过这样的烦恼...你就是目标用户

- 硬盘上存了一堆音乐，但没有一个很好的播放器. 🖴

- 下载了体积非常大的无损音乐，在设备间移动很困难. 🎵

- 想要在其他 电脑/手机 上听 电脑/服务器 上储存的音乐. 😋

- 想给你的好友分享本地音乐. 😘

## 使用方法

1. 修改 `config.json` 配置文件中的 `secret` 值

2. 运行后端服务器程序 `msw-open-music.exe` 或者 `msw-open-music`. 服务默认监听 8080 端口。 然后打开 <http://127.0.0.1:8080> 去创建的一个管理员帐号。

前端 HTML 文件存放在 `web/build` 目录下。

### 创建第一个管理员帐号

第一个创建的管理员帐号会被自动激活，其他后续创建的管理员帐号需要管理员手动激活。

请前往注册页面，选择角色为 管理员，然后注册第一个管理员帐号。

#### config.json

- `secret` 字符串类型。用来加密 session 会话。

- `database_name` 字符串类型。`sqlite3` 数据库的文件名。如果不存在，会自动创建。
- `addr` 字符串类型。监听地址和端口。
- `ffmpeg_config_list` 列表类型。预设的 `ffmpeg` 配置文件。包含 `ffmpegConfig` 对象。
- `file_life_time` 整数类型（秒）。临时文件的生命周期。如果临时文件超过这个时间没有被访问，那么将会被自动删除。
- `cleaner_internal` 整数类型（秒）。`tmpfs` 检查临时文件的间隔时间。
- `root` 字符串类型。存放临时文件的目录。默认是 `/tmp`。**Windows用户请修改成可用的目录**。如果不存在，将会被自动创建。
- `permission` 各个 API 的权限等级。
  - `0` 无需任何权限。
  - `1` 需要管理员（最高级别）权限等级。
  - `2` 需要普通用户权限等级，也就是说，管理员和普通用户都有权访问此等级的 API ，而 匿名用户 则没有权限访问。
  - 如果你想避免 API 被滥用，可以调整下面 5 个与串流相关的 API 权限等级。
    - `/get_file` 使用 `io.copy()` 方法串流
    - `/get_file_direct` 使用 `http.serveFile()` 方法串流
    - `/get_file_stream` 调用 `ffmpeg` 并串流其标准输出 `stdout`
    - `/prepare_file_stream_direct` 调用 `ffmpeg` 预转码一个文件
    - `/get_file_stream_direct` 使用 `http.serveFile()` 获取预转码结束的临时文件
  - 其他在 `config.json` 中没有设定的 API 将默认拥有 `0` 的权限等级。

对于 Windows 用户，请确保 `ffmpeg` 正确安装并设置环境变量。

## 开发

欢迎任何 issue / pull request / feature request

### 主要变更历史

- `v1.0.0` 第一个版本。核心串流功能可用。
- `v1.1.0` 使用 `React` 重构前端。
- `v1.2.0` 数据库 DBMS 课程作业。添加 用户、标签、评论 和其他功能。

### ER Diagram

Database Entities Relationship Diagram

![ER Diagram](erdiagram.png)

- `avatar` 目前没有在使用。

- 第一次运行程序时，程序会自动创建一个 ID 为 `1` 的匿名用户。所有未登陆的用户都会自动登陆到这个账户。

- `tmpfs` 储存在内存中，每次重新启动后端程序将会清空记录的信息。

### 关于 tmpfs

如果前端的播放器勾选了 `预转码` 选项，后端程序会先将文件转码到临时目录中，转码完成后再串流文件。这么做可以实现断点续传，解决由于网络波动导致 `ffmpeg` 管道链接断开而终止转码的问题。

默认的临时文件夹目录是 `/tmp`，这是 Linux 系统中通用的临时目录。默认的生存时间是 600 秒（10 分钟）。如果超过这个时间没有访问该临时文件，那么后端程序将会自动删除它。

### 后端 API 设计

一个不需要返回任何有用数据的 API 将会返回下面的 JSON 对象

```json
{
  "status": "OK"
}
```

当错误发生时，后端会返回如下格式的 JSON 对象。`error` 是对错误信息的详细描述文本。

```json
{
  "error": "Wrong password"
}
```

不需要传递参数的 API 使用 `GET` 方法，否则使用 `POST` 方法。（忽略 RESTFUL 设计）

后端使用 cookies 来实现用户会话管理。任何不带 cookies 的请求会被认为是由 匿名用户 发送的（也就是 ID 为 `1` 的用户）

一些重要的源代码文件

- `pkg/api/api.go` 定义各个 API 的 URL 和对应函数。

- `pkg/database/sql_stmt.go` 定义 SQL 语句和做一些初始化工作。

- `pkg/database/struct.go` 定义 JSON 和 数据库对象 的 数据结构。
