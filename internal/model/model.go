package model

type DataUsrReg struct {
        Mail string `json:"mail"`
        Pwd  string `json:"pwd"`
}

type DataUsrGet struct {
        IdUsr string `json:"id_usr"`
        Mail  string `json:"mail"`
        Pwd   string `json:"pwd"`
        St    string `json:"estatus"`
        Opt   string `json:"opt"`
}

type DataListProduct struct {
        IdP      int32   `json:"id"`
        IdUsr    int32   `json:"id_usr"`
        IdProd   int32   `json:"id_product"`
        Opt      string  `json:"opt"`
        Name     string  `json:"name"`
        Sku      string  `json:"sku"`
        Quantity int32   `json:"quantity"`
        Price    float32 `json:"price"`
        St       string  `json:"estatus"`
}

type DataUsr struct {
        IdUsr string `json:"id_usr"`
        Mail  string `json:"mail"`
        Opt   string `json:"opt"`
        St    string `json:"estatus"`
}

type DataProduct struct {
        Name     string `json:"name"`
        Sku      string `json:"sku"`
        Quantity string `json:"quantity"`
        Price    string `json:"price"`
        IdUsr    string `json:"id_usr"`
        Mail     string `json:"mail"`
}

type DataProductSearch struct {
        Id   string `json:"id"`
        Name string `json:"name"`
        Sku  string `json:"sku"`
        PI   string `json:"pi"`
        PF   string `json:"pf"`
}
