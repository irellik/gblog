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

使用Nginx做服务器，反向代理Golang，详情页路由 /1.html使用Nginx rewrite 到 /post/1
