package usr

import (
        "encoding/json"
        "errors"
        "fmt"
        "log"

        "comprarmas.com.mx/internal/model"

        "comprarmas.com.mx/internal/database"
)

func LoginUsr(sMail string, sPwd string) (string, string, error) {
        blnOK := false
        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("LoginUsr_Error cnn DB " + err.Error())
                return "", "", err
        }
        defer db.Close()

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("LoginUsr_Error cnn DB " + err.Error())
                return "", "", err
        }

        defer db.Close()
        sSQL := "SELECT DISTINCT a.id_Usr, a.Mail, a.Pwd, a.Estatus, a.Opt "
        sSQL += "FROM Usr a "
        sSQL += "WHERE (a.Mail = '" + sMail + "') "
        sSQL += "AND (a.Pwd = '" + sPwd + "') "
        sSQL += "AND (a.Estatus = 'ACTIVO') "

        resDat, err := db.Query(sSQL)
        if err != nil {
                return "", "", err
        } else {
                for resDat.Next() {
                        var datUsr = model.DataUsrGet{}
                        resDat.Scan(
                                &datUsr.IdUsr, &datUsr.Mail, &datUsr.Pwd, &datUsr.St, &datUsr.Opt,
                        )
                        blnOK = true
                        return fmt.Sprint(datUsr.IdUsr), fmt.Sprint(datUsr.Opt), nil
                }
        }
        if blnOK == false {
                err = errors.New("Usuario no encontrado.")
                return "", "", err
        }
        return "", "", nil
} //END LoginUsr

func GetUsr(sMail string, sSeller string) error {
        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("GetUsr_Error cnn DB " + err.Error())
                return err
        }
        defer db.Close()
        var errOpt error

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("GetUsr_Error cnn DB " + err.Error())
                return err
        }

        defer db.Close()
        sSQL := "SELECT DISTINCT a.Mail "
        sSQL += "FROM Usr a "
        sSQL += "WHERE (a.Mail = '" + sMail + "') "
        if sSeller != "" {
                sSQL += "AND (a.Opt='" + sSeller + "')"
        }

        resDat, err := db.Query(sSQL)
        if err != nil {
                return err
        } else {
                for resDat.Next() {
                        var datUsr = model.DataUsrGet{}
                        resDat.Scan(
                                &datUsr.Mail, &datUsr.St,
                        )
                        if datUsr.Mail == "" {
                                errOpt = errors.New("Ya existe el usuario.")
                                return errOpt
                        }
                        err = error(nil)
                }
        }
        return nil
} //END GetUsr

func RegisterUsr(sMail string, sPwd string) (string, error) {
        var lastId int64

        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("RegisterUsr_Error " + err.Error())
                return "", err
        }
        defer db.Close()

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("RegisterUsr_Error_Ping " + err.Error())
                return "", err
        }

        defer db.Close()
        sSQL := "INSERT INTO Usr(Mail, Pwd) values('" + sMail + "', '" + sPwd + "') "

        resDB, err := db.Exec(sSQL)
        if err != nil {
                log.Println("No fue posible registrar el usuario: " + err.Error())

                return "", err
        } else {
                lastId, err = resDB.LastInsertId()
                log.Println(resDB)
        }
        return fmt.Sprint(lastId), nil
} //END RegisterUsr

func GetUsrsSeller() ([]map[string]interface{}, error) {
        var allData []map[string]interface{}

        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("GetUsrSeller_Error cnn DB " + err.Error())
                return nil, err
        }
        defer db.Close()

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("GetUsrSeller_Error cnn DB " + err.Error())
                return nil, err
        }

        defer db.Close()
        sSQL := "SELECT DISTINCT a.id_Usr, a.Mail "
        sSQL += "FROM Usr a "
        sSQL += "WHERE (a.Opt='seller')"

        resDat, err := db.Query(sSQL)
        if err != nil {
                return nil, err
        } else {
                for resDat.Next() {
                        var datUsr = model.DataUsrGet{}
                        resDat.Scan(
                                &datUsr.IdUsr, &datUsr.Mail,
                        )

                        mapJson := make(map[string]interface{})
                        mapJson["id_usr"] = fmt.Sprint(datUsr.IdUsr)
                        mapJson["mail"] = fmt.Sprint(datUsr.Mail)
                        allData = append(allData, mapJson)
                }
        }

        sJson, err := json.Marshal(allData)
        fmt.Println(string(sJson))
        return allData, nil
} //END GetUsrSeller
