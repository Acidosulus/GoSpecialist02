package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// DB Driver visible to whole program
var DB *sql.DB

// StationResource holds information about locations
type StationResource struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	OpeningTime string `json:"opening_time"`
	ClosingTime string `json:"closing_time"`
}

// GetStation returns the station detail
func GetStation(c *gin.Context) {
	var station StationResource
	id := c.Param("station_id")
	err := DB.QueryRow("select ID, NAME, CAST(OPENING_TIME as CHAR), CAST(CLOSING_TIME as CHAR) from station where id=?", id).Scan(&station.ID, &station.Name, &station.OpeningTime, &station.ClosingTime)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.JSON(200, gin.H{
			"result": station,
		})
	}
}

// CreateStation handles the POST
func CreateStation(c *gin.Context) {
	var station StationResource
	// Parse the body into our resrource
	if err := c.BindJSON(&station); err == nil {
		// Format Time to Go time format
		statement, _ := DB.Prepare("insert into station (NAME, OPENING_TIME, CLOSING_TIME) values (?, ?, ?)")
		result, _ := statement.Exec(station.Name, station.OpeningTime, station.ClosingTime)
		if err == nil {
			newID, _ := result.LastInsertId()
			station.ID = int(newID)
			c.JSON(http.StatusOK, gin.H{
				"result": station,
			})
		} else {
			c.String(http.StatusInternalServerError, err.Error())
		}
	} else {
		c.String(http.StatusInternalServerError, err.Error())
	}
}

// RemoveStation handles the removing of resource
func RemoveStation(c *gin.Context) {
	id := c.Param("station-id")
	statement, _ := DB.Prepare("delete from station where id=?")
	_, err := statement.Exec(id)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		c.String(http.StatusOK, "")
	}
}

// -------------------------------------------------------------------------------------------------------------------------------------------

// Работник
type TEMPLOYEES struct {
	ID         int    `json:"ID"`
	FIO        string `json:"FIO"`
	DEPARTMENT string `json:"DEPARTMENT"`
	POSITION   string `json:"POSITION"`
}

// время входа-выхода
type TWORKINGTIME struct {
	ID       int    `json:"ID"`
	EMPLOYEE int    `json:"EMPLOYEE"`
	IN       string `json:"IN"`
	OUT      string `json:"OUT"`
}

// Добавление нового работника
func PostEmployee(context *gin.Context) {
	var employee, empl TEMPLOYEES

	bts, _ := ioutil.ReadAll(context.Request.Body) //тело запроса
	fmt.Println(string(bts))

	if err := context.BindJSON(&employee); err == nil {
		Database.QueryRow("select ID, FIO, DEPARTMENT, POSITION from EMPLOYEES where FIO=?", employee.FIO).Scan(&empl.ID, &empl.FIO, &empl.DEPARTMENT, &empl.POSITION)
		if empl.ID != 0 {
			context.String(http.StatusInternalServerError, "Запись уже есть в базе данных, даннные не добавлены")
		} else {
			statement, _ := Database.Prepare("insert into EMPLOYEES (FIO, DEPARTMENT, POSITION) values (?, ?, ?)")
			result, _ := statement.Exec(employee.FIO, employee.DEPARTMENT, employee.POSITION)
			if err == nil {
				newID, _ := result.LastInsertId()
				employee.ID = int(newID)
				context.JSON(http.StatusOK, gin.H{"result": employee})
			} else {
				context.String(http.StatusInternalServerError, err.Error())
			}
		}
	} else {
		context.String(http.StatusInternalServerError, err.Error())
	}
}

// Получение данных работника по EMPLOYEES.ID
func GetEmployeeByID(context *gin.Context) {
	var employee TEMPLOYEES
	id := context.Param("employee_id")
	err := Database.QueryRow("select ID, FIO, DEPARTMENT, POSITION from EMPLOYEES where id=?", id).Scan(&employee.ID, &employee.FIO, &employee.DEPARTMENT, &employee.POSITION)
	if err != nil {
		log.Println(err)
		context.JSON(500, gin.H{
			"error": err.Error(),
		})
	} else {
		context.JSON(200, gin.H{
			"result": employee,
		})
	}

}

func Init_data_Base() *sql.DB {
	db, err := sql.Open("sqlite3", "./databases/accounting_of_working_time.db")
	if err != nil {
		log.Println("DB driver failed!")
	}
	return db
}

// возвращает идентификато работника по его ФИО
func GetEmployeeIDbyFIO(fio string) (int, error) {
	var empl TEMPLOYEES
	Database.QueryRow("select ID, FIO, DEPARTMENT, POSITION from EMPLOYEES where FIO=?", fio).Scan(&empl.ID, &empl.FIO, &empl.DEPARTMENT, &empl.POSITION)
	if empl.ID != 0 {
		return empl.ID, nil
	} else {
		return 0, fmt.Errorf("Не найдено ФИО работника")

	}
}

// если это возможно, фиксируем работника приступившим к работе, возвращает row_id записи в таблице workingtime, в случае удачной фиксации
func EmployeeIn(fio string) (int, error) {
	var workingtime TWORKINGTIME
	employee_id, err := GetEmployeeIDbyFIO(fio)
	if err == nil { // проверяем, что такой работник существует
		Database.QueryRow("SELECT ID, EMPLOYEE, [IN], OUT FROM WORKINGTIME where EMPLOYEE=? and (out is null)", employee_id).Scan(&workingtime.ID, &workingtime.EMPLOYEE, &workingtime.IN, &workingtime.OUT)
		if workingtime.ID > 0 { // если запрос вернул что-то, то работник сейчас на работе и его нельзя учёсть пришедним на работу ещё раз
			return 0, fmt.Errorf("Работник уже находится на работе с " + workingtime.IN)
		} else {
			statement, _ := Database.Prepare("insert into WORKINGTIME (EMPLOYEE, [IN]) values (?, datetime())")
			result, _ := statement.Exec(employee_id)
			newID, _ := result.LastInsertId()
			return int(newID), nil
		}
	} else {
		return 0, err
	}
}

// если это возможно, фиксируем работника прекратившим работу, возвращает row_id записи в таблице workingtime, в случае удачной фиксации
func EmployeeOut(fio string) (int, error) {
	var workingtime TWORKINGTIME
	employee_id, err := GetEmployeeIDbyFIO(fio)
	if err == nil { // проверяем, что такой работник существует
		Database.QueryRow("SELECT ID, EMPLOYEE, [IN], OUT FROM WORKINGTIME where EMPLOYEE=? and out is null", employee_id).Scan(&workingtime.ID, &workingtime.EMPLOYEE, &workingtime.IN, &workingtime.OUT)
		if workingtime.ID > 0 { // если запрос вернул что-то, то работник сейчас на работе и его можно зафиксировать прекратившим работу
			statement, _ := Database.Prepare("update WORKINGTIME set out = datetime() where EMPLOYEE=? and (out is null)")
			statement.Exec(employee_id)
			return workingtime.ID, nil
		} else {
			return 0, fmt.Errorf("Работник сейча не находится на работе, его невозможно зафиксировать прекратившим работу")
		}
	} else {
		return 0, err
	}
}

var Database *sql.DB

func main() {

	Database = Init_data_Base()
	fmt.Println(EmployeeIn("Зосимова Аркадия Ильинишна"))
	fmt.Println(EmployeeOut("Зосимова Аркадия Ильинишна"))
	fmt.Println("================================================================================================================================")

	router := gin.Default()
	// Add routes to REST verbs
	router.GET("/v1/employee/:employee_id", GetEmployeeByID)
	router.POST("/v1/employee", PostEmployee) //принимает добавляемы в базу данные в виден json
	router.POST("/v1/in/:employee_fio", PostEmployee)
	//r.DELETE("/v1/stations/:station_id", RemoveStation)

	router.Run("127.0.0.1:8000")
}
