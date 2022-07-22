# MSW Open Music Project

[![CI](https://github.com/heimoshuiyu/msw-open-music/actions/workflows/build.yml/badge.svg)](https://github.com/heimoshuiyu/msw-open-music/actions/workflows/build.yml)

> The best way to search for a music is to load up a huge playlist and shuffle until you find it.

A 💪 light weight ⚡️ blazingly fast 🖥️ cross platform personal music streaming platform. Manage your existing music files and enjoy them on any devices.

Front-end web application build with `react.js` and `water.css`, back-end build with `golang` and `sqlite`.

## Introduction

Screenshot

![demo1](demo1.jpg)

### Features

- 🔎 Index your existing music files, and record file name and folder information.

- 📕 Use folder 📁 tag 🏷️ review 💬 to manage your music.

- 🌐 Provide a light weight web application with multi-language support.

- 👥 Multi-user support.

- 🔥 Call `ffmpeg` with customizable preset to stream your music.

- 🔗 Share music with others!

### Try it if you...

- Already saved a lot of music files on disk. 🖴

- Downloaded tons of huge lossless music. 🎵

- Wants to stream your music files from PC/Server to PC/phone. 😋

- Wants to share your stored music. 😘

## Usage

1. Modify the `secret` in `config.json`

2. Run back-end server `msw-open-music.exe` or `msw-open-music`. Server will listen on 8080 port by default. Then open <http://127.0.0.1:8080> to setup first admin account.

The front-end HTML files are under `web/build`

### Setup first admin account

The first administrator account will be active automatically, other administrator accounts need active manually.

Go to register page, select the role to admin, and register the first admin account.

#### config.json

- `secret` string type. Secret to encrypt the session.

- `database_name` string type. The filename of `sqlite3` database. Will create if that file doesn't exist.
- `addr` string type. The listen address and port.
- `ffmpeg_config_list` list type, include `ffmpegConfig` object.
- `file_life_time` integer type (second). Life time for temporary file. If the temporary file is not accessed for more than this time, back-end server will delete this file.
- `cleaner_internal` integer type (second). Interval for `tmpfs` checking temporary file.
- `root` string type. Directory to store temporary files. Default is `/tmp`, **please modify this directory if you are using Windows.** Directory will be created if not exists.

For windows user, make sure you have `ffmpeg` installed.

## Development

Any issues or pull requests are welcome.

### Major changes log

- `v1.0.0` First version. Implement the core streaming function.
- `v1.1.0` Use `React` to rewrite the font-end web pages.
- `v1.2.0` Add user, tag, review and other functions for DBMS course project.

### ER Diagram

Database Entities Relationship Diagram

![ER Diagram](erdiagram.png)

- `avatar` is not using currently

- The first time you run the program, the server will create an anonymous user with id `1`. All users who are not logged in will be automatically logged in to this account.

- `tmpfs` is store in memory, which will be empty everytime server restart.

### About tmpfs

If the `Prepare` mode is enabled in the font-wed player, back-end server will convert the whole file into the temporary folder, then serve file. This can avoid `ffmpeg` pipe break problem cause by unstable network connection while streaming audio.

The default temporary folder is `/tmp`, which is a `tmpfs` file system in Linux operating system. Default life time for temporary files is 600 seconds (10 minutes). If the temporary file is not accessed for more than this time, back-end server will delete this file.

### Back-end API design

API does not need to respond any data will return the following JSON object.

```json
{
    "status": "OK"
}
```

Sometime errors happen, server will return the following JSON object, which `error` is the detailed error message.

```json
{
    "error": "Wrong password"
}
```

API does not need to send any data should use `GET` method, otherwise use `POST` method.

Server use cookies to authenticate a user. Any request without cookies will be consider from an anonymous user (aka. user with ID `1`).

Some important source code files:

- `pkg/api/api.go` define URL

- `pkg/database/sql_stmt.go` define SQL queries and do the init job.

- `pkg/database/struct.go` define JSON structures for database entities.
