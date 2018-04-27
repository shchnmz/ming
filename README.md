# ming

[![Build Status](https://travis-ci.org/shchnmz/ming.svg?branch=master)](https://travis-ci.org/shchnmz/ming)
[![Go Report Card](https://goreportcard.com/badge/github.com/shchnmz/ming)](https://goreportcard.com/report/github.com/shchnmz/ming)
[![GoDoc](https://godoc.org/github.com/shchnmz/ming?status.svg)](https://godoc.org/github.com/shchnmz/ming)

ming是一个[Golang](https://golang.org)包，主要提供了将明日系统的数据导入至redis的功能。

#### 限制
* ming基于[ming800](https://github.com/northbright/ming800)
* 适合单机版本且只有1个校区的版本

#### 校区约定
因为受限1个校区，所以多个校区的可以通过在"课程类别"添加校区信息的做法来实现多个校区

* 课程类别命名约定

          课程类别（校区）

* 例子

          // 括号为中文全角括号
          四年级（校区A）
          校区：校区A
          类别：四年级

#### 同步后的redis中的keys
* 所有校区
  key: `"ming:campuses"`, type: ordered set, value: 校区, score: timestamp.

* 所有课程
  key: `"ming:categories"`, type: ordered set, value: 课程, score: timestamp.

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

* 时间段对应的班级
  key: `ming:$CAMPUS:$CATEGORY:$PERIOD:classes`, type: ordered set, value: 班级, score: timestamp.

* 所有学生
  key: `ming:students`, type: ordered set, value: `$NAME:$PHONE_NUM`.

* 一个学生所在的班级
  key: `ming:$NAME:$PHONE_NUM:classes`, type: ordered set, value: `$CAMPUS:$CATEGORY:$CLASS`.

* 所有电话
  key: `ming:phone_nums"`, type: ordered set, value: 联系电话.

* 联系电话对应的学生
  key: `ming:$PHONE_NUM:students`, type: ordered set, value: 学生.

* 一个班级中所有学生
  key: `ming:$CAMPUS:$CATEGORY:$CLASS:students`, type: ordered set, value: 学生姓名.

#### 例子

将明日系统数据导入到Redis中.

        // Create a ming.DB instance.
        db := ming.DB{RedisServer: RedisServer, RedisPassword: RedisPassword}

        // Sync Redis from ming server.
        // ServerURL is ming800 server URL.
        if err = db.SyncFromMing(ServerURL, Company, User, Password); err != nil {
                return
        } 

可以参考[tools/ming800-to-redis/](./tools/ming800-to-redis)

#### Documentation
* [API Reference](https://godoc.org/github.com/shchnmz/ming)

#### License
* [MIT License](LICENSE)
