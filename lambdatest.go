package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
    "html"
)

type Settings struct {
    ClientId string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
    Port string `json:"port"`
}

type ImgurData struct {
    Data []struct {
        AccountID      interface{} `json:"account_id"`
        AccountURL     interface{} `json:"account_url"`
        Animated       bool        `json:"animated"`
        Bandwidth      int         `json:"bandwidth"`
        CommentCount   int         `json:"comment_count"`
        CommentPreview interface{} `json:"comment_preview"`
        Datetime       int         `json:"datetime"`
        Description    interface{} `json:"description"`
        Downs          int         `json:"downs"`
        Favorite       bool        `json:"favorite"`
        Height         int         `json:"height"`
        ID             string      `json:"id"`
        IsAlbum        bool        `json:"is_album"`
        Link           string      `json:"link"`
        Nsfw           bool        `json:"nsfw"`
        Points         int         `json:"points"`
        Score          int         `json:"score"`
        Section        string      `json:"section"`
        Size           int         `json:"size"`
        Title          string      `json:"title"`
        Topic          string      `json:"topic"`
        TopicID        int         `json:"topic_id"`
        Type           string      `json:"type"`
        Ups            int         `json:"ups"`
        Views          int         `json:"views"`
        Vote           interface{} `json:"vote"`
        Width          int         `json:"width"`
    } `json:"data"`
    Status  int  `json:"status"`
    Success bool `json:"success"`
}


func main() {
    configFile, err := os.Open("config.json")
    if err != nil {
        fmt.Println("opening config file", err.Error())
    }

    settings := Settings{}
    jsonParser := json.NewDecoder(configFile)
    if err = jsonParser.Decode(&settings); err != nil {
        fmt.Println("parsing config file", err.Error())
    }

    fmt.Printf("%s %s %s", settings.ClientId, settings.ClientSecret, settings.Port)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        x, err := getRandomImgur(settings, r.URL.Query().Get("text"))
        if err != nil {
            fmt.Println(err)
        }

        var imgurData ImgurData
        err = json.Unmarshal(x, &imgurData)
        if err != nil {
            return
        }
        fmt.Println(imgurData.Data[0].Link)
        fmt.Fprintf(w, imgurData.Data[0].Link)
    })

    log.Fatal(http.ListenAndServe(settings.Port, nil))

}

func getRandomImgur(settings Settings, searchString string) (result []byte, err error) {
    client := &http.Client{}

    req, err := http.NewRequest("GET", "https://api.imgur.com/3/gallery/search/top/?q=" + html.EscapeString(searchString), nil)
    clientId := "Client-ID " + settings.ClientId
    req.Header.Add("Authorization", clientId)
    resp, err := client.Do(req)
    body, err := ioutil.ReadAll(resp.Body)
    return body, err
}

