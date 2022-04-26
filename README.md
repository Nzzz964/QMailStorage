# QMailStorage.

> 使用 `QQ邮箱` 备份文件

## 😅 目前功能

- 上传指定文件

## 🙏 Todo

- [x] 分段存储文件
  - [X] I/O 性能调优

## 🙌 使用

1. 下载对应平台的 `Release`
2. 修改配置文件

```json
{
    "server"   : "smtp.qq.com:465",     // SMTP 服务器
    "username" : "admin@example.com",   // 登录账号
    "password" : "password",            // 登录密码
    "from"     : "admin@example.com",   // 发件人
    "to"       : "admin@example.com",   // 收件人
    "chunksize": 52428800               // 分片大小 (50M
}
```

3. 运行

```cmd
qmailstorage-windows-amd64.exe -h     // 获取帮助

Usage of qmailstorage-windows-amd64.exe:
  -c string
        指定配置文件 (default "config.json")
  -d string
        文件描述
  -f string
        指定上传文件
```

```cmd
qmailstorage-windows-amd64.exe -c config.json -d "描述" -f kawaii.zip
```

## 👌 文件还原

### Windows

```cmd
copy /b kawaii.zip.part1 + kawaii.zip.part2 + kawaii.zip.part3 kawaii.zip
```

### Linux

```bash
cat kawaii.zip.part2 kawaii.zip.part3 >> kawaii.zip.part1
mv  kawaii.zip.part1 kawaii.zip
```
