package main

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 不要忘了导入数据库驱动
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var db *gorm.DB

var sqlDB *sqlx.DB

// StudentsInit 连接学生数据库并进行初始化
func StudentsInit() {
	dsn := "root:redmi@qwe123@tcp(100.79.174.8:30123)/students?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("数据库连接失败：", err)
	} else {
		fmt.Println("数据库连接成功")
	}
	var tables []string
	if tables, err = db.Migrator().GetTables(); err != nil {
		fmt.Println("查询所有数据库表集合失败：", err)
	} else {
		fmt.Println("查询所有数据库表集合成功", tables)
	}
	if !Exists(tables, "students") {
		if err = db.Migrator().CreateTable(&Students{}); err != nil {
			fmt.Println("数据库初始化失败: ", err)
		} else {
			fmt.Println("students", "数据库初始化成功")
		}
	}
}

// TradeInit 连接交易数据库并进行初始化
func TradeInit() {
	dsn := "root:redmi@qwe123@tcp(100.79.174.8:30123)/trade?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("数据库连接失败：", err)
	} else {
		fmt.Println("数据库连接成功")
	}
	var tables []string
	if tables, err = db.Migrator().GetTables(); err != nil {
		fmt.Println("查询所有数据库表集合失败：", err)
	} else {
		fmt.Println("查询所有数据库表集合成功", tables)
	}
	if !Exists(tables, "accounts") {
		if err = db.Migrator().CreateTable(&Accounts{}); err != nil {
			fmt.Println("数据库初始化失败: ", err)
		} else {
			fmt.Println("accounts", "数据库初始化成功")
		}
	}
	if !Exists(tables, "transactions") {
		if err = db.Migrator().CreateTable(&Transactions{}); err != nil {
			fmt.Println("数据库初始化失败: ", err)
		} else {
			fmt.Println("transactions", "数据库初始化成功")
		}
	}
}

func Exists(str []string, target string) bool {
	for _, v := range str {
		if v == target {
			return true
		}
	}
	return false
}

// ## SQL语句练习
//### 题目1：基本CRUD操作
//- 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、
// age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。
//  - 要求 ：
//    - 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
//    - 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
//    - 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
//    - 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。

type Students struct {
	ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null;comment:学生姓名"`
	Age       uint8          `json:"age" gorm:"not null;default:0;comment:学生年龄"`
	Gender    string         `json:"gender" gorm:"type:varchar(20);not null;default:'未知';comment:学生性别"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:删除时间"`
}

//### 题目2：事务语句
//- 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和
// transactions 表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
//  - 要求 ：
//    - 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
//    如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

type Accounts struct {
	ID      uint64  `json:"id" gorm:"primaryKey;autoIncrement;comment:ID"`
	Balance float64 `json:"balance" gorm:"not null;default:0;comment:账户余额"`
}

type Transactions struct {
	ID            uint64  `json:"id" gorm:"primaryKey;autoIncrement;comment:ID"`
	FromAccountID uint64  `json:"from_account_id" gorm:"not null;default:0;comment:转出账户ID"`
	ToAccountID   uint64  `json:"to_account_id" gorm:"not null;default:0;comment:转入账户ID"`
	Amount        float64 `json:"amount" gorm:"not null;default:0;comment:转账金额"`
}

func AccountsTrade(tran Transactions) (string, error) {
	userA := Accounts{}
	// 检查账户 A 的余额是否足够
	db.Where("id = ?", tran.FromAccountID).First(&userA)
	if userA.Balance < tran.Amount {
		return "余额不足", nil
	}
	// 查询 B 账户余额
	userB := Accounts{}
	db.Where("id = ?", tran.ToAccountID).First(&userB)
	// 开启事物进行转账操作
	if err := db.Transaction(func(tx *gorm.DB) error {
		// 扣除 A 账户余额
		if err := tx.Model(&Accounts{}).Where("id = ?", tran.FromAccountID).Updates(map[string]interface{}{
			"ID":      userA.ID,
			"Balance": userA.Balance - tran.Amount,
		}).Error; err != nil {
			return errors.New("rollback Update userA: " + err.Error())
		}
		// 增加 B 账户余额
		if err := tx.Model(Accounts{}).Where("id = ?", tran.ToAccountID).Updates(map[string]interface{}{
			"ID":      userB.ID,
			"Balance": userB.Balance + tran.Amount,
		}).Error; err != nil {
			return errors.New("rollback Update userB: " + err.Error())
		}
		// 转账记录留存
		if err := tx.Model(Transactions{}).Create(&tran).Error; err != nil {
			return errors.New("rollback Create tran: " + err.Error())
		}
		return nil
	}); err != nil {
		return "", err
	}
	return "转账成功", nil
}

// ## Sqlx入门
//### 题目1：使用SQL扩展库进行查询
//- 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
//  - 要求 ：
//    - 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
//    - 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

// EmployeesInitDB 员工数据库连接及初始化
func EmployeesInitDB() (err error) {
	dsn := "root:redmi@qwe123@tcp(100.79.174.8:30123)/employees?charset=utf8mb4&parseTime=True&loc=Local"
	// 也可以使用MustConnect连接不成功就panic
	sqlDB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	return
}

type Employees struct {
	ID         uint64  `db:"id" gorm:"primaryKey;autoIncrement;comment:ID"`
	Name       string  `db:"name" gorm:"type:varchar(100);not null;comment:名字"`
	Department string  `db:"department" gorm:"type:varchar(100);not null;comment:部门"`
	Salary     float32 `db:"salary" gorm:"type:varchar(100);not null;comment:薪资"`
}

func CreateTable() {
	// 表初始化
	sqlStr := "CREATE TABLE IF NOT EXISTS employees (id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',name VARCHAR(100) NOT NULL COMMENT '员工姓名', department VARCHAR(100) NOT NULL COMMENT '所属部门', salary DECIMAL(10, 2) NOT NULL COMMENT '月薪(保留两位小数)', created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间', updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',PRIMARY KEY (id),INDEX idx_department (department),INDEX idx_salary (salary)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='员工信息表';"
	if _, err := sqlDB.Exec(sqlStr); err != nil {
		fmt.Printf("CREATE TABLE failed, err:%v\n", err)
		return
	}
}

func CreateEmployees(employees Employees) {
	sqlStr := "insert into employees(name, department, salary) values (?,?,?)"
	_, err := sqlDB.Exec(sqlStr, employees.Name, employees.Department, employees.Salary)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	return
}

// FiltrationEmployees 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
func FiltrationEmployees(department string) ([]Employees, error) {
	var employees []Employees
	sqlStr := "select id, name, department, salary from employees where department = ?"
	if err := sqlDB.Select(&employees, sqlStr, department); err != nil {
		return employees, err
	}
	return employees, nil
}

// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employees 结构体中。
func SalaryMax() ([]Employees, error) {
	var employees []Employees
	sqlStr := "select id, name, department, salary from employees ORDER BY salary DESC LIMIT 1"
	if err := sqlDB.Select(&employees, sqlStr); err != nil {
		return employees, err
	}
	return employees, nil
}

//### 题目2：实现类型安全映射
//- 假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
//  - 要求 ：
//    - 定义一个 Book 结构体，包含与 books 表对应的字段。
//    - 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

type Books struct {
	Books []Book
}

type Book struct {
	ID     uint64  `db:"id"`
	Title  string  `db:"title"`  // 书名
	Author string  `db:"author"` // 作者
	Price  float64 `db:"price"`  // 价格
}

// BookInitDB 书籍数据库连接及初始化
func BookInitDB() (err error) {
	dsn := "root:redmi@qwe123@tcp(100.79.174.8:30123)/books?charset=utf8mb4&parseTime=True&loc=Local"
	// 也可以使用MustConnect连接不成功就panic
	sqlDB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	return
}

// tables 创建
func CreateBooksTables() {
	sqlStr := "CREATE TABLE IF NOT EXISTS books ( id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '书籍ID', title VARCHAR(100) NOT NULL COMMENT '书名', author VARCHAR(50) NOT NULL COMMENT '作者', price DECIMAL(6,2) NOT NULL COMMENT '价格', PRIMARY KEY (id), INDEX idx_price (price)\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='书籍信息表';"
	if _, err := sqlDB.Exec(sqlStr); err != nil {
		fmt.Printf("CREATE TABLE failed, err:%v\n", err)
		return
	}
}

// 增加书籍数据
func CreateBooks(book Book) {
	sqlStr := "insert into books(title, author, price) values (?,?,?)"
	_, err := sqlDB.Exec(sqlStr, book.Title, book.Author, book.Price)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	return
}

// 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
func QueryBooksByPrice(minPrice float64) (Books, error) {
	strStr := `
		SELECT id, title, author, price 
		FROM books 
		WHERE price > ?
		ORDER BY price DESC
	`
	// 创建Book结构体切片用于存储结果
	var books Books

	// 执行查询并将结果映射到结构体切片
	if err := sqlDB.Select(&books.Books, strStr, minPrice); err != nil {
		return books, fmt.Errorf("查询书籍失败: %w", err)
	}
	return books, nil
}

// ## 进阶gorm
//### 题目1：模型定义
//- 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
//  - 要求 ：
//    - 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
//    - 编写Go代码，使用Gorm创建这些模型对应的数据库表。

func BlogsInitDB() {
	dsn := "root:redmi@qwe123@tcp(100.79.174.8:30123)/blogs?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		fmt.Println("数据库连接失败：", err)
	} else {
		fmt.Println("数据库连接成功")
	}
	var tables []string
	if tables, err = db.Migrator().GetTables(); err != nil {
		fmt.Println("查询所有数据库表集合失败：", err)
	} else {
		fmt.Println("查询所有数据库表集合成功", tables)
	}
	if !Exists(tables, "users") {
		if err = db.Migrator().CreateTable(&User{}); err != nil {
			fmt.Println("数据库初始化失败: ", err)
		} else {
			fmt.Println("users", "数据库初始化成功")
		}
	}
	if !Exists(tables, "posts") {
		if err = db.Migrator().CreateTable(&Post{}); err != nil {
			fmt.Println("数据库初始化失败: ", err)
		} else {
			fmt.Println("posts", "数据库初始化成功")
		}
	}
	if !Exists(tables, "comments") {
		if err = db.Migrator().CreateTable(&Comment{}); err != nil {
			fmt.Println("数据库初始化失败: ", err)
		} else {
			fmt.Println("comments", "数据库初始化成功")
		}
	}
}

type Blogs struct{}

type Operate struct {
	ID        uint64         `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:删除时间"`
}

type User struct {
	Operate
	Name    string `json:"name" gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;charset=utf8mb4;comment:姓名"`
	Age     uint8  `json:"age" gorm:"not null;default:0;comment:年龄"`
	Gender  string `json:"gender" gorm:"type:varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;default:未知;comment:性别"`
	PostNum uint64 `json:"post_num" gorm:"not null;default:0;comment:文章数量"`
	Posts   []Post `json:"posts" gorm:"foreignKey:UserID"`
}

type Post struct {
	Operate
	UserID        uint64    `json:"id" gorm:"comment:主键ID"`
	Name          string    `json:"name" gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;comment:文章名称"`
	CommentNum    uint64    `json:"comment_num" gorm:"not null;default:0;comment:评论数量"`
	Content       string    `json:"content" gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;comment:内容"`
	Comments      []Comment `json:"comments" gorm:"foreignKey:PostID"`
	CommentStatus string    `json:"comment_status" gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;comment:评论状态"`
}

type Comment struct {
	Operate
	PostID  uint64 `json:"id" gorm:"comment:主键ID"`
	Content string `json:"content" gorm:"type:varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;not null;comment:内容"`
}

func (b Blogs) CreateUser(user User) {
	db.Create(&user)
}

func (b Blogs) CreatePost(post Post) {
	db.Create(&post)
}

func (b Blogs) CreateComment(Comment Comment) {
	db.Create(&Comment)
}

func (b Blogs) DeleteComment(Comment Comment) {
	db.Model(&Comment).Where("id = ?", Comment.ID).Delete(&Comment)
}

// ### 题目2：关联查询
// - 基于上述博客系统的模型定义。
//   - 要求 ：
//   - 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
//   - 编写Go代码，使用Gorm查询评论数量最多的文章信息。

// GetAllPostComment 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
func (b Blogs) GetAllPostComment(userID uint64) (*User, error) {
	var user User
	err := db.Preload("Posts.Comments"). // 嵌套预加载：用户->文章->评论
		Where("id = ?", userID).
		First(&user).Error

	if err != nil {
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	return &user, nil
}

// 查询结果结构体
type PostCommentCount struct {
	PostID uint64 `gorm:"column:post_id"`
	Count  int    `gorm:"column:count"`
}

// GetPostNumMax 编写Go代码，使用Gorm查询评论数量最多的文章信息。
func (b Blogs) GetPostNumMax() (*Post, error) {
	var result Post
	var postComCount PostCommentCount
	if err := db.Model(&Comment{}).
		Select("post_id, COUNT(*) as count").
		Group("post_id").
		Order("count DESC").
		Limit(1).
		Scan(&postComCount).Error; err != nil {
		return nil, fmt.Errorf("查询 Comment 错误: %w", err)
	}
	if err := db.Model(Post{}).Where("id = ?", postComCount.PostID).Find(&result).Error; err != nil {
		return nil, fmt.Errorf("查询 Post 错误: %w", err)
	}
	return &result, nil
}

//### 题目3：钩子函数
//- 继续使用博客系统的模型。
//  - 要求 ：
//    - 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
//    - 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

// AfterCreate 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	log.Printf("检测到新文章创建 (ID: %d, 用户ID: %d), 更新用户文章数量...", p.ID, p.UserID)

	// 更新用户的文章数量字段
	result := tx.Model(&User{}).
		Where("id = ?", p.UserID).
		Update("post_num", gorm.Expr("post_num + 1"))

	if result.Error != nil {
		log.Printf("更新用户文章数量失败: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("警告: 未找到用户ID: %d", p.UserID)
	}

	log.Printf("用户ID: %d 的文章数量已更新", p.UserID)
	return nil
}

// AfterCreate 为 Post 模型添加一个钩子函数，在评论创建时自动更新文章的评论数量统计字段。
func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	log.Printf("检测到新评论创建 (ID: %d, 文章ID: %d), 更新用户文章数量...", c.ID, c.PostID)

	// 更新文章的评论数量字段
	result := tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_num", gorm.Expr("comment_num + 1"))

	if result.Error != nil {
		log.Printf("更新用户文章数量失败: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("警告: 未找到文章ID: %d", c.PostID)
	}

	log.Printf("文章ID: %d 的文章数量已更新", c.PostID)
	return nil
}

// AfterDelete 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
// 在评论删除后触发的钩子函数
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	log.Printf("检测到评论删除 (ID: %d, 文章ID: %d), 检查文章评论状态...", c.ID, c.PostID)

	// 统计文章当前的评论数量
	var count int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		log.Printf("统计评论数量失败: %v", err)
		return err
	}

	log.Printf("文章ID: %d 剩余评论数: %d", c.PostID, count)

	// 更新文章信息
	if count == 0 {
		if err := db.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_num", count).Update("comment_status", "无评论").Error; err != nil {
			return err
		}
	} else {
		if err := db.Model(&Post{}).Where("id = ?", c.PostID).Update("comment_num", count).Update("comment_status", "有评论").Error; err != nil {
			return err
		}
	}

	log.Printf("文章ID: %d 状态更新完成", c.PostID)
	return nil
}

func main() {
	// 1.1 学生数据库初始化
	StudentsInit()
	// 1.1.1 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"
	user := Students{Name: "张三", Age: 20, Gender: "三年级", CreatedAt: time.Now()}
	result := db.Create(&user)
	fmt.Println(result.RowsAffected)

	// 1.1.2 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	var findStudent []Students
	db.Where("age > ?", 18).Find(&findStudent)
	fmt.Println(findStudent)

	// 1.1.3 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	db.Model(Students{}).Where("name = ?", "张三").Updates(&Students{Gender: "四年级"})

	// 1.1.4 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	user = Students{Name: "李四", Age: 13, Gender: "一年级", CreatedAt: time.Now()}
	result = db.Create(&user)
	fmt.Println(result.RowsAffected)

	savedStudent := Students{}
	db.Debug().Find(&savedStudent)
	fmt.Println(savedStudent)

	db.Where("age < ?", 18).Delete(&Students{})

	// 2.1
	// 交易数据库初始化
	TradeInit()
	// 前置条件 创建A B 两个账户存 100元
	userA := Accounts{Balance: 100}
	db.Create(&userA)
	userB := Accounts{Balance: 100}
	db.Create(&userB)
	// 2.1.1 从账户 A 向账户 B 转账 100 元的操作
	trade, err := AccountsTrade(Transactions{FromAccountID: 1, ToAccountID: 2, Amount: 100})
	fmt.Println(trade, err)
	// 3.1
	// 员工数据库初始化
	if err := EmployeesInitDB(); err != nil {
		fmt.Println(err)
	}
	// 初始化表
	CreateTable()
	// 增加员工信息
	CreateEmployees(Employees{Name: "张三", Department: "技术部", Salary: 35000})
	CreateEmployees(Employees{Name: "张二", Department: "技术部", Salary: 50000})
	CreateEmployees(Employees{Name: "李四", Department: "财务部", Salary: 10000})
	// 3.1.1 只要技术部的信息
	fmt.Println(FiltrationEmployees("技术部"))
	// 3.1.2 查看薪资最高的同事
	fmt.Println(SalaryMax())
	// 3.2 书籍创建
	// 连接数据库
	if err := BookInitDB(); err != nil {
		fmt.Println(err)
	}
	// 初始化书籍表
	CreateBooksTables()
	// 创建数据
	CreateBooks(Book{Title: "七个习惯", Price: 80})
	CreateBooks(Book{Title: "Go程序设计语言", Author: "Alan A. A. Donovan", Price: 70})
	// 3.2.1 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。
	fmt.Println(QueryBooksByPrice(50))
	// 博客系统数据库及数据表初始化
	BlogsInitDB()
	blogs := Blogs{}
	blogs.CreateUser(User{Name: "张三", Age: 30, Gender: "男"})
	blogs.CreateUser(User{Name: "李四", Age: 35, Gender: "男"})
	blogs.CreatePost(Post{Name: "Golang 修养之路", UserID: 1, Content: "Hello World!"})
	blogs.CreatePost(Post{Name: "Golang 程序设计", UserID: 2, Content: "Hello Golang!"})
	blogs.CreateComment(Comment{PostID: 1, Content: "good"})
	blogs.CreateComment(Comment{PostID: 2, Content: "good"})
	blogs.CreateComment(Comment{PostID: 2, Content: "good"})
	// 3.3.1 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
	fmt.Println(blogs.GetAllPostComment(1))
	// 3.3.2 编写Go代码，使用Gorm查询评论数量最多的文章信息。
	fmt.Println(blogs.GetPostNumMax())
	// 3.4 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。
	blogs.DeleteComment(Comment{Operate: Operate{ID: 2}, PostID: 2, Content: "good"})
}
