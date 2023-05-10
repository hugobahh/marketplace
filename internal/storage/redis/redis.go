package redis

import (
        "log"

        "comprarmas.com.mx/internal/database"
)

func ChkUsr(sKey string) error {

        //Redis cliente
        clientRedis := database.RedisCnn_ENV()
        // check connection status
        err := database.Ping(clientRedis)
        if err != nil {
                log.Println(err)
                return err
        }

        blnOK, err := database.Get(clientRedis, sKey)
        if blnOK == false {
                log.Println(err)
                return err
        }

        //c.JSON(200, gin.H{"code": "ChkUsr", "message": "OK.", "status_code": 200})
        return nil
} //FIN de ChkUsr

//=======================================================
func AddUsr(sKey string) error {
        //Redis cliente
        clientRedis := database.RedisCnn_ENV()
        // check connection status
        err := database.Ping(clientRedis)
        if err != nil {
                log.Println(err)
                //c.JSON(400, gin.H{"code": "SearchToken", "message": "Error al hacer ping a redis.", "status_code": 400})
                return err
        }
        //Agrega una clave-valor
        err = database.Set(clientRedis, sKey, "ACTIVO")
        if err != nil {
                log.Println(err)
                return nil
        }

        //Agregar el tiempo de expiracion 1h
        _, err = database.AddExpTime(clientRedis, sKey)
        return nil
} //FIN de AddUsr
