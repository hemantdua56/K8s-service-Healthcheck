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
        // c "healthcheck/src/config"
        "github.com/spf13/viper"
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
      Build string
      HealthStatus string

}


func main() {

    // Set the file name of the configurations file
    viper.SetConfigName("config")

    // Set the path to look for the configurations file
    viper.AddConfigPath(".")

    // Enable VIPER to read Environment Variables
    viper.AutomaticEnv()

    viper.SetConfigType("yml")  
    
    if err := viper.ReadInConfig(); err != nil {
      log.Printf("Error reading config file, %s", err)
    }

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
       var urls = viper.GetStringSlice("urls")

       var temp = [1][]Data{}
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
             dt = append(dt, Diff{Microservice: temp[0][i].Microservice,Build: temp[0][i].Build, HealthStatus: temp[0][i].Healthcode })
       }

       log.Println(dt)
       templ.ExecuteTemplate(w, "submitted.html", dt)

 }
