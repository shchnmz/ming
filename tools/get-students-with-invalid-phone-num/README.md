# get-students-with-invalid-phone-num

get-students-with-invalid-phone-num是一个输出明日系统无效联系电话的程序。它是使用[Golang](https://golang.org)写的。

#### 如何使用

1. 在和可执行文件相同的文件夹下，创建一个`config.json`的配置文件：

        {
            "server_url": "http://localhost:8080",
            "company": "my_company",
            "user": "Frank",
            "password": "my_password",
        }

* `"server_url"`是明日系统的URL地址。
* `"company"`是使用明日系统的组织或公司名称。
* `"user"`,`"password"`是明日系统的账号和密码。

2. 运行`get-students-with-invalid-phone-num`

        ./get-students-with-invalid-phone-num
