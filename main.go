package main
import (
        "time"
        "log"
        "io/ioutil"
        "crypto/tls"
        "net/http"
        "sync"
        "html/template"
        "encoding/json"
    )
  var (
        err error
        templ *template.Template
)

type Data struct {
      Microservice string `json:"Microservice"`
      Build string `json:"Build"`
      Healthcode string `json:"Healthcode"`
      Healthstatus string `json:"Healthstatus"`

}

type Diff struct{
      Microservice string
      BsInt string

}


func main() {

    templ, err = templ.ParseGlob("templates/*.html")
    log.Println(templ)
    if err != nil {
    log.Println(err)
  	}

    http.HandleFunc("/", hello)

    log.Printf("Starting server for testing HTTP POST...\n")
    if err := http.ListenAndServe(":8000", nil); err != nil {
        log.Fatal(err)
     }

 }

 var myClient = &http.Client{Timeout: 15 * time.Second}


 func hello(w http.ResponseWriter, r *http.Request){
       if r.URL.Path != "/" {
         http.Error(w, "404 not found.", http.StatusNotFound)
         return
       }
       //var dt []Diff
       temp := [<NUMBER of ENVIRONEMTNES>][]Data{}
       urls := []string{"https://<ENV_URL>/healthcheck"}
       var wg sync.WaitGroup
       wg.Add(len(urls))

       for i := 0; i < len(urls); i++ {
            go func(i int) {
               defer wg.Done()
               url := urls[i]
               transCfg := &http.Transport{
               TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // ignore expired SSL certificates
               }
               myClient = &http.Client{Transport: transCfg}

               req, err := http.NewRequest(http.MethodGet,url,nil)
               if err != nil {
                   log.Fatal(err)
               }
               r, err := myClient.Do(req)
               if err != nil {
                   log.Fatal(err)
               }
               body, readErr := ioutil.ReadAll(r.Body)
               if readErr != nil {
                   log.Fatal(readErr)
               }
               jsonErr:= json.Unmarshal(body,&temp[i])
               if jsonErr != nil {
                 log.Printf("error decoding sakura response: %v", jsonErr)
                  if e, ok := jsonErr.(*json.SyntaxError); ok {
                      log.Printf("syntax error at byte offset %d", e.Offset)
                  }
                  log.Printf("sakura response: %q", body)
                  log.Println(url)
					result := make([]Data, 50)
                 	temp[i] = result
                  //return jsonErr
               }

               log.Println(temp[i])
              l:=len(temp[0])
              log.Println(l)

               }(i)
             }
             wg.Wait()
       var dt []Diff
       
       for i := 0; i < len(temp[0]); i++ {
             dt = append(dt, Diff{Microservice: temp[0][i].Microservice,<ENV NAME>: temp[0][i].Build, Health<ENV NAME>: temp[0][i].Healthcode })
       }

       log.Println(dt)
       templ.ExecuteTemplate(w, "submitted.html", dt)

 }
