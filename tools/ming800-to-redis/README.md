# ming800-to-redis

ming800-to-redis是一个用来把当前学期中课程，学生信息从明日系统同步到redis的程序。它是使用[Golang](https://golang.org)写的。

#### 提示
* ming800-to-redis会在同步数据前运行`FLUSHDB`来清空redis的数据库。
* 请确认当前的redis是用来同步明日系统数据。

#### 同步后的redis中的keys
* 所有校区
  key: `"campuses"`, type: ordered set, value: 校区.

* 校区对应的课程
  key: `$CAMPUS:categories`, type: ordered set, value: 课程.

* 课程对应的校区
  key: `$CATEGORY:classes`, type: ordered set, value: 校区.

* 班级对应的教师
  key: `$CAMPUS:$CATEGORY:$CLASS:teachers`, type: ordered set, value: 校区.

* 班级的上课时间段
  key: `$CAMPUS:$CATEGORY:$CLASS:period`, type: string, value: 上课时间段(如果多个，只取第一个).

* 课程对应的所有时间段
  key: `$CAMPUS:$CATEGORY:periods`, type: ordered set, value: 上课时间段.

* 所有学生
  key: `students`, type: ordered set, value: `$NAME:$PHONE_NUM`.

* 一个学生所在的班级
  key: `$NAME:$PHONE_NUM:classes`, type: ordered set, value: `$CAMPUS:$CATEGORY:$CLASS`.

* 所有电话
  key: `phones`, type: ordered set, value: 联系电话.

* 联系电话对应的学生
  key: `$PHONE_NUM:students`, type: ordered set, value: 学生.

* 一个班级中所有学生
  key: `$CAMPUS:$CATEGORY:$CLASS:students`, type: ordered set, value: 学生姓名.

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
