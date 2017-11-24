# ming800-to-redis

ming800-to-redis是一个用来把当前学期中课程，学生信息从明日系统同步到redis的程序。它是使用[Golang](https://golang.org)写的。

#### 提示
* ming800-to-redis会在同步数据前运行`FLUSHDB`来清空redis的数据库。
* 请确认当前的redis是用来同步明日系统数据。

#### 如何使用

1. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "server_url": "http://localhost:8080",
            "company": "my_company",
            "user": "Frank",
            "password": "my_password",
            "redis_server": "localhost:6379",
            "redis_password": ""
        }

* `"server_url"`是明日系统的URL地址。
* `"company"`是使用明日系统的组织或公司名称。
* `"user"`,`"password"`是明日系统的账号和密码。
* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。

2. 运行`ming800-to-redis`

        ./ming800-to-redis
