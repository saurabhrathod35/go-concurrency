package main

import "net/http"

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}
func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}
func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {

}
func (app *Config) LogOutPage(w http.ResponseWriter, r *http.Request) {

}
func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}
func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {

}
func (app *Config) ActivateAccountPage(w http.ResponseWriter, r *http.Request) {

}
