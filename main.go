package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Khai báo biến toàn cục cho kết nối DB
var DB *gorm.DB

// Định nghĩa model Student
type Student struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email" gorm:"unique"`
}

// Kết nối tới Database PostgreSQL thông qua GORM
func ConnectDatabase() {
	var err error
	// Chuỗi kết nối tới PostgreSQL
	dsn := "postgresql://student_lhl9_user:Ghm6I1HlTxIpgCSm6aqzWJ37kh6LZ7ng@dpg-cs9lukq3esus739hnch0-a.singapore-postgres.render.com/student_lhl9"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Tự động tạo bảng cho model Student
	DB.AutoMigrate(&Student{})
	fmt.Println("Database connected and student table migrated!")
}

// Lấy tất cả sinh viên
func GetStudents(c *gin.Context) {
	var students []Student
	DB.Find(&students)
	c.JSON(http.StatusOK, students)
}

// Lấy một sinh viên theo ID
func GetStudent(c *gin.Context) {
	var student Student
	id := c.Param("id")
	if err := DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found!"})
		return
	}
	c.JSON(http.StatusOK, student)
}

// Tạo sinh viên mới
func CreateStudent(c *gin.Context) {
	var student Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create student!"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// Cập nhật thông tin sinh viên theo ID
func UpdateStudent(c *gin.Context) {
	var student Student
	id := c.Param("id")

	if err := DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found!"})
		return
	}

	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	DB.Save(&student)
	c.JSON(http.StatusOK, student)
}

// Xóa sinh viên theo ID
func DeleteStudent(c *gin.Context) {
	var student Student
	id := c.Param("id")

	if err := DB.First(&student, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found!"})
		return
	}

	DB.Delete(&student)
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully!"})
}

func main() {
	// Kết nối tới cơ sở dữ liệu
	ConnectDatabase()

	// Khởi tạo router Gin
	r := gin.Default()
	// Cấu hình CORS cho phép React gửi yêu cầu
	r.Use(cors.Default())

	// Định nghĩa các route CRUD cho Student
	r.GET("/students", GetStudents)          // Lấy danh sách tất cả sinh viên
	r.GET("/students/:id", GetStudent)       // Lấy thông tin sinh viên theo ID
	r.POST("/students", CreateStudent)       // Tạo sinh viên mới
	r.PUT("/students/:id", UpdateStudent)    // Cập nhật thông tin sinh viên theo ID
	r.DELETE("/students/:id", DeleteStudent) // Xóa sinh viên theo ID

	// Chạy server trên cổng 8080
	r.Run(":8080")
}
