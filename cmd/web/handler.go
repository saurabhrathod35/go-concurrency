package main

import (
	"final-project/data"
	"fmt"
	"html/template"
	"net/http"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}
func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}
func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {
	_ = app.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// get email and password
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	user, err := app.Modelas.User.GetByEmail(email)
	if err != nil {
		app.Session.Put(r.Context(), "error", "invalid creds")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// check password
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "invalid creds")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	if !validPassword {
		msg := Message{
			To:      email,
			Subject: "Failed to login attampt",
			Data:    "Invalid login attampt",
		}
		app.sendEmail(msg)
		app.Session.Put(r.Context(), "error", "invalid creds")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// log user in
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)
	app.Session.Put(r.Context(), "flash", "Successfull Login")
	// redirect user
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
func (app *Config) LogOutPage(w http.ResponseWriter, r *http.Request) {
	app.Session.Destroy(r.Context())
	app.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}
func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}
	// create user
	u := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
		IsAdmin:   0,
	}
	_, err = u.Insert(u)
	if err != nil {
		app.Session.Put(r.Context(), "error", "unable to create user.")
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// send email

	url := fmt.Sprintf("http://localhost:9091/activate?email=%s", u.Email)
	signUrl := GenerateTokenFromString(url)
	app.InfoLog.Println(signUrl)
	msg := Message{
		To:       u.Email,
		Subject:  "Welcome",
		Template: "conformation",
		Data:     template.HTML(signUrl),
	}
	app.sendEmail(msg)
	app.Session.Put(r.Context(), "flash", "Confirmation email Sent")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
func (app *Config) ActivateAccountPage(w http.ResponseWriter, r *http.Request) {

	// validate url
	url := r.RequestURI
	testURL := fmt.Sprintf("http://localhost:9091%s", url)
	okay := VerifyToken(testURL)
	if !okay {
		app.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// activate user
	u, err := app.Modelas.User.GetByEmail(r.URL.Query().Get("email"))
	if err != nil {
		app.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u.Active = 1
	err = u.Update()
	if err != nil {
		app.Session.Put(r.Context(), "error", "Failed to update user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	app.Session.Put(r.Context(), "flash", "Account activated")
	http.Redirect(w, r, "/login", http.StatusSeeOther)

	// generate invoice

	// send email with attachment

	// send email with pdf
}
