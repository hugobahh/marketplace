package database

import (
        "fmt"
        "log"
        "strconv"
        "time"

        "context"

        "go.mongodb.org/mongo-driver/mongo"
        "go.mongodb.org/mongo-driver/mongo/options"

        "database/sql"

        "comprarmas.com.mx/internal/secrets"

        _ "github.com/go-sql-driver/mysql"

        "github.com/go-redis/redis"
)

func CnnMongo() (cClient *mongo.Client, e error) {
        sCnnDB := stringCnnMongo()
        //sCnnDB := "mongodb://127.0.0.1:27017/"
        //sCnnDB := "mongodb://docker:mongopw@172.24.100.14:27017"
        client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(sCnnDB))
        fmt.Println(client)
        if err != nil {
                log.Println("CnnMongo_Err: " + err.Error())
                panic(err)
        } else {
                fmt.Println("Conexion realizada")

                err = client.Ping(context.TODO(), nil)
                if err != nil {
                        log.Println("CnnMongo_ping_Err: " + err.Error())
                        return nil, err
                }

                return client, err
        }
} //FIN de cnnMongo

//=====================  Archivo Conf ==========================================
func stringCnnMongo() string {
        sIP := secrets.LoadSecrets("MONGO_HOST")
        sPort := secrets.LoadSecrets("MONGO_PORT")
        sUsr := secrets.LoadSecrets("MONGO_USER")
        sPWD := secrets.LoadSecrets("MONGO_PWD")

        //mongodb://127.0.0.1:27017/
        //mongodb://docker:mongopw@172.24.100.14:27017
        //sCnn := "mongodb://" + sIp + ":" + fmt.Sprintf(sPort) + "/"
        sCnn := "mongodb://" + sUsr + ":" + sPWD + "@" + sIP + ":" + fmt.Sprintf(sPort)

        return sCnn
} //FIN stringCnnMongo

func Close_ClientCtx(client *mongo.Client, ctx context.Context) {
        defer func() {
                if err := client.Disconnect(ctx); err != nil {
                        panic(err)
                }
        }()
}

func Close(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {
        defer cancel()

        defer func() {
                if err := client.Disconnect(ctx); err != nil {
                        panic(err)
                }
        }()
}

//cnn MYSQL
func CnnDB_ENV() (db *sql.DB, e error) {
        sPwd := ""
        sIP := ""

        sIP = secrets.LoadSecrets("DB_HOST")
        sPort := secrets.LoadSecrets("DB_PORT")
        sUsr := secrets.LoadSecrets("DB_USER")
        sTmp := secrets.LoadSecrets("DB_PWD2")
        if sTmp != "" {
                sPwd = secrets.LoadSecrets("DB_PWD1") + "#" + secrets.LoadSecrets("DB_PWD2")
        } else {
                sPwd = secrets.LoadSecrets("DB_PWD1")
        }
        sDB := secrets.LoadSecrets("DB_NAME")

        db, err := sql.Open("mysql", sUsr+":"+sPwd+"@tcp("+sIP+":"+sPort+")/"+sDB+"?parseTime=true")
        if err != nil {
                //log.Println(err.Error())
                return nil, err
        }
        return db, nil
}

//=========================================================
func CnnDB_Select() (db *sql.DB, e error) {
        sPwd := ""
        sIP := ""

        sIP = secrets.LoadSecrets("DB_HOST_2")
        sPort := secrets.LoadSecrets("DB_PORT_2")
        sUsr := secrets.LoadSecrets("DB_USER_2")
        sTmp := secrets.LoadSecrets("DB_PWD2_2")
        if sTmp != "" {
                sPwd = secrets.LoadSecrets("DB_PWD1_2") + "#" + secrets.LoadSecrets("DB_PWD2_2")
        } else {
                sPwd = secrets.LoadSecrets("DB_PWD1_2")
        }
        sDB := secrets.LoadSecrets("DB_NAME_2")

        db, err := sql.Open("mysql", sUsr+":"+sPwd+"@tcp("+sIP+":"+sPort+")/"+sDB+"?parseTime=true")
        if err != nil {
                //log.Println(err.Error())
                return nil, err
        }
        return db, nil
}

//=========================================================
//=========================================================
func CnnDB_MSSQL() (db *sql.DB, e error) {
        //sPort := "1433"
        sPwd := ""
        sServer := secrets.LoadSecrets("MSSQL_HOST")
        //sPort := secrets.LoadSecrets("MSSQL_PORT")
        //sUsr := secrets.LoadSecrets("MSSQL_USER")
        //sDB := secrets.LoadSecrets("MSSQL_DB")
        sTmp := secrets.LoadSecrets("MSSQL_PWD2")
        if sTmp != "" {
                sPwd = secrets.LoadSecrets("MSSQL_PWD1") + "#" + secrets.LoadSecrets("DB_PWD2_2")
        } else {
                sPwd = secrets.LoadSecrets("MSSQL_PWD1")
        }

        //sCnnString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d",
        //      sServer, sUsr, sPwd, sPort)

        sCnnString := "sqlserver://sa:" + sPwd + "@" + sServer + ":1433?database=master&connection+timeout=30"
        // Create connection pool
        db, err := sql.Open("mssql", sCnnString)
        if err != nil {
                //log.Println(err.Error())
                log.Println("CnnDB_MSSQL_Err:" + err.Error())
                defer db.Close()
                return nil, err
        }
        return db, nil
}

//cnn MSSQL
func CnnDB_MSSQL_TEST() (db *sql.DB, e error) {
        //println(sql.Drivers())
        sPwd := ""
        sServer := secrets.LoadSecrets("MSSQL_HOST")
        //sPort := secrets.LoadSecrets("MSSQL_PORT")
        //sUsr := secrets.LoadSecrets("MSSQL_USER")
        sDB := secrets.LoadSecrets("MSSQL_DB")
        sPwd = secrets.LoadSecrets("MSSQL_PWD1")

        cnnString := "server=" + sServer + ";Initial Catalog=" + sDB + ";user id=sa;password=" + sPwd
        //"server=14f0d4f.online-server.cloud;Initial Catalog=PAGINAR;user id=sa;password=Doqsoft#2023"
        db, err := sql.Open("mssql", cnnString)
        err = db.Ping()
        if err != nil {
                //fmt.Println("Failed to ping: ", err.Error())
                log.Println("CnnDB_MSSQL_TEST_Err:" + err.Error())
        }

        //rows, err := db.Query("select * from Estructura where id_Estructura=?", 2)
        //if err != nil {
        //      fmt.Println(err)
        //}
        return db, nil
}

func CnnDB() (db *sql.DB, e error) {
        sDbSvr := secrets.LoadSecrets("DB_SERVER")
        if sDbSvr == "MSSQL" {
                dbOK, err := CnnDB_MSSQL_TEST()
                log.Println(dbOK)
                if err != nil {
                        log.Println(err.Error())
                }
                return dbOK, nil
        }
        dbOK, err := CnnDB_ENV()
        log.Println(dbOK)
        if err != nil {
                log.Println(err.Error())
        }
        return dbOK, nil
}

func RedisCnn_ENV() *redis.Client {
        sIP := secrets.LoadSecrets("REDIS_HOST")
        sPort := secrets.LoadSecrets("REDIS_PORT")
        //sUsr := utils.LoadSecrets("REDIS_USER")
        sPwd := secrets.LoadSecrets("REDIS_PWD")
        //sDB := utils.LoadSecrets("DB_NAME")

        client := redis.NewClient(&redis.Options{
                Addr:     sIP + ":" + sPort,
                Password: sPwd,
                DB:       10,
        })
        return client
}

func Ping(client *redis.Client) error {
        pong, err := client.Ping().Result()
        if err != nil {
                return err
        }
        fmt.Println(pong, err)
        // Output: PONG <nil>
        return nil
}

func Set(client *redis.Client, sKey string, sValue string) error {
        err := client.Set(sKey, sValue, 0).Err()
        if err != nil {
                return err
        }
        return nil
}

func Get(client *redis.Client, sKey string) (blnOK bool, err error) {
        nameVal, err := client.Get(sKey).Result()
        if nameVal != "" {
                return true, nil
        }
        if err != nil {
                return false, err
        }
        return true, nil
}

func AddExpTime(client *redis.Client, sKey string) (blnOK bool, err error) {
        nT := secrets.LoadSecrets("EXP_TIME")
        nH, _ := strconv.ParseInt(nT, 10, 64)
        nTmp := int32(nH)
        //yourTime := Int31n(nTmp)
        _, err = client.Expire(sKey, time.Duration(nTmp)*time.Hour).Result()
        if err != nil {
                return false, err
        }
        return true, nil
}
