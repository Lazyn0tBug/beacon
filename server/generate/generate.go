package main

import (
	"context"

	"github.com/Lazyn0tBug/beacon/server/generate/method"
	"github.com/Lazyn0tBug/beacon/server/initialize"
	"github.com/Lazyn0tBug/beacon/server/model"
	"github.com/Lazyn0tBug/beacon/server/model/system"
	"github.com/Lazyn0tBug/beacon/server/utils"

	"gorm.io/gen"
)

func main() {
	Logger := utils.GetLogger()
	db := initialize.DB(context.Background())
	if db == nil {
		Logger.Error("failed to connect database")
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	// gormdb, _ := gorm.Open(mysql.Open("root:@(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"))
	g.UseDB(db) // reuse your gorm db

	g.ApplyBasic(model.User{}, model.Case{}, model.Customer{}, model.Role{}, model.Permission{}, model.Hospital{}, model.Doctor{}, model.MedicalRecord{}, model.Appointment{}, model.MemberLevel{}, model.ServiceItem{}) // g.ApplyBasic(g.GenerateModel("User"),

	g.ApplyBasic(system.JwtBlacklist{}, system.JwtInActive{})

	// apply diy interfaces on structs or table models
	g.ApplyInterface(func(method.UserMethod) {}, model.User{}) // struct test will be ignored

	// g.ApplyBasic(
	// 	// Generate struct `User` based on table `users`
	// 	g.GenerateModel("users"),

	// 	// Generate struct `Employee` based on table `users`
	// g.GenerateModelAs("User", "Employee"),

	// 	// Generate struct `User` based on table `users` and generating options
	// 	g.GenerateModel("users", gen.FieldIgnore("address"), gen.FieldType("id", "int64")),

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(model.Method) {}, &model.User{})
	// g.ApplyInterface(func(model.CompMethod) {}, &model.Company{}) // struct test will be ignored
	// g.ApplyInterface(func(Querier) {}, g.GenerateModel("User"))

	// Generate the code
	g.Execute()
	// 迁移 schema
	// err = db.AutoMigrate(&Product{}, &model.User{})
	// if err != nil {
	// 	fmt.Println("failed to migrate database")
	// }
	// // Create
	// db.Create(&Product{Code: "D42", Price: 100})
}
