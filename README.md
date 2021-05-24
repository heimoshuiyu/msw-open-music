# MSW Open Music Project

## 简介

Fork from `msw-file`，目前是一个音乐播放器。

## 编译 & 构建

### 编译后端

`go build`

如无任何输出，说明构建成功，可执行程序位于 `msw-open-music`

### 构建前端

`make`

说明：`Makefile` 脚本中的代码会在 `dist` 目录生成用于生产环境的前端 web 文件。这个脚本做的事情是简单地将 `vue` `vue-router` 等 js 文件替换成生产版本。并使用 `minify` 工具处理 `css` 和 `html` 文件。

## 使用

### 后端使用

初次使用请配置 `api_config.json`， **最重要的是配置 `token`** 。

#### api_config.json 说明

- `database_name` 字符串类型，指定 sqlite3 单文件数据库的位置，如果不存在则会自动创建。
- `addr` api 服务监听端口，该参数会被传入 `http.Serve.Addr`
- `token` 字符串，作为管理密码
- `ffmpeg_configs`，字典，其键是 ffmpeg 配置的名称，其值是放入 `ffmpeg -i input.mp3 -vn [此处] -f matroska -` 的参数，类型是字符串。 **注意：** 前端会按键名来排序配置列表，并以列表中的第一项作为默认配置。

### 前端使用

前端文件引用均使用相对路径，将前端文件放到同一目录下即可。

前端在调用后端 api 时使用的是绝对路径，例如 `/api/v1/hello`。如需更改，可以修改后端 `api.go` 中的 `apiMux` 和 `mux` 的相关属性。

## 后端 API 文档

说明中带有 `stream` 或 `流` 相关字样的，说明该 API 以 `io.Copy` 方式传输文件，不支持断点续传

无需返回数据的 API 将返回 OK，某些 API 可能会在 `status` 字段中返回详细的执行信息。

```json
{
    "status": "OK"
}
```



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

- `/web/*` 返回程序同目录下 web 文件夹中的内容

  此 api 仅用于方便开发，项目根目录中 web 文件夹中的内容并不是生产用（for production）的 js 文件，这个 API 不应该用来提供前端的 web 服务，web 服务应该由其他程序负责（例如 apache caddy nginx 等）