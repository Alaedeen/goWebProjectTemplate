package handlers

import (
	"github.com/Alaedeen/goWebProjectTemplate/helpers"
	"encoding/json"
	"net/http"
	"strconv"
	models "github.com/Alaedeen/goWebProjectTemplate/models"
	"github.com/Alaedeen/goWebProjectTemplate/repository"
	"crypto/sha1"
	"gopkg.in/gomail.v2"
	"github.com/sethvargo/go-password/password"
)

// UserHandler ...
type UserHandler struct {
	Repo repository.UserRepository
}

func responseFormatter (code int, status string, data interface{}, response *models.Response) {
	response.Code = code
	response.Status = status
	response.Data=data
}

// GetUsers ...
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	var responseWithCount models.ResponseWithCount
	role := r.URL.Query()["role"][0]
	offset,err0 := strconv.Atoi(r.URL.Query()["offset"][0])
	if err0 != nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err0.Error(),&response)
		responseWithCount.Response=response
		responseWithCount.Count=0
		json.NewEncoder(w).Encode(responseWithCount)
		return
	}
	limit , err:= strconv.Atoi(r.URL.Query()["limit"][0])
	if err != nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err.Error(),&response)
		responseWithCount.Response=response
		responseWithCount.Count=0
		json.NewEncoder(w).Encode(responseWithCount)
		return
	}
	result,err1,count := h.Repo.GetUsers(role,offset,limit) 
	if err1 !=nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err1.Error(),&response)
		responseWithCount.Response=response
		responseWithCount.Count=0
		json.NewEncoder(w).Encode(responseWithCount)
		return
	}
	var users []models.UserResponse
	var user models.UserResponse
	for _,res := range result {
		user.Roles= user.Roles[:0]
		helpers.UserResponseFormatter(res,&user)
		users= append(users,user)
	} 
	responseFormatter(200,"OK",users,&response)
	responseWithCount.Response=response
	responseWithCount.Count=count
	json.NewEncoder(w).Encode(responseWithCount)
}

// GetUsersByName ...
func (h *UserHandler) GetUsersByName(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	name:= r.URL.Query()["name"][0]
	role:= r.URL.Query()["role"][0]
	var response models.Response
	var responseWithCount models.ResponseWithCount
	offset,err0 := strconv.Atoi(r.URL.Query()["offset"][0])
	if err0 != nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err0.Error(),&response)
		responseWithCount.Response=response
		responseWithCount.Count=0
		json.NewEncoder(w).Encode(responseWithCount)
		return
	}
	limit , err:= strconv.Atoi(r.URL.Query()["limit"][0])
	if err != nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err.Error(),&response)
		responseWithCount.Response=response
		responseWithCount.Count=0
		json.NewEncoder(w).Encode(responseWithCount)
		return
	}
	result,err,count:= h.Repo.GetUsersByName(name,role,offset,limit)
	if err != nil {
		responseFormatter(404,"NOT FOUND",err.Error(),&response)
		responseWithCount.Response=response
		responseWithCount.Count=0
		json.NewEncoder(w).Encode(responseWithCount)
		return
	}
	var users []models.UserResponse
	var user models.UserResponse
	for _,res := range result {
		user.Roles= user.Roles[:0]
		helpers.UserResponseFormatter(res,&user)
		users= append(users,user)
	} 
	responseFormatter(200,"OK",users,&response)
	responseWithCount.Response=response
	responseWithCount.Count=count
	json.NewEncoder(w).Encode(responseWithCount)
}


// GetUser ...
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()["id"]
	var response models.Response
	id, err := strconv.Atoi(params[0])
	if err!= nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}
	result,err1 := h.Repo.GetUser(uint(id))
	if err1!=nil {
		responseFormatter(404,"NOT FOUND",err1.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}
	var user models.UserResponse
	helpers.UserResponseFormatter(result,&user)
	responseFormatter(200,"OK",user,&response)
	json.NewEncoder(w).Encode(response)
}

// GetUserBy ...
func (h *UserHandler) GetUserBy(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params:= r.URL.Query()
	var keys []string
	var values []interface{}
	var response models.Response
	for key,value := range params {
		keys = append(keys,key)
		val , err := strconv.Atoi(value[0])
		if err != nil {
			values = append(values, value[0])
		}else{
			values = append(values, uint(val))
		}
	}
	result,err:= h.Repo.GetUserBy(keys,values)
	if err != nil {
		responseFormatter(404,"NOT FOUND",err.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}
	var user models.UserResponse
	helpers.UserResponseFormatter(result,&user)
	responseFormatter(200,"OK",user,&response)
	json.NewEncoder(w).Encode(response)
}

// Login ...
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params:= r.URL.Query()
	var keys []string
	var values []interface{}
	var responseWithToken  models.ResponseWithToken
	var response models.Response
	for key,value := range params {
		keys = append(keys,key)
		val , err := strconv.Atoi(value[0])
		if err != nil {
			values = append(values, value[0])
		}else{
			values = append(values, uint(val))
		}
	}
	result,err:= h.Repo.GetUserBy(keys,values)
	if err != nil {
		responseFormatter(404,"NOT FOUND",err.Error(),&response)
		responseWithToken.Response=response
		responseWithToken.Token=""
	
		json.NewEncoder(w).Encode(responseWithToken)
		return
	}
	var user models.UserResponse
	helpers.UserResponseFormatter(result,&user)
	var role string
	if len(user.Roles)==1 {
		role = "user"
	}else if len(user.Roles)==2 {
		role = "admin"
	}else{
		role = "super admin"
	}
	token,err:= helpers.GenerateJWT(result.Name,role)
	responseFormatter(200,"OK",user,&response)
	responseWithToken.Response=response
	responseWithToken.Token=token
	json.NewEncoder(w).Encode(responseWithToken)
}

// CreateUser ...
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	var User models.User
	var UserRequest models.UserRequest
	var response models.Response
	var responseWithToken  models.ResponseWithToken
	err:=json.NewDecoder(r.Body).Decode(&UserRequest)
	if err != nil {
		responseFormatter(400,"BAD REQUEST",err.Error(),&response)
		responseWithToken.Response=response
		responseWithToken.Token=""
	
		json.NewEncoder(w).Encode(responseWithToken)
		return
	}
	helpers.UserRequestFormatter(UserRequest,&User)
	result,err1 := h.Repo.CreateUser(User)
	if err1 != nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err1.Error(),&response)
		responseWithToken.Response=response
		responseWithToken.Token=""
	
		json.NewEncoder(w).Encode(responseWithToken)
		return
	}
	var user models.UserResponse
	helpers.UserResponseFormatter(result,&user)
	token,err:= helpers.GenerateJWT(result.Name , "user")
	responseFormatter(201,"CREATED",user,&response)
	responseWithToken.Response=response
	responseWithToken.Token=token
	
	
	json.NewEncoder(w).Encode(responseWithToken)
}

// DeleteUser ...
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	params := r.URL.Query()["id"]
	var response models.Response
	id, err := strconv.Atoi(params[0])

	if err != nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}
	err1 := h.Repo.DeleteUser(uint(id))
	if err1!=nil {
		responseFormatter(404,"NOT FOUND",err1.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}
	responseFormatter(200,"OK","USER DELETED",&response)
	json.NewEncoder(w).Encode(response)
}

// UpdateUser ...
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")

	params := r.URL.Query()["id"]
	var response models.Response
	id, err1 := strconv.Atoi(params[0])
	if err1 != nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err1.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}
	

    var m map[string]interface{}
	m = make(map[string]interface{})
	var password string
	r.ParseMultipartForm(10 << 20)

	for key,value := range r.Form {
		if key=="password" {
			crypt := sha1.New()
			password= value[0]
			crypt.Write([]byte(password))
			m[key]=crypt.Sum(nil)
		}else {	
			if key!="id" {
				if value[0] == "true" {
					m[key]= true
				}else if value[0] == "false" {
					m[key]= false
				}else{
					val, err1 := strconv.Atoi(value[0])
					if err1 != nil {
						m[key]=value[0]
					}else {
						m[key]=val
					}
				}
			}
		}
	}
	err2 := h.Repo.UpdateUser(m,uint(id))
	if err2 !=nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err2.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}
	
	responseFormatter(200,"OK","USER UPDATED",&response)
	json.NewEncoder(w).Encode(response)
}

// ResetPassword ...
func (h *UserHandler) ResetPassword(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Content-Type", "application/json")
	var response models.Response
	
	email := r.URL.Query()["email"][0]
	var keys []string
	keys = append(keys,"email")
	var values []interface{}
	values = append(values,email)
	user,err:= h.Repo.GetUserBy(keys,values)
	if err != nil {
		responseFormatter(404,"NOT FOUND",err.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}
	res, err := password.Generate(10, 2, 2, false, false)
	if err != nil {
		responseFormatter(200,"OK",err,&response)
		json.NewEncoder(w).Encode(response)
	}
	var m1 map[string]interface{}
	m1 = make(map[string]interface{})
	crypt := sha1.New()
	crypt.Write([]byte(res))
	m1["password"]=crypt.Sum(nil)
	err2 := h.Repo.UpdateUser(m1,user.ID)
	if err2 !=nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err2.Error(),&response)
		json.NewEncoder(w).Encode(response)
		return
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "alaedeen.stark7@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Password Reset!")
	m.SetBody("text/html", `<div style="text-align: center">
								<h1 >Hello</h1>
								<h3>`+user.Name+`</h3>
								<p>your new password is : <b style="color : red">`+res+`</b> </p>
								<h4>Thank you for using our app!</h4>
								<p>Yours sincerely.</p>
								<p>360 video editor team.</p>
							</div>`)

	d := gomail.NewDialer("smtp.gmail.com", 587, "alaedeen.stark7@gmail.com", "Ala1995Stark")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		responseFormatter(500,"INTERNAL SERVER ERROR",err,&response)
		json.NewEncoder(w).Encode(response)
	}
	responseFormatter(200,"OK","PASSWORD UPDATED",&response)
	json.NewEncoder(w).Encode(response)
}