# get-phones-which-have-2-or-more-students

get-phones-which-have-2-or-more-students是一个输出明日系统中一个联系电话对应2个以上学生的所有联系电话的程序。它是使用[Golang](https://golang.org)写的。

#### 如何使用
1. 确认已经运行过[ming800-to-redis](../ming800-to-redis)将明日系统的数据导入到redis中。

2. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "redis_server": "localhost:6379",
            "redis_password": ""
        }

* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。

3. 运行`get-phones-which-have-2-or-more-students`

        ./get-phones-which-have-2-or-more-students
