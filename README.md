# MSW Open Music Project

[![CI](https://github.com/heimoshuiyu/msw-open-music/actions/workflows/build.yml/badge.svg)](https://github.com/heimoshuiyu/msw-open-music/actions/workflows/build.yml)

## Introduction

A light weight personal music streaming platform.

![demo1](demo1.jpg)

[toc]

## How to build

### Build the back-end server

`make linux` or `make windows`

The executable file is named `msw-open-music` or `msw-open-music.exe`

### Build the font-end web pages

To build production web page `make web`

This command will go into `web` directory and install `node_modules`. Then execute `npm run build` command. The built web pages is under `web/build` directory.

To start the development, run `cd web` and `npm start`

## Usage

Start back-end server. Server will listen on 8080 port.

Build the font-end web page, then go to <http://127.0.0.1:8080>

By default:

- URL matched `/api/*` will process by back-end server.
- Others URL matched `/*` will be served files under `web/build/`

### Run back-end server

Configuration file is  `config.json`， **Please modify your `token`** 。

Default `ffmpeg_threads` is 1. Seems value larger than 1 will not increase the audio encode speed.

#### config.json description

- `database_name` string type. The filename of `sqlite3` database. Will create if that file doesn't exist.
- `addr` string type. The listen address and port.
- `token` string type. Password.
- `ffmpeg_config_list` list type, include `ffmpegConfig` object.
- `file_life_time` integer type (second). Life time for temporary file. If the temporary file is not accessed for more than this time, back-end server will delete this file.
- `cleaner_internal` integer type (second). Interval for `tmpfs` checking temporary file.
- `root` string type. Directory to store temporary files. Default is `/tmp`, **please modify this directory if you are using Windows.**

### Run font-end web page

Open your web browser to <http://127.0.0.1:8080> you will see the web pages.

## About tmpfs

If the `Prepare` mode is enabled in the font-wed player, back-end server will convert the whole file into the temporary folder, then serve file using native method. This can avoid ffmpeg pipe break problem cause by unstable network connection while streaming audio.

The default temporary folder is `/tmp`, which is a `tmpfs` file system in Linux operating system. Default life time for temporary files is 600 seconds (10 minutes). If the temporary file is not accessed for more than this time, back-end server will delete this file.

## Change log

- `v1.0.0` First version. Ready to use in production environment.
- `v1.1.0` Use `React` to rewrite the font-end web pages (Previous using `Vue`).

## Back-end API references

API named `stream` means it transfer data using `io.Copy`, which **DO NOT** support continue getting a partially-downloaded audio.

API does not need to respond any data will return the following JSON object.

```json
{
    "status": "OK"
}
```

### Anonymous API

Anonymous API can be called by anonymous.

- `/api/v1/hello` Just for test purpose.

- `/api/v1/get_file` Get a file with `stream` mode.

  - Request example

    ```json
    {
        "id": 123
    }
    ```

- `/api/v1/get_file_direct` Get a file with standart `http` methods, implement by `http.ServeFile` method.

  - Request example

    `/api/v1/get_file_direct?id=30`

- `/api/v1/search_files` Search files by filename.

  - Request example

    ```json
    {
        "filename": "miku",
        "limit": 10,
        "offset" 0
    }
    ```

    Search all files' name like `%miku%`. `%` is the wildcard in SQL. For example, `"filename": "miku%hatsune"` can match `hatsune miku`.

    `limit` Numbers of files in the respond. Should be within 1 - 10;

    `offset` It is the offset of the result, related to the page turning function.

  - Respond example

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
                "id": 31,
                "folder_id": 100,
                "folder_name": "wonderful",
                "filename": "memories (instrunment).flac",
                "filesize": 1248531
            }
        ]
    }
    ```

    `id` Identification of file.

    `folder_id` Identification of folder.

    `foldername` Folder name where the file in.

    `filename` File name.

    `filesize` File size, unit is byte.

- `/api/v1/search_folders` Search folders.

  - Request example.

    ```json
    {
        "foldername": "miku",
        "limit": 10,
        "offset": 0,
    }
    ```

    Search all folders' name like `%miku%`. `%` is the wildcard in SQL. For example, `"filename": "miku%hatsune"` can match `hatsune miku`.

    `limit` Numbers of files in the respond. Should be within 1 - 10;

    `offset` It is the offset of the result, related to the page turning function.

  - Respond example

    ```json
    {
        "folders": [
            {
                "id": 100,
                "foldername": "folder name"
            },
            {
                "id": 101,
                "foldername": "folder name (another)"
            }
        ]
    }
    ```

    `id` Identification of folder.

    `foldername` Folder name.

- `/api/v1/get_files_in_folder` Get files in a specify folder.

  - Request example.

    ```json
    {
        "folder_id": 123,
        "limit": 10,
        "offset": 0
    }
    ```

  - Respond example.

    Same with `/api/v1/search_files`

- `/api/v1/get_random_files` Randomly get 10 files.

  - Request example.

    GET `/api/v1/get_random_files`

  - Respond example.

    Same with `/api/v1/search_files`

- `/api/v1/get_file_stream`

  Stream file with a ffmpeg config name.

  - Request example.

    GET `/api/v1/get_file_stream?id=123&config=OPUS%20128k`

- `/api/v1/get_ffmpeg_config_list`

  Get ffmpeg config list

  - Request example

    GET `/api/v1/get_ffmpeg_config_list`

  - Respond example

    ```json
    {
        "ffmpeg_config_list": [
            {"name": "OPUS 256k", "args": "-c:a libopus -ab 256k"},
            {"name": "WAV", "args": "-c:a wav"}
        ]
    }
    ```

- `/api/v1/feedback` Send a feedback.

  - Request example

    ```json
    {
        "feedback": "some suggestions..."
    }
    ```

  - Respond OK.

- `/api/v1/get_file_info` Get information of a specify file.

  - Request example.

    ```json
    {
        "ID": 123
    }
    ```

  - Respond example.

    ```json
    {
        "id": 30,
        "folder_id": 100,
        "folder_name": "wonderful",
        "filename": "memories.flac",
        "filesize": 1048576
    },
    ```

- `/api/v1/get_file_stream_direct` Get a ffmpeg converted file with native http method. This API support continue getting a partially-downloaded audio. Note, you should call `/api/v1/prepare_file_stream_direct` first and wait for its respond, then call this API.

  - Request example

    GET `/api/v1/get_file_stream_direct?id=123&config=OPUS%20128k`

- `/api/v1/prepare_file_stream_direct` Ask server to convert a file with specific ffmpeg config name. When the conver process is finished, server will reply with the converted file size.

  - Request example

    ```json
    {
        "id": 123,
        "config_name": "OPUS 128k"
    }
    ```

  - Respond example

    ```json
    {
        "filesize": 1973241
    }
    ```

### API needs token

- `/api/v1/walk` Walk directory, add all files and folders to database.

  - Request example

    ```json
    {
        "token": "your token",
        "root": "/path/to/root",
        "pattern": [".wav", ".flac"]
    }
    ```

    `token` The token in `config.json` file.

    `root` Root directory server will walk throught

    `pattern` A list of pattern that files ends with. Only files matched a pattern in list will be add to database.

  - Respond OK

- `/api/v1/reset` Rest the **files and folders table**

  - Request example

    ```json
    {
        "token": "your token"
    }
    ```

  - Respond OK

- `/api/v1/add_ffmpeg_config` Add ffmpeg config.

  Will be changed in future.

  - Request example

    ```json
    {
        "token": "your token",
        "name": "OPUS",
        "ffmpeg_config": {
            "args": "-c:a libopus -ab 256k"
        }
    }
    ```

    `name` Name of the ffmpeg config.

    `ffmpeg_config`

    `args`

  - Respond OK

## Font-end API references

Currently only few APIs in font-end.

- `/#/files/39/share`

  Share a specific file.

- `/#/folders/2614`

  Show files in a specific folder.
