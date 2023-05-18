rem restish "http://127.0.0.1:8000/v1/employee/1"


rem restish POST "http://127.0.0.1:8000/v1/employee" DEPARTMENT:"Центральный Аппарат",FIO:"Сидоров Сидор Сидорович",POSITION:"Снабжение основными материалами"
rem restish POST "http://127.0.0.1:8000/v1/employee" DEPARTMENT:"Центральный Аппарат",FIO:"Зосимова Аркадия Ильинишна",POSITION:"Бухгалтер"

rem restish http://127.0.0.1:8000/v1/employee/1
rem restish http://127.0.0.1:8000/v1/employee/2
rem restish http://127.0.0.1:8000/v1/employee/3

rem restish POST http://127.0.0.1:8000/v1/in/"Иванов Иван Иванович"
rem restish POST http://127.0.0.1:8000/v1/in/"Зосимова Аркадия Ильинишна"

rem restish POST http://127.0.0.1:8000/v1/out/"Иванов Иван Иванович"
rem restish POST http://127.0.0.1:8000/v1/out/"Зосимова Аркадия Ильинишна"

restish http://127.0.0.1:8000/v1/worktime/"Зосимова Аркадия Ильинишна"/"2023-05-11 00:00:00"/"2023-05-19 23:59:59"

restish http://127.0.0.1:8000/v1/worktime/"Иванов Иван Иванович"/"2023-05-18 02:30:00"/"2023-05-19 23:59:59"
