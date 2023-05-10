package product

import (
        "encoding/json"
        "errors"
        "fmt"
        "log"
        "time"

        "comprarmas.com.mx/internal/model"

        "comprarmas.com.mx/internal/database"
)

var DatProduct = new(model.DataProductSearch)

func RegisterProduct(sName string, sSku string, sQuantity string, sPrice string) (string, error) {

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
        sSQL := "INSERT INTO Product(Name, Sku, Quantity, Price) "
        sSQL += "values('" + sName + "', '" + sSku + "', " + sQuantity + ", " + sPrice + ") "

        resDB, err := db.Exec(sSQL)
        if err != nil {
                //log.Println("No fue posible registrar el producto: " + err.Error())
                err = errors.New("No fue posible registrar el producto: " + err.Error())
                return "", err
        } else {
                lastId, errDB := resDB.LastInsertId()
                log.Println(resDB)
                return fmt.Sprint(lastId), errDB
        }
} //END RegisterProduct

func RegisterUsrProduct(sIdProd string, sIdUsr string) error {
        dtFechaHoy := time.Now()
        sTmp := dtFechaHoy.String()
        sFA := sTmp[0:19]

        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("RegisterUsrProduct_Error " + err.Error())
                return err
        }
        defer db.Close()

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("RegisterUsrProduct_Error_Ping " + err.Error())
                return err
        }

        defer db.Close()
        sSQL := "INSERT INTO UsrProduct(id_Usr, id_Product, RegisterDate) "
        sSQL += "values(" + sIdUsr + ", '" + sIdProd + "', '" + sFA + "') "

        resDB, err := db.Exec(sSQL)
        if err != nil {
                //log.Println("No fue posible registrar el producto: " + err.Error())
                err = errors.New("No fue posible registrar el producto asociado al usuario: " + err.Error())
                return err
        } else {
                //lastId, errDB := resDB.LastInsertId()
                log.Println(resDB)
                return nil
        }
} //END RegisterUsrProduct

func ListProducts(sId string, sOpt string) ([]map[string]interface{}, []string, error) {
        var aList []string
        var allData []map[string]interface{}

        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("GetListProducts_Error cnn DB " + err.Error())
                return nil, aList, err
        }
        defer db.Close()

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("GetListProducts_Error_Ping " + err.Error())
                return nil, aList, err
        }

        defer db.Close()

        sSQL := "SELECT DISTINCT c.id_Product, a.id_Usr, c.id_Product, a.Opt, c.Name, c.Quantity, c.Price, c.Sku, c.St FROM Usr a "
        sSQL += "INNER JOIN UsrProduct b ON (a.id_Usr=b.id_Usr) "
        sSQL += "INNER JOIN Product c ON (b.id_Product=c.id_Product) "
        sSQL += "WHERE (c.St = 'ACTIVE') "
        if sOpt != "admin" {
                sSQL += "AND (a.id_Usr = " + sId + ") "
        }

        resDat, err := db.Query(sSQL)
        if err != nil {
                return nil, aList, err
        } else {
                for resDat.Next() {
                        var datList = model.DataListProduct{}
                        resDat.Scan(
                                &datList.IdP, &datList.IdUsr, &datList.IdProd, &datList.Opt, &datList.Name,
                                &datList.Quantity, &datList.Price, &datList.Sku, &datList.St,
                        )
                        aList = append(aList, fmt.Sprint(datList.IdP), fmt.Sprint(datList.IdUsr), datList.Opt, datList.Name, fmt.Sprint(datList.Quantity), datList.Sku, datList.St)
                        mapJson := make(map[string]interface{})
                        mapJson["id"] = fmt.Sprint(datList.IdP)
                        mapJson["id_usr"] = fmt.Sprint(datList.IdUsr)
                        mapJson["id_product"] = fmt.Sprint(datList.IdProd)
                        mapJson["opt"] = datList.Opt
                        mapJson["name"] = datList.Name
                        mapJson["sku"] = datList.Sku
                        mapJson["quantity"] = fmt.Sprint(datList.Quantity)
                        mapJson["price"] = fmt.Sprint(datList.Price)
                        mapJson["st"] = fmt.Sprint(datList.St)
                        allData = append(allData, mapJson)
                }
        }

        sJson, err := json.Marshal(allData)
        fmt.Println(string(sJson))
        return allData, aList, nil
} //END ListProducts

func ExistProduct(sName string, sSku string) error {
        blnOK := false

        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("ExistProduct_Error cnn DB " + err.Error())
                return err
        }
        defer db.Close()

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("ExistProduct_Error_Ping " + err.Error())
                return err
        }

        defer db.Close()

        sSQL := "SELECT DISTINCT c.id_Product, c.Name, c.Quantity, c.Price, c.Sku, c.St FROM Product c "
        sSQL += "WHERE (c.St = 'ACTIVE') "
        if sName != "" {
                sSQL += "AND (c.Name = '" + sName + "') "
        }
        if sSku != "" {
                sSQL += "AND (c.sKu = '" + sSku + "') "
        }
        resDat, err := db.Query(sSQL)
        if err != nil {
                return err
        } else {
                for resDat.Next() {
                        var datList = model.DataListProduct{}
                        resDat.Scan(
                                &datList.IdP, &datList.Name, &datList.Quantity, datList.Price, &datList.Sku, &datList.St,
                        )
                        log.Println(&datList.Name)
                        blnOK = true
                }
        }
        if blnOK == false {
                return nil
        } else {
                err = errors.New("El producto ya existe.")
                return err
        }
} //END ExistProduct

func CancelProduct(sId string) error {
        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("CancelProduct_Error cnn DB " + err.Error())
                return err
        }
        defer db.Close()

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("CancelProduct_Error_Ping " + err.Error())
                return err
        }

        defer db.Close()
        resDB, err := db.Exec("UPDATE Product SET St = 'CANCEL' WHERE id_Product = ?", sId)
        if err != nil {
                log.Println("No fue posible cancelar el producto: " + err.Error())
                err = errors.New("No fue posible cancelar el producto: " + err.Error())
                return err
        } else {
                log.Println(resDB)
                return nil
        }
} //END CancelProduct

func ListProductsSearch(sIdUsr string, sName string, sSku string, sPI string, sPF string) ([]map[string]interface{}, error) {
        var allData []map[string]interface{}

        db, err := database.CnnDB_ENV()
        if err != nil {
                log.Println("ListProductsSearch_Error cnn DB " + err.Error())
                return nil, err
        }
        defer db.Close()

        // Ahora vemos si tenemos conexión
        err = db.Ping()
        if err != nil {
                log.Println("ListProductsSearch_Error_Ping " + err.Error())
                return nil, err
        }

        defer db.Close()

        sSQL := "SELECT DISTINCT c.id_Product, a.id_Usr, c.id_Product, a.Opt, c.Name, c.St, c.Quantity, c.Price, c.Sku "
        sSQL += "FROM Usr a "
        sSQL += "INNER JOIN UsrProduct b ON (a.id_Usr=b.id_Usr) "
        sSQL += "INNER JOIN Product c ON (b.id_Product=c.id_Product) "
        sSQL += "WHERE (c.St = 'ACTIVE') "
        if sIdUsr != "" {
                sSQL += "AND (a.id_Usr = " + sIdUsr + ") "
        }
        if sName != "" {
                sSQL += "AND (c.Name LIKE '" + sName + "%') "
        }
        if sSku != "" {
                sSQL += "AND (c.sKu LIKE '" + sSku + "%') "
        }
        if sPI != "" {
                sSQL += "AND (c.Price >= " + sPI + ") "
        }
        if sPF != "" {
                sSQL += "AND (c.Price <= " + sPF + ") "
        }

        resDat, err := db.Query(sSQL)
        if err != nil {
                return nil, err
        } else {
                for resDat.Next() {
                        var datList = model.DataListProduct{}
                        resDat.Scan(
                                &datList.IdP, &datList.IdUsr, &datList.IdProd, &datList.Opt, &datList.Name,
                                &datList.St, &datList.Quantity, &datList.Price, &datList.Sku,
                        )
                        mapJson := make(map[string]interface{})
                        mapJson["id"] = fmt.Sprint(datList.IdP)
                        mapJson["id_usr"] = fmt.Sprint(datList.IdUsr)
                        mapJson["id_product"] = fmt.Sprint(datList.IdProd)
                        mapJson["opt"] = datList.Opt
                        mapJson["name"] = datList.Name
                        mapJson["sku"] = datList.Sku
                        mapJson["quantity"] = fmt.Sprint(datList.Quantity)
                        mapJson["price"] = fmt.Sprint(datList.Price)
                        mapJson["st"] = fmt.Sprint(datList.St)
                        allData = append(allData, mapJson)
                }
        }

        sJson, err := json.Marshal(allData)
        fmt.Println(string(sJson))
        return allData, nil
} //END ListProductsSearch
