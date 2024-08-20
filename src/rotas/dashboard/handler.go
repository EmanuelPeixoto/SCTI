package dashboard

import (
  DB "SCTI/database"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "html/template"
)

type Courses struct {
  Seg string `json:"seg"`
  Ter string `json:"ter"`
  Qua string `json:"qua"`
  Qui string `json:"qui"`
  Sex string `json:"sex"`
}

type DashboardData struct {
  IsVerified bool
  HTMLContent template.HTML
}

func GetDashboard(w http.ResponseWriter, r *http.Request) {
  cookie, err := r.Cookie("accessToken")
  if err != nil {
    // fmt.Println("Error Getting cookie:", err)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
    return
  }

  if cookie.Value == "-1" {
    // fmt.Println("Invalid accessToken")
    http.Redirect(w, r, "/login", http.StatusSeeOther)
  }

  MakeHTML();

  email := DB.GetEmail(cookie.Value)
  standing := DB.GetStanding(email)
  htmlContent := template.HTML(MakeHTML())

  data := DashboardData{
    IsVerified: standing,
    HTMLContent: htmlContent,
  }

  tmpl, err := template.ParseFiles("template/dashboard.gohtml")
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  tmpl.ExecuteTemplate(w, "dashboard", data)
}

func PostDashboard(w http.ResponseWriter, r *http.Request) {
  auth, err := r.Cookie("accessToken")
  if err != nil {
    // fmt.Println("Error Getting cookie:", err)
    http.Redirect(w, r, "/login", http.StatusSeeOther)
    return
  }

  if auth.Value == "-1" {
    // fmt.Println("Invalid accessToken")
    http.Redirect(w, r, "/login", http.StatusSeeOther)
  }
  var courses Courses
  if r.Header.Get("Content-Type") == "application/json" {
    err := json.NewDecoder(r.Body).Decode(&courses)
    if err != nil {
      log.Printf("Error decoding JSON: %v", err)
      http.Error(w, "Invalid JSON", http.StatusBadRequest)
      return
    }
  } else {
    // Lida com application/x-www-form-urlencoded
    if err := r.ParseForm(); err != nil {
      log.Printf("Error parsing form: %v", err)
      http.Error(w, "Invalid form data", http.StatusBadRequest)
      return
    }
    courses.Seg = r.FormValue("seg")
    courses.Ter = r.FormValue("ter")
    courses.Qua = r.FormValue("qua")
    courses.Qui = r.FormValue("qui")
    courses.Sex = r.FormValue("sex")
    fmt.Fprintf(w, "%s\n%s\n%s\n%s\n%s", courses.Seg, courses.Ter, courses.Qua, courses.Qui, courses.Sex)
  }
}

func MakeHTML()string{
  activities, err := DB.GetAllActivities()
  html := "<ul class=\"courses\">\n"
  if err != nil {
    fmt.Print(err.Error())
  } else {
    for _, a := range activities {
      html += fmt.Sprintf("<li class=\"atividades\"> ID: %v | hora: %v | sala: %v | descrição: %v | dia: %v | palestrante: %v | vagas: %v | tópico: %v tipo de atividade: %v <button class=\"cadastrar\" value=\"%v\" hx-post=\"/cadastrar\" hx-trigger=\"click\" hx-swap=\"none\">Cadastrar</button></li>\n",
      a.Activity_id,
      a.Time,
      a.Room,
      a.Description,
      a.Day,
      a.Speaker,
      a.Spots,
      a.Topic,
      a.Activity_type,
      a.Activity_id,
    )
  }
}
html += "</ul>"
return html
}

func RegisterRoutes(mux *http.ServeMux) {
  mux.HandleFunc("GET /dashboard", GetDashboard)
  mux.HandleFunc("POST /dashboard", PostDashboard)
  mux.HandleFunc("POST /cadastrar", PostCadastros)
  mux.HandleFunc("POST /send-verification-email", VerifyEmail)
}
