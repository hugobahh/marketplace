
package controller

import (
        "bytes"
        "crypto/rand"
        "crypto/rsa"
        "encoding/json"
        "errors"
        "fmt"
        "log"
        "os"

        "comprarmas.com.mx/internal/model"
        "comprarmas.com.mx/internal/storage/product"
        "comprarmas.com.mx/internal/storage/redis"
        "comprarmas.com.mx/internal/storage/usr"

        "github.com/gofiber/fiber/v2"
)

var (
        privateKey *rsa.PrivateKey
)

func PostLogin(ctx *fiber.Ctx) error {
        log.Println("LoginUsr ...")
        var err error
        sId := ""
        sOpt := ""
        log.Println(sId + sOpt)

        file, _ := os.Create("cb2.log")
        log.SetOutput(file)
        defer file.Close()

        //sMail := ctx.FormValue("Email")
        //sPwd := ctx.FormValue("password")
        //log.Println(sMail, sPwd)
        log.Println("Login_1")
        if ctx.Get("Content-Type") != "application/json" {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": err.Error(),
                        "msg":   "Content-Type header is not application/json.",
                })
        }
        log.Println("Login_2")
        var datUsr = model.DataUsrReg{}
        eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
        //if eJson != nil {
        //      return err
        //}
        log.Println("Login_3")
        err = eJson.Decode(&datUsr)
        if err != nil {
                return err
        }

        log.Println("Login_4")
        //Autenticate
        sId, sOpt, err = usr.LoginUsr(datUsr.Mail, datUsr.Pwd)
        if err != nil {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": err.Error(),
                        "id":    "",
                })
        } else {
                log.Println("Login_5")
                //Register Redis
                err := redis.AddUsr("Id_" + sId)
                if err != nil {
                        return ctx.Status(401).JSON(fiber.Map{
                                "Error": err.Error(),
                                "msg":   "No fue posible registrar al usuario en redis.",
                        })
                }
                log.Println("Login_6")
                return ctx.Status(200).JSON(fiber.Map{
                        "OK":  "Access",
                        "id":  sId,
                        "opt": sOpt,
                })
        }
}

func PostRegisterUsr(ctx *fiber.Ctx) error {
        log.Println("PostRegisterUsr ...")
        rng := rand.Reader
        var err error
        sToken := ""
        sId := ""
        log.Println(sToken + sId)

        privateKey, err = rsa.GenerateKey(rng, 2048)
        if err != nil {
                log.Fatalf("rsa.GenerateKey: %v", err)
        }

        if ctx.Get("Content-Type") != "application/json" {
                //err.Error() := "Content-Type header is not application/json"
                return err
        }

        var datUsr = model.DataUsrReg{}
        eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
        //if eJson != nil {
        //      return err
        //}
        err = eJson.Decode(&datUsr)
        if err != nil {
                return err
        }
        //Buscar si ya existe el usuario
        err = usr.GetUsr(datUsr.Mail, "")
        if err != nil {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": err.Error(),
                        "id":    "0",
                })
        }
        //Registrar el Usuario
        sId, err = usr.RegisterUsr(datUsr.Mail, datUsr.Pwd)
        if err != nil {
                return ctx.JSON(fiber.Map{
                        "Error": err.Error(),
                        "id":    "0",
                })
        } else {
                // Buscar el dato completo que se registro
                return ctx.JSON(fiber.Map{
                        "OK": "Usuario registrado",
                        "id": sId,
                })
        }
} //END PostRegisterUsr

func GetListProducts(ctx *fiber.Ctx) error {
        var err error

        if ctx.Get("Content-Type") != "application/json" {
                err = errors.New("Content-Type header is not application/json")
                return err
        }
        log.Println(ctx.Body())
        var datUsr = model.DataUsr{}
        eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
        err = eJson.Decode(&datUsr)
        if err != nil {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": err.Error(),
                })
        }
        //CHK Redis
        err = redis.ChkUsr("Id_" + datUsr.IdUsr)
        if err != nil {
                return ctx.JSON(fiber.Map{
                        "Error": err.Error(),
                        "msg":   "Se ha terminado la sesión del usuario.",
                })
        }

        //Obtener la lista de productos
        mapJson, aChk, errLP := product.ListProducts(datUsr.IdUsr, datUsr.Opt)
        if errLP != nil {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": errLP.Error(),
                })
        } else {
                sJson, err := json.Marshal(mapJson)
                if err != nil {
                        return ctx.Status(401).JSON(fiber.Map{
                                "Error": err.Error(),
                        })
                }
                log.Println(string(sJson))
                log.Println(aChk)
                return ctx.JSON(string(sJson))
        }
        //return ctx.JSON(aChk)
}

func PostCancelProduct(ctx *fiber.Ctx) error {
        log.Println("RegisterUsr ...")
        sId := ctx.Params("id")
        var err error

        //CHK Redis
        //err = redis.ChkUsr("Id_" + sId)
        //if err != nil {
        //      return ctx.Status(401).JSON(fiber.Map{
        //              "Error": err.Error(),
        //              "msg":   "Se ha terminado la sesión del usuario.",
        //      })
        //}

        err = product.CancelProduct(sId)
        if err != nil {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": err.Error(),
                        "mag":   "No fue posible cancelar el producto.",
                })
        } else {
                // Buscar el dato completo que se registro
                return ctx.Status(200).JSON(fiber.Map{
                        "process": "Cancelar.",
                        "msg":     "Regisro cancelado.",
                })
        }
}

func PostRegProduct(ctx *fiber.Ctx) error {
        log.Println("PostRegProductLogin ...")
        var err error

        if ctx.Get("Content-Type") != "application/json" {
                err = errors.New("Content-Type header is not application/json")
                return err
        }

        var datProd = model.DataProduct{}
        eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
        //if eJson != nil {
        //      return err
        //}
        err = eJson.Decode(&datProd)
        if err != nil {
                return err
        }
        //Buscar si el usuario esta en session
        err = redis.ChkUsr(datProd.IdUsr)
        if err != nil {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": err.Error(),
                        "msg":   "Se ha terminado la sesión del usuario.",
                })
        }

        err = product.ExistProduct(datProd.Name, datProd.Sku)
        if err != nil {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": err.Error(),
                        "msg":   "El producto ya existe.",
                })
        }

        //Register product
        sId, err := product.RegisterProduct(datProd.Name, datProd.Sku, fmt.Sprint(datProd.Quantity), fmt.Sprint(datProd.Price))
        if err != nil {
                return ctx.Status(202).JSON(fiber.Map{
                        "Error": "No fue registrado el producto",
                        "msg":   "El producto ya existe",
                })
        } else {
                err = product.RegisterUsrProduct(sId, fmt.Sprint(datProd.IdUsr))
                if err != nil {
                        //err = errors.New("No fue posible registrar el producto.")
                        return ctx.Status(401).JSON(fiber.Map{
                                "Error": err.Error(),
                                "msg":   "No fue posible registrar el producto.",
                        })
                }
                return ctx.Status(200).JSON(fiber.Map{
                        "OK":  "Access",
                        "msg": "El producto se registro correctamente.",
                })
        }
} //END PostRegProduct

func PostProductSearch(ctx *fiber.Ctx) error {
        log.Println("PostRegProductLogin ...")
        var err error

        if ctx.Get("Content-Type") != "application/json" {
                err = errors.New("Content-Type header is not application/json")
                return err
        }

        var datProd = model.DataProductSearch{}
        eJson := json.NewDecoder(bytes.NewReader(ctx.Body()))
        //if eJson != nil {
        //      return err
        //}
        err = eJson.Decode(&datProd)
        if err != nil {
                return err
        }
        //Buscar si el usuario esta en session
        //err = redis.ChkUsr("")
        //if err != nil {
        //      return ctx.JSON(fiber.Map{
        //              "Error": err.Error(),
        //              "msg":   "Se ha terminado la sesión del usuario.",
        //      })
        //}
        mapJson, err := product.ListProductsSearch(datProd.Id, datProd.Name, datProd.Sku, datProd.PI, datProd.PF)
        if err != nil {
                return ctx.Status(401).JSON(fiber.Map{
                        "Error": err.Error(),
                        "msg":   "No fue posible encontrar productos.",
                })
        } else {
                sJson, err := json.Marshal(mapJson)
                if err != nil {
                        return ctx.Status(401).JSON(fiber.Map{
                                "Error": err.Error(),
                        })
                }
                log.Println(string(sJson))
                return ctx.Status(200).JSON(string(sJson))
        }
} //END PostProductSearch

func PostProductAdmin(ctx *fiber.Ctx) error {
        log.Println("PostRegProductLogin ...")
        var err error

        //Buscar si el usuario esta en session
        //err = redis.ChkUsr("")
        //if err != nil {
        //      return ctx.JSON(fiber.Map{
        //              "Error": err.Error(),
        //              "msg":   "Se ha terminado la sesión del usuario.",
        //      })
        //}
        mapJson, err := usr.GetUsrsSeller()
        if err != nil {
                return ctx.JSON(fiber.Map{
                        "Error": err.Error(),
                        "msg":   "No fue posible encontrar productos.",
                })
        } else {
                sJson, err := json.Marshal(mapJson)
                if err != nil {
                        return ctx.JSON(fiber.Map{
                                "Error": err.Error(),
                        })
                }
                log.Println(string(sJson))
                return ctx.JSON(string(sJson))
        }
}
