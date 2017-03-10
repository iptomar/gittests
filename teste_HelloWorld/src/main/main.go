package main
import(
	"net/http"
	"github.com/gorilla/mux"
	"text/template"
	"os"
	"log"
	"bufio"
	"strings"
)

func main() {
	webService()
}

var staticPages = populateStaticPages()

func webService(){
	gorillaRoute := mux.NewRouter()

	gorillaRoute.HandleFunc("/", page)
	gorillaRoute.HandleFunc("/{page_alias}", page)
	
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/js/", serveResource)

	http.Handle("/", gorillaRoute)
	http.ListenAndServe(":5000", nil)
}

func page (w http.ResponseWriter, r *http.Request){
	urlParams := mux.Vars(r)
	page_alias := urlParams["page_alias"]
	if page_alias == ""{
		page_alias = "home"
	}
	layout_page := staticPages.Lookup(page_alias+".html")
	if layout_page == nil{
		layout_page = staticPages.Lookup("404.html")
		w.WriteHeader(404)
	}
	layout_page.Execute(w, nil)

	//w.Write([]byte("Hello World"))
}

//---------------------------------------------------------------------------
// Retrieve all files under a given folder
func populateStaticPages() *template.Template{
	result := template.New("templates")
	templatePaths := new([]string)

	basePath := "pages"
	templateFolder, _:= os.Open(basePath)
	defer templateFolder.Close()
	templatePathsRaw, _:= templateFolder.Readdir(-1)

	for _, pathInfo := range templatePathsRaw{
		log.Println(pathInfo.Name())
		*templatePaths = append(*templatePaths, basePath+"/"+pathInfo.Name())
	}
	result.ParseFiles (*templatePaths...)
	return result
}
//----------------------------------------------------------------------------

//---------------------------------------------------------------------------
// Serve resources of type js, css and img files
func serveResource (w http.ResponseWriter, r *http.Request){
	path := "public/"+ r.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css"){
		contentType = "text/css; charset=utf-8"
	} else if strings.HasSuffix(path, ".png"){
		contentType = "image/png; charset=utf-8"
	} else if strings.HasSuffix(path, ".jpg"){
		contentType = "image/jpg; charset=utf-8"
	} else if strings.HasSuffix(path, ".js"){
		contentType = "application/javascript; charset=utf-8"
	} else {
		contentType = "text/plain; charset=utf-8"
	}
	log.Println(path)
	f, err := os.Open(path)
	if err == nil{
		defer f.Close()
		w.Header().Add("Content-Type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else{
		w.WriteHeader(404)
	}
}
//---------------------------------------------------------------------------