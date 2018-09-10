# niu admin

* based on [gopub](https://github.com/lisijie/gopub)

安装步骤：

1. 创建数据库，将install.mgo导入mongodb。
2. 修改 conf/app.conf 中相关的配置。
3. 使用命令 `./service.sh start` 启动，如果无法启动，检查主程序 gopub 是否具有可执行权限，使用 `chmod +x ./gopub` 增加权限。
4. 使用 `http://localhost:8000` 访问。
5. 后台默认帐号为 `admin`，密码为 `admin888`。 
