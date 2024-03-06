package main

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Good struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Photo       string `json:"photo"`
	Articul     string `json:"articul"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Sizes       string `json:"sizes"`
	CategoryId  int
	Isnew       int
	Count_likes int
}

type CategoryCount struct {
	Count int `json:"count"`
}

type UserReg struct {
	Login string
	Email string
	Pass  string
}

type UserLog struct {
	Login string
	Pass  string
}

type UserTrue struct {
	Status bool   `json:"status"`
	Name   string `json:"name"`
}

type Token struct {
	Token string
}

type ResponseStat struct {
	Status bool   `json:"status"`
	Token  string `json:"token"`
}

type Order struct {
	Token  string
	Adress string
	Goods  string
}

type OrdersGet struct {
	Id     int
	Goods  string
	Adress string
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/app/get", handleGet)
	mux.HandleFunc("/app/getCount", handleGetCount)
	mux.HandleFunc("/app/registration", handleRegistration)
	mux.HandleFunc("/app/signin", handleSignIn)
	mux.HandleFunc("/app/check_tn", handleCheckToken)
	mux.HandleFunc("/app/order", handleOrder)
	mux.HandleFunc("/app/get_orders", handleGetOrders)
	mux.HandleFunc("/", handleWelcome)
	mux.HandleFunc("/signin", handleSign)
	mux.HandleFunc("/catalog", handleCatalog)
	mux.HandleFunc("/register", handleRegister)
	mux.HandleFunc("/cart", handlerCart)
	mux.HandleFunc("/card", handlerCard)
	mux.HandleFunc("/user", handlerUser)
	mux.HandleFunc("/about", handlerAbout)
	mux.HandleFunc("/delivery", handlerDelivery)
	mux.HandleFunc("/users/create", handleCreateUser)
	mux.HandleFunc("/users/read", handleReadUser)
	mux.HandleFunc("/users/update", handleUpdateUser)
	mux.HandleFunc("/users/delete", handleDeleteUser)
	mux.HandleFunc("/goods/create", handleCreateGood)
	mux.HandleFunc("/goods/read", handleGetGoods)
	mux.HandleFunc("goods/update", handleUpdateGood)
	mux.HandleFunc("/goods/delete", handleDeleteGood)
	mux.HandleFunc("/orders/read", handleReadOrders)
	mux.HandleFunc("orders/create", handleCreateOrder)
	mux.HandleFunc("orders/update", handleUpdateOrder)
	mux.HandleFunc("orders/delete", handleDeleteOrder)
	assetsDir := "C:/Users/diasi/OneDrive/Desktop/clothing-store-main/frontend"
	mux.Handle("/frontend/", http.StripPrefix("/frontend/", http.FileServer(http.Dir(assetsDir))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(filepath.Join(assetsDir, "css")))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(filepath.Join(assetsDir, "images")))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir(filepath.Join(assetsDir, "js")))))

	log.Println("Server is running on :8000...")
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
func renderHTML(w http.ResponseWriter, filename string) {
	absPath := filepath.Join("C:\\Users\\diasi\\OneDrive\\Desktop\\clothing-store-main\\frontend", filename)
	tmpl, err := template.ParseFiles(absPath)
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, filename)
}
func handleWelcome(w http.ResponseWriter, r *http.Request) {
	// Assuming your HTML file is in the "templates" directory
	renderHTML(w, "welcome-page.html")
}

func handleCatalog(w http.ResponseWriter, r *http.Request) {
	// Assuming your HTML file is in the "templates" directory
	renderHTML(w, "catalog.html")
}
func handleSign(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "signin.html")
}
func handleRegister(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "registration.html")
}
func handlerCart(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "cart.html")
}
func handlerCard(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "card.html")
}
func handlerUser(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "user_acc.html")
}
func handlerAbout(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "about.html")
}
func handlerDelivery(w http.ResponseWriter, r *http.Request) {
	renderHTML(w, "delivery.html")
}
func handleGet(w http.ResponseWriter, r *http.Request) {
	// Allow cross-domain requests
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Create a connection to the database
	db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	if err != nil {
		log.Println("Error connecting to the database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Get the query parameters
	categoryFilter := ""
	categoryId := r.URL.Query().Get("category_id")
	if categoryId != "" {
		categoryFilter = " AND `category_id` = " + categoryId
	}

	newFilter := ""
	isNew := r.URL.Query().Get("is_new")
	if isNew == "1" {
		newFilter = " AND `is_new` = 1"
	}

	bestFilter := ""
	isBest := r.URL.Query().Get("is_best")
	if isBest == "1" {
		bestFilter = " AND `count_likes` > 10"
	}

	idFilter := ""
	id := r.URL.Query().Get("id")
	if id != "" {
		idFilter = " AND `id` IN(" + id + ") "
	}

	// Send the query to the database
	query := "SELECT * FROM inordic.goods WHERE 1" + categoryFilter + idFilter + newFilter + bestFilter
	log.Println("Database Query:", query) // Add this line for logging
	result, err := db.Query(query)
	if err != nil {
		log.Println("Error executing the database query:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer result.Close()

	// Get and output the result
	goods := []Good{}
	for result.Next() {
		good := Good{}
		result.Scan(&good.Id, &good.Title, &good.Photo, &good.Articul, &good.Price, &good.Description, &good.Sizes, &good.CategoryId, &good.Isnew, &good.Count_likes)
		goods = append(goods, good)
	}

	// Log the retrieved data
	log.Println("Retrieved Data:", goods)

	// Encode to JSON and send the response
	jsonData, err := json.Marshal(goods)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handleGetCount(w http.ResponseWriter, r *http.Request) {
	// Allow cross-origin requests
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create a connection to the database
	db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer db.Close()

	// Create an array to store the count
	count := []CategoryCount{}

	// Define a map for category count
	list := map[int]string{
		1: "category_id",
		2: "category_id",
		3: "category_id",
		4: "is_new",
		5: "count_likes",
	}

	// Iterate to get the count for each category
	for i := 0; i < len(list); i++ {
		y := strconv.Itoa(i + 1)
		sign := " = "
		if list[i+1] == "is_new" {
			y = strconv.Itoa(1)
		}
		if list[i+1] == "count_likes" {
			sign = " > "
			y = strconv.Itoa(10)
		}
		rows, err := db.Query("SELECT COUNT(*) FROM `goods` WHERE `" + list[i+1] + "`" + sign + y)
		if err != nil {
			fmt.Println("Error executing the database query:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Fetch and store the count
		for rows.Next() {
			countSre := CategoryCount{}
			err := rows.Scan(&countSre.Count)
			if err != nil {
				fmt.Println("Error scanning row:", err)
				continue
			}
			count = append(count, countSre)
		}
	}

	// Convert array to JSON
	countJson, err := json.Marshal(count)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Send the JSON response
	fmt.Fprintf(w, string(countJson))
}

func handleRegistration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//подключаемся к нашей базе данных
	db, _ := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	//считать данные из тела запроса
	body, _ := ioutil.ReadAll(r.Body)
	//перелисть их в структуру
	regData := UserReg{}
	json.Unmarshal(body, &regData)
	//проверяем пользователя на оригинальность
	countCommon := 0
	resultCount, _ := db.Query("SELECT COUNT(*) FROM `users` WHERE `login` = ? OR `email` = ?", regData.Login, regData.Email)
	for resultCount.Next() {
		resultCount.Scan(&countCommon)
	}
	respStat := ResponseStat{}
	if countCommon == 0 {
		userToken := md5.Sum([]byte(regData.Login))
		userPass := md5.Sum([]byte(regData.Pass))
		userPassLine := fmt.Sprintf("%x", userPass)
		userTokenLine := fmt.Sprintf("%x", userToken)
		_, err := db.Query("INSERT INTO `users`(`login`, `email`, `password`, `token`) VALUES(?,?,?,?)", regData.Login, regData.Email, userPassLine, userTokenLine)
		if err != nil {
			fmt.Println(err)
		}
		//вернуть результат
		respStat = ResponseStat{Status: true, Token: userTokenLine}
	} else {
		respStat = ResponseStat{Status: false, Token: "0"}
	}
	jsonrespStat, _ := json.Marshal(respStat)
	fmt.Fprintf(w, string(jsonrespStat))

}

func handleSignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//подключаемся к нашей базе данных
	db, _ := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	//считать данные из тела запроса
	body, _ := ioutil.ReadAll(r.Body)
	//перелисть их в структуру
	logData := UserLog{}
	json.Unmarshal(body, &logData)
	//проверяем пользователя на оригинальность
	countCommon := 0
	userPass := md5.Sum([]byte(logData.Pass))
	userPassLine := fmt.Sprintf("%x", userPass)
	resultCount, _ := db.Query("SELECT COUNT(*) FROM `users` WHERE (`login` = ? OR `email` = ?) AND `password` = ?", logData.Login, logData.Login, userPassLine)
	for resultCount.Next() {
		resultCount.Scan(&countCommon)
	}
	respStat := ResponseStat{}
	if countCommon == 1 {
		token := ""
		tokenData, _ := db.Query("SELECT token FROM `users` WHERE `login` = ? OR `email` = ?", logData.Login, logData.Login)
		for tokenData.Next() {
			tokenData.Scan(&token)
		}
		respStat = ResponseStat{Status: true, Token: token}
		//вернуть результат
	} else {
		respStat = ResponseStat{Status: false, Token: "0"}
	}
	jsonrespStat, _ := json.Marshal(respStat)
	fmt.Fprintf(w, string(jsonrespStat))
}

func handleCheckToken(w http.ResponseWriter, r *http.Request) {
	//разрешаем посещать наш сайт всем
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//подключаемся к нашей базе данных
	db, _ := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	//считать данные из тела запроса
	body, _ := ioutil.ReadAll(r.Body)
	//перелисть их в структуру
	token := Token{}
	json.Unmarshal(body, &token)
	//проверяем пользователя на оригинальность
	result, _ := db.Query("SELECT COUNT(*) FROM `users` WHERE `token`= ? ", token.Token)
	resultData := 0
	for result.Next() {
		result.Scan(&resultData)
	}
	response := UserTrue{}
	if resultData != 0 {
		resultName, _ := db.Query("SELECT `login` FROM `users` WHERE `token` = ?", token.Token)
		name := ""
		for resultName.Next() {
			resultName.Scan(&name)
		}
		response = UserTrue{Status: true, Name: name}
		responseJson, _ := json.Marshal(response)
		fmt.Fprintf(w, string(responseJson))
	} else {
		response = UserTrue{Status: false, Name: ""}
		responseJson, _ := json.Marshal(response)
		fmt.Fprintf(w, string(responseJson))
	}
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, _ := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	//считать данные из тела запроса
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(body)
	order := Order{}
	json.Unmarshal(body, &order)
	user_id_data := 0
	user_id, _ := db.Query("SELECT `id` FROM `users` WHERE `token`=?", order.Token)
	fmt.Println(order)
	for user_id.Next() {
		user_id.Scan(&user_id_data)
	}
	db.Exec("INSERT INTO `orders`(`goods`, `adress`, `user_id`) VALUES(?, ?, ?)", order.Goods, order.Adress, user_id_data)
}

func handleGetOrders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db, _ := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	//считать данные из тела запроса
	body, _ := ioutil.ReadAll(r.Body)
	token := ""
	json.Unmarshal(body, &token)
	result, _ := db.Query("SELECT `order_id`, `goods`, `adress` FROM `users` JOIN `orders` ON users.id=orders.user_id WHERE users.token=?", token)
	orders := []OrdersGet{}
	for result.Next() {
		order := OrdersGet{}
		result.Scan(&order.Id, &order.Goods, &order.Adress)
		orders = append(orders, order)
	}
	jsonData, _ := json.Marshal(orders)
	fmt.Fprintf(w, string(jsonData))
}

// CRUD FOR USERS
func generateToken() string {
	token := fmt.Sprintf("%d", time.Now().UnixNano())
	hashedToken := md5.Sum([]byte(token))
	return hex.EncodeToString(hashedToken[:])
}
func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Распаковка JSON-запроса в структуру UserReg
		db, _ := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
		var newUser UserReg
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		defaultToken := generateToken()

		// Хеширование пароля в формате md5
		hashedPassword := fmt.Sprintf("%x", md5.Sum([]byte(newUser.Pass)))

		// Вставка нового пользователя в базу данных
		_, err = db.Exec("INSERT INTO users (login, email, password, token) VALUES (?, ?, ?, ?)",
			newUser.Login, newUser.Email, hashedPassword, defaultToken)
		if err != nil {
			log.Println("Error creating user:", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		// Handle other HTTP methods or return an error
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// ReadUser (GET)
func handleReadUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query all users from the database
	db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		log.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT password, login, email FROM users")
	if err != nil {
		http.Error(w, "Error readin222g users", http.StatusInternalServerError)
		log.Println("Error reading users:", err)
		return
	}
	defer rows.Close()

	// Build a slice of users
	var users []UserReg
	for rows.Next() {
		var user UserReg
		err := rows.Scan(&user.Pass, &user.Login, &user.Email)
		if err != nil {
			http.Error(w, "Error scanning users", http.StatusInternalServerError)
			log.Println("Error scanning users:", err)
			return
		}
		users = append(users, user)
	}

	// Encode users to JSON and send the response
	jsonData, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		log.Println("Error encoding JSON:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// UpdateUser (PUT)
func handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	// Распаковка JSON-запроса в структуру UserReg
	db, _ := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	var updatedUser UserReg
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Проверка существования пользователя
	var existingUser UserReg
	err = db.QueryRow("SELECT * FROM users WHERE login=?", updatedUser.Login).
		Scan(&existingUser.Login, &existingUser.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error checking user existence", http.StatusInternalServerError)
		return
	}

	// Обновление пользователя в базе данных по логину
	_, err = db.Exec("UPDATE users SET email=? WHERE login=?",
		updatedUser.Email, updatedUser.Login)
	if err != nil {
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteUser (DELETE)
func handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	// Получение логина пользователя из запроса
	db, _ := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	login := r.URL.Query().Get("login")
	if login == "" {
		http.Error(w, "Login parameter is required", http.StatusBadRequest)
		return
	}

	// Удаление пользователя из базы данных по логину
	_, err := db.Exec("DELETE FROM users WHERE login=?", login)
	if err != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// CRUD ДЛЯ GOODS
func handleGetGoods(w http.ResponseWriter, r *http.Request) {
	// Fetch all goods from the database
	db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM goods")
	if err != nil {
		http.Error(w, "Error reading goods", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build a list of goods
	var goods []Good
	for rows.Next() {
		var g Good
		err := rows.Scan(&g.Id, &g.Title, &g.Photo, &g.Articul, &g.Price, &g.Description, &g.Sizes, &g.CategoryId, &g.Isnew, &g.Count_likes)
		if err != nil {
			http.Error(w, "Error scanning goods", http.StatusInternalServerError)
			return
		}
		goods = append(goods, g)
	}

	// Encode to JSON and send the response
	jsonData, err := json.Marshal(goods)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func handleCreateGood(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body into a Good object
	var newGood Good
	err := json.NewDecoder(r.Body).Decode(&newGood)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Insert the new good into the database
	db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO goods (title, photo, articul, price, description, sizes, category_id, is_new, count_likes) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		newGood.Title, newGood.Photo, newGood.Articul, newGood.Price, newGood.Description, newGood.Sizes, newGood.CategoryId, newGood.Isnew, newGood.Count_likes)
	if err != nil {
		http.Error(w, "Error creating good", http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly created good
	newID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting the ID of the created good", http.StatusInternalServerError)
		return
	}

	// Set the ID in the response
	newGood.Id = int(newID)

	// Encode to JSON and send the response
	jsonData, err := json.Marshal(newGood)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonData)
}

func handleUpdateGood(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the request body into a Good object
	var updatedGood Good
	err := json.NewDecoder(r.Body).Decode(&updatedGood)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Update the good in the database
	db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE goods SET title=?, photo=?, articul=?, price=?, description=?, sizes=?, category_id=?, is_new=?, count_likes=? WHERE id=?",
		updatedGood.Title, updatedGood.Photo, updatedGood.Articul, updatedGood.Price, updatedGood.Description, updatedGood.Sizes, updatedGood.CategoryId, updatedGood.Isnew, updatedGood.Count_likes, updatedGood.Id)
	if err != nil {
		http.Error(w, "Error updating good", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleDeleteGood(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the ID from the request URL
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID parameter is required", http.StatusBadRequest)
		return
	}

	// Convert the ID to an integer
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameter", http.StatusBadRequest)
		return
	}

	// Delete the good from the database
	db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	if err != nil {
		http.Error(w, "Error connecting to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM goods WHERE id=?", id)
	if err != nil {
		http.Error(w, "Error deleting good", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

//CRUD ДЛЯ ORDERS

// Create (POST)
func handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Decode JSON request to OrdersGet struct
		var newOrder OrdersGet
		err := json.NewDecoder(r.Body).Decode(&newOrder)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
		// Insert new order into the database
		_, err = db.Exec("INSERT INTO orders (goods, adress) VALUES (?, ?)", newOrder.Goods, newOrder.Adress)
		if err != nil {
			http.Error(w, "Error creating order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Read (GET)
func handleReadOrders(w http.ResponseWriter, r *http.Request) {
	// Query orders from the database
	db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
	rows, err := db.Query("SELECT order_id, goods, adress FROM orders")
	if err != nil {
		http.Error(w, "Error reading orders", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Build a list of orders
	var orders []OrdersGet
	for rows.Next() {
		var o OrdersGet
		err := rows.Scan(&o.Id, &o.Goods, &o.Adress)
		if err != nil {
			http.Error(w, "Error scanning orders", http.StatusInternalServerError)
			return
		}
		orders = append(orders, o)
	}

	// Encode to JSON and send the response
	jsonData, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// Update (PUT)
func handleUpdateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		// Decode JSON request to OrdersGet struct
		var updatedOrder OrdersGet
		err := json.NewDecoder(r.Body).Decode(&updatedOrder)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Update order in the database
		db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
		_, err = db.Exec("UPDATE orders SET goods=?, adress=? WHERE order_id=?", updatedOrder.Goods, updatedOrder.Adress, updatedOrder.Id)
		if err != nil {
			http.Error(w, "Error updating order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Delete (DELETE)
func handleDeleteOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		// Extract order ID from the request
		id := r.URL.Query().Get("id")
		if id == "" {
			http.Error(w, "ID parameter is required", http.StatusBadRequest)
			return
		}

		// Delete order from the database
		db, err := sql.Open("mysql", "root:Samsungtab7!@tcp(127.0.0.1:3306)/inordic")
		_, err = db.Exec("DELETE FROM orders WHERE order_id=?", id)
		if err != nil {
			http.Error(w, "Error deleting order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
