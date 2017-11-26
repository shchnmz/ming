package ming

import (
	"fmt"
)

var (
	// RedisKeyPrefix is for all keys in redis.
	RedisKeyPrefix = "ming"
)

// GetKeyOfCampuses returns the key to store all campuses.
func GetKeyOfCampuses() string {
	return fmt.Sprintf("%v:campuses", RedisKeyPrefix)
}

// GetKeyOfCampusCategories returns the key to store categories of the campus.
func GetKeyOfCampusCategories(campus string) string {
	return fmt.Sprintf("%v:%v:categories", RedisKeyPrefix, campus)
}

// GetKeyOfCampusCategoryClasses returns the key to store classes of the combination:
// campus and category.
func GetKeyOfCampusCategoryClasses(campus, category string) string {
	return fmt.Sprintf("%v:%v:%v:classes", RedisKeyPrefix, campus, category)
}

// GetKeyOfCategoryCampuses returns the key to store campuses of the category.
func GetKeyOfCategoryCampuses(category string) string {
	return fmt.Sprintf("%v:%v:campuses", RedisKeyPrefix, category)
}

// GetKeyOfTeachers returns the key to store all teachers.
func GetKeyOfTeachers() string {
	return fmt.Sprintf("%v:teachers", RedisKeyPrefix)
}

// GetKeyOfCampusCategoryClassTeachers returns the key to store teachers of the combination:
// campus, category and class.
func GetKeyOfCampusCategoryClassTeachers(campus, category, class string) string {
	return fmt.Sprintf("%v:%v:%v:%v:teachers", RedisKeyPrefix, campus, category, class)
}

// GetKeyOfTeacherClasses returns the key to store classes of the teacher.
func GetKeyOfTeacherClasses(teacher string) string {
	return fmt.Sprintf("%v:%v:classes", RedisKeyPrefix, teacher)
}

// GetKeyOfCampusCategoryClassPeriod returns the key to store period of the combination:
// campus, category and class.
func GetKeyOfCampusCategoryClassPeriod(campus, category, class string) string {
	return fmt.Sprintf("%v:%v:%v:%v:period", RedisKeyPrefix, campus, category, class)
}

// GetKeyOfCampusCategoryPeriods returns the key to store periods of the combination:
// campus and cateogry.
func GetKeyOfCampusCategoryPeriods(campus, category string) string {
	return fmt.Sprintf("%v:%v:%v:periods", RedisKeyPrefix, campus, category)
}

// GetKeyOfStudents returns the key to store all students.
func GetKeyOfStudents() string {
	return fmt.Sprintf("%v:students", RedisKeyPrefix)
}

// GetKeyOfStudentNamePhoneNumClasses returns the key to store classes of the combination:
// student name and student phone num.
func GetKeyOfStudentNamePhoneNumClasses(name, phoneNum string) string {
	return fmt.Sprintf("%v:%v:%v:classes", RedisKeyPrefix, name, phoneNum)
}

// GetKeyOfPhoneNums returns the key to store all phone numbers.
func GetKeyOfPhoneNums() string {
	return fmt.Sprintf("%v:phone_nums", RedisKeyPrefix)
}

// GetKeyOfPhoneNumStudents returns the key to store students of the phone number.
func GetKeyOfPhoneNumStudents(phoneNum string) string {
	return fmt.Sprintf("%v:%v:students", RedisKeyPrefix, phoneNum)
}

// GetKeyOfCampusCategoryClassStudents returns the key to store students of the combination:
// campus, category and class.
func GetKeyOfCampusCategoryClassStudents(campus, category, class string) string {
	return fmt.Sprintf("%v:%v:%v:%v:students", RedisKeyPrefix, campus, category, class)
}
