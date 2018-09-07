# gblog
Golang写的博客，移植hexo的next主题

## 安装部署

```
go get -u github.com/kardianos/govendor
git clone https://github.com/irellik/gblog.git
cd gblog
govendor sync
go build
./gblog
```