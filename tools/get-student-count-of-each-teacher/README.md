# get-student-count-of-each-teacher

get-student-count-of-each-teacher是一个输出明日系统中每个教师对应的学生数量的程序。它是使用[Golang](https://golang.org)写的。

#### 如何使用
1. 确认已经运行过[ming800-to-redis](../ming800-to-redis)将明日系统的数据导入到redis中。

2. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "redis_server": "localhost:6379",
            "redis_password": ""
        }

* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。

3. 运行`get-student-count-of-each-teacher`

        ./get-student-count-of-each-teacher
