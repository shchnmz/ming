package ming

import (
	"fmt"
	//"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/northbright/ming800"
	"github.com/northbright/redishelper"
)

// DB sync data from ming system.
type DB struct {
	RedisServer   string
	RedisPassword string
}

// ParseCategory gets campus and real category from category string.
//
//	Param:
//	    category: raw category string like this: 初一（中山）
//	Return:
//	    campus, category. e.g. campus: 中山,category: 初一
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
//
//	period: period string. e.g. "星期一09:00-11:30".
//
// Return:
//
//	computed score.
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

// ClassHandler implements ming800.WalkDB interface.
// It'll be called when a class is found.
func (db *DB) ClassHandler(class *ming800.Class) error {
	pipedConn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return fmt.Errorf("Dial() error: %v", err)
	}
	defer pipedConn.Close()

	pipedConn.Do("MULTI")

	campus, category := ParseCategory(class.Category)
	if category == "" && campus == "" {
		return fmt.Errorf("Failed to parse category and campus: %v", class.Category)
	}

	// Get timestamp as score for redis ordered set.
	t := strconv.FormatInt(time.Now().UnixNano(), 10)

	k := "ming:campuses"
	pipedConn.Send("ZADD", k, t, campus)

	k = "ming:categories"
	pipedConn.Send("ZADD", k, t, category)

	k = fmt.Sprintf("ming:%v:categories", campus)
	pipedConn.Send("ZADD", k, t, category)

	k = fmt.Sprintf("ming:%v:campuses", category)
	pipedConn.Send("ZADD", k, t, campus)

	k = fmt.Sprintf("ming:%v:%v:classes", campus, category)
	pipedConn.Send("ZADD", k, t, class.Name)

	for _, teacher := range class.Teachers {
		k = "ming:teachers"
		pipedConn.Send("ZADD", k, t, teacher)

		k = fmt.Sprintf("ming:%v:%v:%v:teachers", campus, category, class.Name)
		pipedConn.Send("ZADD", k, t, teacher)

		k = fmt.Sprintf("ming:%v:classes", teacher)
		v := fmt.Sprintf("%v:%v:%v", campus, category, class.Name)
		pipedConn.Send("ZADD", k, t, v)
	}

	if len(class.Periods) >= 1 {
		period := class.Periods[0]
		score := GetPeriodScore(period)

		k = fmt.Sprintf("ming:%v:%v:%v:period", campus, category, class.Name)
		pipedConn.Send("SET", k, period)

		k = fmt.Sprintf("ming:%v:%v:periods", campus, category)
		pipedConn.Send("ZADD", k, score, period)

		k = fmt.Sprintf("ming:%v:%v:%v:classes", campus, category, period)
		pipedConn.Send("ZADD", k, score, class.Name)
	}

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

// StudentHandler implements ming800.WalkDB interface.
// It'll be called when a student is found.
func (db *DB) StudentHandler(class *ming800.Class, student *ming800.Student) error {
	// Check if phone number: 11-digit or 8-digit.
	if !ValidPhoneNum(student.PhoneNum) {
		fmt.Printf("Invalid Phone Num: %s,%s,%s,%s\n", class.Category, class.Name, student.Name, student.PhoneNum)
		return nil
	}

	// Student contact phone may have '.' suffix, remove it.
	student.PhoneNum = strings.TrimRight(student.PhoneNum, `.`)

	// Get another redis connection for pipelined transaction.
	pipedConn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return err
	}
	defer pipedConn.Close()

	pipedConn.Do("MULTI")

	// Get timestamp as store for redis ordered set.
	t := strconv.FormatInt(time.Now().UnixNano(), 10)

	// Get campus, category.
	campus, category := ParseCategory(class.Category)
	if category == "" && campus == "" {
		return fmt.Errorf("Failed to parse category and campus: %v", class.Category)
	}

	k := "ming:students"
	v := fmt.Sprintf("%v:%v", student.Name, student.PhoneNum)
	pipedConn.Send("ZADD", k, t, v)

	k = fmt.Sprintf("ming:%v:%v:classes", student.Name, student.PhoneNum)
	v = fmt.Sprintf("%v:%v:%v", campus, category, class.Name)
	pipedConn.Send("ZADD", k, t, v)

	k = "ming:phone_nums"
	pipedConn.Send("ZADD", k, t, student.PhoneNum)

	k = fmt.Sprintf("ming:%v:students", student.PhoneNum)
	pipedConn.Send("ZADD", k, t, student.Name)

	k = fmt.Sprintf("ming:%v:%v:%v:students", campus, category, class.Name)
	v = fmt.Sprintf("%v:%v", student.Name, student.PhoneNum)
	pipedConn.Send("ZADD", k, t, v)

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

// SyncFromMing sync data included all current campuses, categories, students from ming800 to redis.
//
// Params:
//
//	serverURL: server URL of ming800. e.g. "http://192.168.1.87:8080".
//	company: company or orgnization name of ming800.
//	user: user account of ming800.
//	password: user password of ming800.
func (db *DB) SyncFromMing(serverURL, company, user, password string) error {
	// New a session
	s, err := ming800.NewSession(serverURL, company, user, password)
	if err != nil {
		return fmt.Errorf("NewSession() error: %v", err)
	}

	// Login
	if err = s.Login(); err != nil {
		return fmt.Errorf("Login() error: %v", err)
	}

	// Clear all data before sync.
	if err = db.Clear(); err != nil {
		return err
	}

	// Walk
	// Write your own class and student handler functions.
	// Class and student handler will be called while walking ming800.
	if err = s.Walk(db); err != nil {
		return fmt.Errorf("Walk() error: %v", err)
	}

	// Logout
	if err = s.Logout(); err != nil {
		return fmt.Errorf("Logout() error: %v", err)
	}

	return nil
}

// Clear cleans all existing ming800 data in redis.
// Do it before each new sync.
func (db *DB) Clear() error {
	var (
		err   error
		v     []interface{}
		items []string
	)

	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return err
	}
	defer conn.Close()

	pipedConn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return err
	}
	defer pipedConn.Close()

	pipedConn.Send("MULTI")

	pattern := `ming:*`
	cursor := 0
	for {
		if v, err = redis.Values(conn.Do("SCAN", cursor, "MATCH", pattern, "COUNT", 1000)); err != nil {
			return err
		}

		if _, err = redis.Scan(v, &cursor, &items); err != nil {
			return err
		}

		for _, key := range items {
			pipedConn.Send("DEL", key)
		}

		if cursor == 0 {
			break
		}
	}

	if _, err = pipedConn.Do("EXEC"); err != nil {
		return err
	}

	return nil
}

// GetNamesByPhoneNum searches student names by phone number.
func (db *DB) GetNamesByPhoneNum(phoneNum string) ([]string, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:students", phoneNum)
	names, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return []string{}, err
	}

	return names, nil
}

// GetClassesByNameAndPhoneNum searches classes by student name and phone number.
func (db *DB) GetClassesByNameAndPhoneNum(name, phoneNum string) ([]string, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:%v:classes", name, phoneNum)
	classes, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return []string{}, err
	}

	return classes, nil
}

// ParseClassValue parses the value of class.
//
// Params:
//
//	classValue: class string contains campus, category and real class.
//	            format: $CAMPUS:$CATEGORY:$CLASS e.g. "新校区:一年级:一年级2班"
//
// Returns:
//
//	campus, category, real class.
func ParseClassValue(classValue string) (string, string, string) {
	arr := strings.SplitN(classValue, ":", 3)
	campus := arr[0]
	category := arr[1]
	class := arr[2]

	return campus, category, class
}

// GetClassPeriod gets the period of the combination of campus, category, class.
func (db *DB) GetClassPeriod(campus, category, class string) (string, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:%v:%v:period", campus, category, class)
	period, err := redis.String(conn.Do("GET", k))
	if err != nil && err != redis.ErrNil {
		return "", err
	}
	return period, nil
}

// ValidClass validates if the campus, category, class info match.
func (db *DB) ValidClass(campus, category, class string) (bool, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:%v:classes", campus, category)
	score, err := redis.String(conn.Do("ZSCORE", k, class))
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if score == "" {
		return false, nil
	}

	return true, nil
}

// ValidPeriod validates if the campus, category, period info match.
func (db *DB) ValidPeriod(campus, category, period string) (bool, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return false, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:%v:periods", campus, category)
	score, err := redis.String(conn.Do("ZSCORE", k, period))
	if err != nil && err != redis.ErrNil {
		return false, err
	}

	if score == "" {
		return false, nil
	}

	return true, nil
}

// GetTeachersOfClass gets teachers of the class.
func (db *DB) GetTeachersOfClass(campus, category, class string) ([]string, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:%v:%v:teachers", campus, category, class)
	return redis.Strings(conn.Do("ZRANGE", k, 0, -1))
}

// GetAllClasses gets the names of all classes.
// Return:
// a slice contains the names of all classes.
// Class name format: $campus:$category:$class.
func (db *DB) GetAllClasses() ([]string, error) {
	var allClasses []string

	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	k := "ming:campuses"
	campuses, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return []string{}, err
	}

	for _, campus := range campuses {
		k = fmt.Sprintf("ming:%v:categories", campus)
		categories, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
		if err != nil {
			return []string{}, err
		}

		for _, category := range categories {
			k = fmt.Sprintf("ming:%v:%v:classes", campus, category)
			classes, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
			if err != nil {
				return []string{}, err
			}

			for _, class := range classes {
				allClasses = append(allClasses, fmt.Sprintf("%v:%v:%v", campus, category, class))
			}

		}
	}
	return allClasses, nil
}

// GetAllPeriods gets the names of all classes.
// Return:
// a slice contains the all periods.
// Period name format: $campus:$category:$period.
func (db *DB) GetAllPeriods() ([]string, error) {
	var allPeriods []string

	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	k := "ming:campuses"
	campuses, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return []string{}, err
	}

	for _, campus := range campuses {
		k = fmt.Sprintf("ming:%v:categories", campus)
		categories, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
		if err != nil {
			return []string{}, err
		}

		for _, category := range categories {
			k = fmt.Sprintf("ming:%v:%v:periods", campus, category)
			periods, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
			if err != nil {
				return []string{}, err
			}

			for _, period := range periods {
				allPeriods = append(allPeriods, fmt.Sprintf("%v:%v:%v", campus, category, period))
			}

		}
	}
	return allPeriods, nil
}

// GetClassesPeriods gets the period of each class.
// Return:
// a map. Key: class name($campus:$category:$class), value: period.
func (db *DB) GetClassesPeriods() (map[string]string, error) {
	m := map[string]string{}

	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return map[string]string{}, err
	}
	defer conn.Close()

	k := "ming:campuses"
	campuses, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return map[string]string{}, err
	}

	for _, campus := range campuses {
		k = fmt.Sprintf("ming:%v:categories", campus)
		categories, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
		if err != nil {
			return map[string]string{}, err
		}

		for _, category := range categories {
			k = fmt.Sprintf("ming:%v:%v:classes", campus, category)
			classes, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
			if err != nil {
				return map[string]string{}, err
			}

			for _, class := range classes {
				period, err := db.GetClassPeriod(campus, category, class)
				if err != nil {
					return map[string]string{}, err
				}

				k := fmt.Sprintf("%v:%v:%v", campus, category, class)
				m[k] = period
			}

		}
	}
	return m, nil
}

// GetAllPeriodsOfCategory gets all category's periods for all campuses.
//
// Params:
//
//	category: category which you want to get all periods.
//
// Returns:
//
//	a map contains all periods. key: campus, value: periods.
func (db *DB) GetAllPeriodsOfCategory(category string) (map[string][]string, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return map[string][]string{}, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:campuses", category)
	campuses, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return map[string][]string{}, err
	}

	periodsMap := map[string][]string{}
	for _, campus := range campuses {
		k = fmt.Sprintf("ming:%v:%v:periods", campus, category)
		periods, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
		if err != nil {
			return map[string][]string{}, err
		}

		if len(periods) > 0 {
			periodsMap[campus] = append(periodsMap[campus], periods...)
		}
	}

	return periodsMap, nil
}

// GetTeachers lists all teacher names.
func (db *DB) GetTeachers() ([]string, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	k := "ming:teachers"
	return redis.Strings(conn.Do("ZRANGE", k, 0, -1))
}

// GetStudentsOfTeacher gets all students of given teacher.
//
// Params:
//
//	teacher: name of teacher.
//
// Return:
//
//	students in the format: $STUDENT_NAME:$CONTACT_PHONE_NUM. e.g. 小明:13800138000
func (db *DB) GetStudentsOfTeacher(teacher string) ([]string, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	k := fmt.Sprintf("ming:%v:classes", teacher)
	classes, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
	if err != nil {
		return []string{}, err
	}

	allStudents := []string{}
	for _, class := range classes {
		k = fmt.Sprintf("ming:%v:students", class)
		students, err := redis.Strings(conn.Do("ZRANGE", k, 0, -1))
		if err != nil {
			return []string{}, err
		}

		if len(students) > 0 {
			allStudents = append(allStudents, students...)
		}
	}

	return allStudents, nil
}

// GetAllStudents gets all students.
// Return:
// Slice contains student in the format: $name:$phone_num.
func (db *DB) GetAllStudents() ([]string, error) {
	conn, err := redishelper.Dial(db.RedisServer, db.RedisPassword)
	if err != nil {
		return []string{}, err
	}
	defer conn.Close()

	allStudents := []string{}

	students, err := redis.Strings(conn.Do("ZRANGE", "ming:students", 0, -1))
	if err != nil {
		return []string{}, err
	}

	if len(students) > 0 {
		allStudents = append(allStudents, students...)
	}

	return allStudents, nil
}
