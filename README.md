# MSW Open Music Project

## 简介

Fork from `msw-file`，目前是一个音乐播放器。

![demo1](demo1.jpg)

[toc]

## 编译 & 构建

### 编译后端

`make linux` 或 `make windows`

如无任何输出，说明构建成功，可执行程序位于 `msw-open-music`

### 构建前端

`make web`

这条命令会在 `web` 目录终安装 `node_modules` 模块，并执行 `build` 脚本。

构建好的静态网页位于 `web/build` 中

## 使用

### 后端使用

初次使用请配置 `config.json`， **最重要的是配置 `token`** 。

默认 ffmpeg 线程 `ffmpeg_threads` 为 1 ，大于 1 以上的值似乎对编码音频没有效果。

#### config.json 说明

- `database_name` 字符串类型，指定 sqlite3 单文件数据库的位置，如果不存在则会自动创建。
- `addr` api 服务监听端口，该参数会被传入 `http.Serve.Addr`
- `token` 字符串，作为管理密码
- `ffmpeg_config_list` 列表，包含 `ffmpegConfig` 对象
- `file_life_time` 临时文件生存时间，超过该时间没有访问该临时文件，tmpfs 将删除此文件。
- `cleaner_internal` 清理器的检查间隔。
- `root` 存放该临时文件目录， **Windows 用户请替换成合适的目录。**

### 前端使用

前端文件引用均使用相对路径，将前端文件放到同一目录下即可。

## 关于临时文件夹的说明

前端播放器中勾选了 `Prepare` ，后端将转码文件到临时文件夹，然后直链提供文件。这有助于修复网路不稳定时 TCP 链接断开，stream 模式下 ffmpeg 中断输出并且不能断点续传的问题。

临时文件夹管理器位于 `internal/pkg/tmpfs` 中，默认删除时间是 10 分钟。10分钟内如果没有对该临时文件的访问，则会删除此临时文件。

## Change log

- `v1.0.0` 首个版本
- `v1.1.0` 使用 react 和 webpack 构建前端

## 后端 API 文档

说明中带有 `stream` 或 `流` 相关字样的，说明该 API 以 `io.Copy` 方式传输文件，不支持断点续传

无需返回数据的 API 将返回 OK，某些 API 可能会在 `status` 字段中返回详细的执行信息。

```json
{
    "status": "OK"
}
```

### 公开 API

- `/api/v1/hello` OK 测试

- `/api/v1/get_file` 以流方式获取文件

  - 请求示例

    ```json
    {
        "id": 123
    }
    ```

- `/api/v1/get_file_direct` http 标准方式获取文件，支持断点续传，由 `http.ServeFile` 实现

  - 请求示例

    `/api/v1/get_file_direct?id=30`

- `/api/v1/search_files` 搜索文件

  - 请求示例

    ```json
    {
        "filename": "miku",
        "limit": 10,
        "offset" 0
    }
    ```

    搜索所有文件名中包含 "miku" 的文件

    `limit` 限制返回结果的数量，该值必须在 0~10 之间

    `offset` 是返回结构的偏移量，用于实现翻页功能。

  - 返回示例

    ```json
    {
        "files": [
            {
                "id": 30,
                "folder_id": 100,
                "folder_name": "wonderful",
                "filename": "memories.flac",
                "filesize": 1048576
            },
            {
                "id": 30,
                "folder_id": 100,
                "folder_name": "wonderful",
                "filename": "memories.flac",
                "filesize": 1048576
            }
        ]
    }
    ```

    `id` 为文件的唯一标识

    `folder_id` 为该文件所在的文件夹标识

    `foldername` 为该文件所在的文件夹名

    `filename` 为该文件名

    `filesize` 为该文件的大小，单位字节

- `/api/v1/search_folders` 搜索文件夹

  - 请求示例

    ```json
    {
        "foldername": "miku",
        "limit": 10,
        "offset": 0,
    }
    ```

    搜索所有文件夹名中包含 "miku" 的文件夹。

    `limit` 限制返回结果的数量，该值必须在 0~10 之间

    `offset` 是返回结构的偏移量，用于实现翻页功能。

  - 返回示例

    ```json
    {
        "folders": [
            {
                "id": 100,
                "foldername": "folder name"
            },
            {
                "id": 100,
                "foldername": "folder name"
            }
        ]
    }
    ```

    `id` 为该文件夹的唯一标识

    `foldername` 为该文件夹的名字

- `/api/v1/get_files_in_folder` 获取指定文件夹中的所有文件

  - 请求示例

    ```json
    {
        "folder_id": 123,
        "limit": 10,
        "offset": 0
    }
    ```

  - 返回示例

    同 `/api/v1/search_files`

- `/api/v1/get_random_files`

  此 API 随机返回 files 表中 10 个文件。请注意，该操作会造成全表查询，在 AMD 2200G CPU 40000条数据记录情况下最大处理量为 100 请求每秒。

  - 请求示例

    直接 GET `/api/v1/get_random_files`

  - 返回示例

    同 `/api/v1/search_files`

- `/api/v1/get_file_stream`

  以流方式返回文件

  - 请求示例

    GET `/api/v1/get_file_stream?id=123`

- `/api/v1/get_ffmpeg_config_list`

  获取 ffmpeg 配置列表

  - 请求示例

    GET `/api/v1/get_ffmpeg_config_list`

  - 返回示例

    ```json
    {
        "ffmpeg_configs": {
            "OPUS 256k": {"args": "-c:a libopus -ab 256k"},
            "WAV": {"args": "-c:a wav"}
        }
    }
    ```

- `/api/v1/feedback` 反馈

  - 请求示例

    ```json
    {
        "feedback": "some suggestions..."
    }
    ```

  - 返回 OK

- `/api/v1/get_file_info` 获取单个文件的信息

  - 请求示例

    ```json
    {
        "ID": 123
    }
    ```

  - 返回示例

    ```json
    {
        "id": 30,
        "folder_id": 100,
        "folder_name": "wonderful",
        "filename": "memories.flac",
        "filesize": 1048576
    },
    ```

- `/api/v1/get_file_stream_direct` 获取已提前转码好的文件，该 API 支持断点续传

  - 请求示例

    GET `/api/v1/get_file_stream_direct?id=123&config=OPUS 128k`

- `/api/v1/prepare_file_stream_direct` 请求提前转码文件，该 API 将返回转码后的文件大小

  - 请求示例

    ```json
    {
        "id": 123,
        "config_name": "OPUS 128k"
    }
    ```

  - 返回示例

    ```json
    {
        "filesize": 1973241
    }
    ```

### 需要 token 的 API

- `/api/v1/walk` 遍历目录，并将文件和文件夹添加到数据库中

  - 请求示例

    ```json
    {
        "token": "your token",
        "root": "/path/to/root",
        "pattern": [".wav", ".flac"]
    }
    ```

    `token` 此 API 需要 token

    `root` 遍历目录

    `pattern` 文件扩展名列表（包含 `.` ），匹配扩展名的文件才会被添加到数据库

  - 返回 OK

- `/api/v1/reset` 重置数据库（feedbacks 不会清空）

  - 请求示例

    ```json
    {
        "token": "your token"
    }
    ```

  - 返回 OK

- `/api/v1/add_ffmpeg_config` 添加 ffmpeg 配置

  注意：目前前端中没有实现此功能

  - 请求示例

    ```json
    {
        "token": "your token",
        "name": "OPUS",
        "ffmpeg_config": {
            "args": "-c:a libopus -ab 256k"
        }
    }
    ```

    `name` 该配置的名字

    `ffmpeg_config` 一个 ffmpeg 的配置

    `args` 该 ffmpeg 配置的参数

  - 返回 OK

- `/*` 返回程序同目录下 `web/build` 文件夹中的内容

## 前端 API 文档

前端只有少量 API ，允许用户直接打开链接就执行某些功能

- `/#/share/39`

  分享文件

- `/#/search-folders/2614`

  显示该文件夹中的文件
