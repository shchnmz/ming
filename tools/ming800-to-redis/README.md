# ming800-to-redis

ming800-to-redis是一个用来把当前学期中课程，学生信息从明日系统同步到redis的程序。它是使用[Golang](https://golang.org)写的。

#### 同步后的redis中的keys
* 所有校区
  key: `"ming:campuses"`, type: ordered set, value: 校区, score: timestamp.

* 校区对应的课程
  key: `ming:$CAMPUS:categories`, type: ordered set, value: 课程, score: timestamp.

* 课程对应的校区
  key: `ming:$CATEGORY:classes`, type: ordered set, value: 校区, score: timestamp.

* 所有教师
  key: `"ming:teachers"`, type: ordered set, value: 教师, score: timestamp.

* 班级对应的教师
  key: `ming:$CAMPUS:$CATEGORY:$CLASS:teachers`, type: ordered set, value: 校区, score: timestamp.

* 教师对应的班级
  key: `ming:$TEACHER:classes`, type: ordered set, value: 班级, score: timestamp.

* 班级的上课时间段
  key: `ming:$CAMPUS:$CATEGORY:$CLASS:period`, type: string, value: 上课时间段(如果多个，只取第一个).

* 课程对应的所有时间段
  key: `ming:$CAMPUS:$CATEGORY:periods`, type: ordered set, value: 上课时间段, score: 上课时间段的权重.
  权重 = `周几*86400 + 开始小时 * 3600 + 开始分钟 * 60`

* 所有学生
  key: `ming:students`, type: ordered set, value: `$NAME:$PHONE_NUM`.

* 一个学生所在的班级
  key: `ming:$NAME:$PHONE_NUM:classes`, type: ordered set, value: `$CAMPUS:$CATEGORY:$CLASS`.

* 所有电话
  key: `ming:phones`, type: ordered set, value: 联系电话.

* 联系电话对应的学生
  key: `ming:$PHONE_NUM:students`, type: ordered set, value: 学生.

* 一个班级中所有学生
  key: `ming:$CAMPUS:$CATEGORY:$CLASS:students`, type: ordered set, value: 学生姓名.

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
