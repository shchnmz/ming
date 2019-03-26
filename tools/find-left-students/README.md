# find-left-students

find-left-students是一个根据输入的CSV中的学生姓名，手机号码，输出明日系统中所有学生与CSV中的学生的差集的程序。它是使用[Golang](https://golang.org)写的。

#### 如何使用
1. 确认已经运行过[ming800-to-redis](../ming800-to-redis)将明日系统的数据导入到redis中。

2. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "redis_server": "localhost:6379",
            "redis_password": "",
            "input_csv_name_column_index": 6,
            "input_csv_phone_num_column_index": 7
        }

* `"redis_server"`,`"redis_password"`是同步的redis的地址和密码。
* `"input_csv_name_column_index"`是输入的CSV中的学生姓名的列的index(0-based)
* `"input_csv_phone_num_column_index"`是输入的CSV中的学生手机号码列的index(0-based)
* 输入CSV的文件需要放在可执行文件的相同文件夹下，命名为`input.csv`

3. 运行`find-left-students`

        ./find-left-students

4. 导出文件
   * CSV(使用`,`分隔，有UTF-8 BOM): `left-students.csv`
