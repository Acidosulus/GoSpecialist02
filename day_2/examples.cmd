rem restish "http://127.0.0.1:8000/v1/employee/1"


rem restish POST "http://127.0.0.1:8000/v1/employee" DEPARTMENT:"����ࠫ�� ������",FIO:"����஢ ����� ����஢��",POSITION:"��������� �᭮��묨 ���ਠ����"
rem restish POST "http://127.0.0.1:8000/v1/employee" DEPARTMENT:"����ࠫ�� ������",FIO:"��ᨬ��� �ઠ��� ��쨭�譠",POSITION:"��壠���"

rem restish http://127.0.0.1:8000/v1/employee/1
rem restish http://127.0.0.1:8000/v1/employee/2
rem restish http://127.0.0.1:8000/v1/employee/3

rem restish POST http://127.0.0.1:8000/v1/in/"������ ���� ��������"
rem restish POST http://127.0.0.1:8000/v1/in/"��ᨬ��� �ઠ��� ��쨭�譠"

rem restish POST http://127.0.0.1:8000/v1/out/"������ ���� ��������"
rem restish POST http://127.0.0.1:8000/v1/out/"��ᨬ��� �ઠ��� ��쨭�譠"

restish http://127.0.0.1:8000/v1/worktime/"��ᨬ��� �ઠ��� ��쨭�譠"/"2023-05-11 00:00:00"/"2023-05-19 23:59:59"

restish http://127.0.0.1:8000/v1/worktime/"������ ���� ��������"/"2023-05-18 02:30:00"/"2023-05-19 23:59:59"
