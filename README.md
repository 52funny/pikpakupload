# PIKPAKUPLOAD

Pikpak 的上传工具

> 首先将`config_example.yml`配置一下, 输入自己的账号密码
>
> 账号要以区号开头 如 `+861xxxxxxxxxx`
>
> 然后将其重命名为`config.yml`

## 使用方法

### 编译

首先你得拥有 go 的环境

[go install guide](https://go.dev/doc/install)

生成可执行文件

```bash
go build
```

### 执行

将本地目录下的所有文件上传至 `pikpak` 根目录/Movies

```bash
./pikpakupload -p Movies .
```

将本地目录下除了后缀名为`mp3`, `jpg`的文件上传至 `pikpak` 根目录/Movies

```bash
./pikpakupload -exn ".mp3$" -exn ".jpg" -p Movies .
```

指定上传的协程数目(默认为16)
```bash
./pikpakupload -c 20 -p Movies .
```