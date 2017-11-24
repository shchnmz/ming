package ming

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/northbright/ming800"
	"github.com/northbright/redishelper"
)

// Processor implements ming800.WalkProcessor interface.
type Processor struct {
	redisServer   string
	redisPassword string
}

// ParseCategory gets campus and real category from category string.
//
//   Param:
//       category: raw category string like this: 初一（中山）
//   Return:
//       campus, category. e.g. campus: 中山,category: 初一
func ParseCategory(category string) (string, string) {
	p := `^(\S+)（(\S+)）$`
	re := regexp.MustCompile(p)
	matched := re.FindStringSubmatch(category)
	if len(matched) != 3 {
		return "", ""
	}
	return matched[2], matched[1]
}

// GetPeriodScore gets the score for the period string.
//
// Params:
//     period: period string. e.g. "星期一09:00-11:30".
// Return:
//     computed score.
func GetPeriodScore(period string) int {
	dayScores := map[string]int{
		"一": 1,
		"二": 2,
		"三": 3,
		"四": 4,
		"五": 5,
		"六": 6,
		"日": 7,
	}

	p := `^星期(\S)(\d{2}):(\d{2})`
	re := regexp.MustCompile(p)
	matched := re.FindStringSubmatch(period)
	if len(matched) != 4 {
		return 0
	}

	day := matched[1]
	if _, ok := dayScores[day]; !ok {
		return 0
	}

	hour, _ := strconv.Atoi(matched[2])
	min, _ := strconv.Atoi(matched[3])

	return dayScores[day]*86400 + hour*3600 + min*60
}

// ClassHandler implements ming800.WalkProcessor interface.
// It'll be called when a class is found.
func (p *Processor) ClassHandler(class ming800.Class) {
	var err error

	defer func() {
		if err != nil {
			log.Printf("classHandler() error: %v", err)
		}
	}()

	pipedConn, err := redishelper.GetRedisConn(p.redisServer, p.redisPassword)
	if err != nil {
		return
	}
	defer pipedConn.Close()

	pipedConn.Do("MULTI")

	campus, category := ParseCategory(class.Category)
	if category == "" && campus == "" {
		err = fmt.Errorf("Failed to parse category and campus: %v", class.Category)
		return
	}

	// Get timestamp as score for redis ordered set.
	t := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Update SET: key: "campuses", value: campuses.
	k := "campuses"
	pipedConn.Send("ZADD", k, t, campus)

	// Update SET: key: campus, value: categories.
	k = fmt.Sprintf("%v:categories", campus)
	pipedConn.Send("ZADD", k, t, category)

	// Update SET: key: category, value: campuses.
	k = fmt.Sprintf("%v:campuses", category)
	pipedConn.Send("ZADD", k, t, campus)

	// Update SET: key: campus + category, value: classes.
	k = fmt.Sprintf("%v:%v:classes", campus, category)
	pipedConn.Send("ZADD", k, t, class.Name)

	for _, teacher := range class.Teachers {
		// Update SET: key: "teachers", value: teachers.
		k = "teachers"
		pipedConn.Send("ZADD", k, t, teacher)

		// Update SET: key: campus + category + class, value: teachers.
		k = fmt.Sprintf("%v:%v:%v:teachers", campus, category, class.Name)
		pipedConn.Send("ZADD", k, t, teacher)

		// Update SET: key: teacher, value: campus + category + class.
		k = fmt.Sprintf("%v:classes", teacher)
		v := fmt.Sprintf("%v:%v:%v", campus, category, class.Name)
		pipedConn.Send("ZADD", k, t, v)
	}

	if len(class.Periods) >= 1 {
		period := class.Periods[0]
		score := GetPeriodScore(period)

		// Update STRING: key: campus + category + class, value: period(1st period).
		k = fmt.Sprintf("%v:%v:%v:period", campus, category, class.Name)
		pipedConn.Send("SET", k, period)

		// Update SET: key: campus + category, value: periods.
		k = fmt.Sprintf("%v:%v:periods", campus, category)
		pipedConn.Send("ZADD", k, score, period)
	}

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return
	}
}

// StudentHandler implements ming800.WalkProcessor interface.
// It'll be called when a student is found.
func (p *Processor) StudentHandler(class ming800.Class, student ming800.Student) {
	var err error

	defer func() {
		if err != nil {
			log.Printf("studentHandler() error: %v", err)
		}
	}()

	// Check if phone number: 11-digit or 8-digit.
	if !ValidPhoneNum(student.PhoneNum) {
		fmt.Printf("%s,%s,%s,%s\n", class.Category, class.Name, student.Name, student.PhoneNum)
		return
	}

	// Student contact phone may have '.' suffix, remove it.
	student.PhoneNum = strings.TrimRight(student.PhoneNum, `.`)

	// Get another redis connection for pipelined transaction.
	pipedConn, err := redishelper.GetRedisConn(p.redisServer, p.redisPassword)
	if err != nil {
		return
	}
	defer pipedConn.Close()

	pipedConn.Do("MULTI")

	// Get timestamp as store for redis ordered set.
	t := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Get campus, category.
	campus, category := ParseCategory(class.Category)
	if category == "" && campus == "" {
		err = fmt.Errorf("Failed to parse category and campus: %v", class.Category)
		return
	}

	// Update SET: key: "students", value: student name + student phone num.
	k := "students"
	v := fmt.Sprintf("%v:%v", student.Name, student.PhoneNum)
	pipedConn.Send("ZADD", k, t, v)

	// Update SET: key: student name + student contact phone num, value: campus + category + class name.
	k = fmt.Sprintf("%v:%v:classes", student.Name, student.PhoneNum)
	v = fmt.Sprintf("%v:%v:%v", campus, category, class.Name)
	pipedConn.Send("ZADD", k, t, v)

	// Update SET: key: "phones", value: student contact phone num.
	k = "phones"
	pipedConn.Send("ZADD", k, t, student.PhoneNum)

	// Update SET: key: student contact phone num, value: student names.
	k = fmt.Sprintf("%v:students", student.PhoneNum)
	pipedConn.Send("ZADD", k, t, student.Name)

	// Update SET: key: campus + category + class, value: student name + student contact phone num.
	k = fmt.Sprintf("%v:%v:%v:students", campus, category, class.Name)
	v = fmt.Sprintf("%v:%v", student.Name, student.PhoneNum)
	pipedConn.Send("ZADD", k, t, v)

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return
	}
}

// Ming2Redis sync all current campuses, categories, students data from ming800 to redis.
// Warning: it'll do a "FLUSHDB" command to flush redis database before sync.
// Make sure this redis db is dedicated to sync data only.
//
// Params:
//     serverURL: server URL of ming800. e.g. "http://192.168.1.87:8080".
//     company: company or orgnization name of ming800.
//     user: user account of ming800.
//     password: user password of ming800.
//     redisServer: redis server address. e.g. ":6379".
//     redisPassword: redis server password.
func Ming2Redis(serverURL, company, user, password, redisServer, redisPassword string) error {
	// New a session
	s, err := ming800.NewSession(serverURL, company, user, password)
	if err != nil {
		return fmt.Errorf("NewSession() error: %v", err)
	}

	// Login
	if err = s.Login(); err != nil {
		return fmt.Errorf("Login() error: %v", err)
	}

	// Warnning: FLUSHDB before every sync.
	// Make sure this redis db is used to sync ming800 data only.
	conn, err := redishelper.GetRedisConn(redisServer, redisPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err = conn.Do("FLUSHDB"); err != nil {
		return err
	}

	// Walk
	// Write your own class and student handler functions.
	// Class and student handler will be called while walking ming800.
	processor := &Processor{redisServer: redisServer, redisPassword: redisPassword}
	if err = s.Walk(processor); err != nil {
		return fmt.Errorf("Walk() error: %v", err)
	}

	// Logout
	if err = s.Logout(); err != nil {
		return fmt.Errorf("Logout() error: %v", err)
	}

	return nil
}
