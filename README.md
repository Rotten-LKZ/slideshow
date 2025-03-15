# Slideshow

旨在用网页现实大量图片的轮播效果，节约人力，易于调节展示形式以及效果。

## 项目文件介绍

```
.
│  .gitignore
│  config.yaml 配置文件
│  index.html  轮播图展示网页
│  LICENSE     本仓库代码以 MIT 协议共享
│  README.md   README
│  style.css   样式
│
├─draft
│      _test.css
│      _test.html 部分测试代码，尝试在轮播图旁边列出当前播放人的所有图片，但是逻辑较为复杂暂时放弃
|
├─server
│      go.mod
│      go.sum
│      main.go 服务端主代码
|
```

## 配置

配置文件示例如下：
···yaml
basic:
    images_path: ''         # 图片存放路径 推荐使用绝对路径
    url: 'http://127.0.0.1' # 服务器地址，请不要加端口号
    port: 5000              # 服务器启用端口

anime:
    anime_duration: 500     # 动画持续时间（单位：毫秒）
    duration: 5000          # 图片持续时间（单位：毫秒）
    background_image: './bg.jpg' # 网页背景图片路径
```

如果路径希望使用相对路径，请注意相对路径地址是在运行 server.exe 的命令行当前所在目录

---

请注意，如果您的服务端和网页想要运行在不同地址，在 `index.html` 下第 16 行 `serverUrl` 变量改为对应服务器地址，比如：
```javascript
const serverUrl = 'http://127.0.0.1:5000' // 请自行加上 http(s) 以及结尾不接斜杠
```
