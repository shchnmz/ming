# get-students-with-invalid-phone-num

get-students-with-invalid-phone-num是一个用来输出明日系统中无效学生联系电话的程序。它是使用[Golang](https://golang.org)写的。

#### 联系电话的正确格式
* 对于8位固定电话
  * 以8位数字开始
  * 可以有0个或者多个`.`(英文半角)作为后缀
    例如：`33001100` `33001100.` `33001100..`

* 对于11位手机
  * 以11位数字开始
  * 可以有0个或者多个`.`(英文半角)作为后缀
    例如：`13800138000` `13800138000.` `13800138000..`

#### 关于存在`.`后缀的原因
* 参考[修改明日系统中已存在的有问题的学生联系电话](https://github.com/shchnmz/worklog/blob/master/software/doc/edit-existing-contact-phone-num-in-ming800.md)

#### 输出格式
* `年级,班级,学生姓名,电话号码`

#### 如何使用

1. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "server_url": "http://localhost:8080",
            "company": "my_company",
            "user": "Frank",
            "password": "my_password"
        }

* `"server_url"`是明日系统的URL地址。
* `"company"`是使用明日系统的组织或公司名称。
* `"user"`,`"password"`是明日系统的账号和密码。

2. 运行`get-students-with-invalid-phone-num`

        // 直接输出:
        ./get-students-with-invalid-phone-num

        // 导出成CSV格式:
        ./get-students-with-invalid-phone-num > ~/1.csv

3. 在Windows下打开CSV文件
    * 在Windows下先用记事本打开，另存为一个文件时选择编码为`UTF-8`
    * 然后使用Excel直接带开
