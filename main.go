package main
import (
	"net/http"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"time"
	"sync"
	"log"
)
func main() {
	fmt.Println("starting server")
	server()
}
type List struct {
	Websites []string `json:"websites"`
 }
 
var Websites = make(map[string]string)
func server(){
    
	go func() {
		for {
			monitor()
		   time.Sleep(5 * time.Second)
		}
	 }()
  
    http.HandleFunc("/POST",WebsitesPostHandler)
	http.HandleFunc("/GET",GetStatus)
	http.HandleFunc("/CHECK",GetSingleHandler)
	http.ListenAndServe("localhost:4000", nil)
	
}

func WebsitesPostHandler(w http.ResponseWriter, r *http.Request) {
	
	//fmt.Fprintf(w, "Welcome to the websites server!\n")
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading JSON ResponseWriter")
	   return
	}
	var list List
	err = json.Unmarshal(body, &list)
	if err != nil {
		log.Printf("Error unmarshalling JSON Post")
	   return
	}
	var Added_to_Map []string;
    flag:=false
	for _, website := range list.Websites {
	   if _, ok := Websites[website]; !ok {
		  Websites[website] = ""
		  Added_to_Map=append(Added_to_Map,website)
		  flag=true
	   }
	}
	if flag {
	fmt.Println("We have successfully Added these Websites to MAP", Added_to_Map);
	}
	//fmt.Print(Websites)
}
func GetStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Websites)
	if err != nil {
	   log.Println(err)
	   return
	}
	
 }
 func GetSingleHandler(w http.ResponseWriter, r *http.Request) {
	website := r.URL.Query().Get("name")
	status, ok := Websites[website]
	if !ok {
	   log.Print("Website not present")
	   return
	}
	EnteredWesite := map[string]string{
	   website: status,
	}
	err := json.NewEncoder(w).Encode(EnteredWesite)
	if err != nil {
	   log.Fatal(err)
	   return
	}
 
 }



func monitor() {	
	
    if len(Websites) == 0 {
		return
	 }
	 var wg sync.WaitGroup
	 var lock sync.Mutex
	 wg.Add(len(Websites))
	 for website := range Websites {
		website := website
		go func() {
		   defer wg.Done()
		   resp, err := http.Get("https://" + website)
		   if err != nil {
			Websites[website] = "DOWN"
		   } else {
			  lock.Lock()
			  if(resp.StatusCode==200){
			  Websites[website] = "UP"}
			//   resp.StatusCode
			  lock.Unlock()
		  }
		}()
	 }
	 wg.Wait()
	//  fmt.Println("Monitoring server");
    //  fmt.Println("Showing websites status",Websites);
}
