package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

// DB Driver visible to whole program
var DB *sql.DB

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
		context.JSON(500, gin.H{"error": err.Error()})
	} else {
		context.JSON(200, gin.H{"result": employee})
	}

}

func Init_data_Base() *sql.DB {
	db, err := sql.Open("sqlite3", "./databases/accounting_of_working_time.db")
	if err != nil {
		log.Println("DB driver failed!")
	}
	return db
}

// возвращает идентификатор работника по его ФИО из БД
func GetEmployeeIDbyFIO(fio string) (int, error) {
	var empl TEMPLOYEES
	Database.QueryRow("select ID, FIO, DEPARTMENT, POSITION from EMPLOYEES where FIO=?", fio).Scan(&empl.ID, &empl.FIO, &empl.DEPARTMENT, &empl.POSITION)
	if empl.ID != 0 {
		return empl.ID, nil
	} else {
		return 0, fmt.Errorf("Не найдено ФИО работника")
	}
}

// фиксирует начало работы работника с переданным ФИО
func PostEmployeeIn(context *gin.Context) {
	fio := context.Param("employee_fio")
	workingtime_id, err := EmployeeIn(fio)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
	} else {
		context.JSON(200, gin.H{"workingtime_id": workingtime_id})
	}
}

// фиксирует конец работы работника с переданным ФИО
func PostEmployeeOut(context *gin.Context) {
	fio := context.Param("employee_fio")
	workingtime_id, err := EmployeeOut(fio)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
	} else {
		context.JSON(200, gin.H{"workingtime_id": workingtime_id})
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

// возвращает количество дней, часов, минут по переданному количеству секунд, принимается, что в рабочем дне 8 часов
func TimeIntervalsFromSeconds(seconds int) (days, hours, minutes int) {
	days = seconds / (60 * 60 * 8)
	hours = (seconds - (days * (60 * 60 * 8))) / (60 * 60)
	minutes = (seconds - (days * (60 * 60 * 8)) - hours*(60*60)) / (60)
	return days, hours, minutes
}

// Возвращает количество отработанного времени по ФИО работника и дата-времени начала и конца анализируемого периода
func GetWorkingTime(context *gin.Context) {
	employee_fio := context.Param("employee_fio")
	date_from := context.Param("datetime_from")
	date_to := context.Param("datetime_to")
	var employee_id, seconds int
	err := Database.QueryRow(`SELECT 
								employee,
								sum(Round((julianday(ifnull(OUT,DATETIME())) - julianday([IN]))*86400,0)) as SECONDS_DIFF -- сумма в секундах всех промежутков, на которых работник работал
								from WORKINGTIME
								where 	[IN]>=? and 
										ifnull(OUT,DATETIME())<=? AND 
										EMPLOYEE = (select id from EMPLOYEES where fio =? limit 1)
								group by employee;`, date_from, date_to, employee_fio).Scan(&employee_id, &seconds)
	days, hours, minutes := TimeIntervalsFromSeconds(seconds)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
	} else {
		context.JSON(200, gin.H{"employee_fio": employee_fio,
			"employee_id":          employee_id,
			"total_seconds_worked": seconds,
			"days_worked":          days,
			"hours_worked":         hours,
			"minutes_worked":       minutes})
	}
}

var Database *sql.DB

func main() {
	Database = Init_data_Base()

	router := gin.Default()
	router.GET("/v1/employee/:employee_id", GetEmployeeByID)                             //возвращает ФИО работника оп его ID
	router.POST("/v1/employee", PostEmployee)                                            //принимает добавляемые в базу данные в виде json
	router.POST("/v1/in/:employee_fio", PostEmployeeIn)                                  //фиксирует начало работы работника с переданным ФИО
	router.POST("/v1/out/:employee_fio", PostEmployeeOut)                                //фиксирует конец работы работника с переданным ФИО
	router.GET("/v1/worktime/:employee_fio/:datetime_from/:datetime_to", GetWorkingTime) //получает фио работника и интервал на котором будет посчитано количество отработанного времени, возвращает json с общим количеством отработанных секунд и количеством отработанных 8часовых рабочих дней, и часов, минут

	router.Run("127.0.0.1:8000")
}
