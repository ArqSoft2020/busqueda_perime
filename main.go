package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//main
func main() {
	a := App{}
	// You need to set your Username and Password here
	// user - password- name container - name db
	a.Initialize("root", "password", "perime-search-db", "perime-search-db")

	a.Run(":1859")

}

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(user, password, dockername string, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, dockername, dbname)

	var err error
	a.DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	//Categoria
	a.Router.HandleFunc("/categorys", a.getCategorys).Methods("GET")
	a.Router.HandleFunc("/category", a.createCategory).Methods("POST")
	a.Router.HandleFunc("/category/{id:[0-9]+}", a.getCategory).Methods("GET")
	a.Router.HandleFunc("/category/{id:[0-9]+}", a.updateCategory).Methods("PUT")
	a.Router.HandleFunc("/category/{id:[0-9]+}", a.deleteCategory).Methods("DELETE")
	//Producto
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}

///bloquejson
//manejo de error
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

//bloquejson

func (a *App) getCategorys(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getCategorys(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createCategory(w http.ResponseWriter, r *http.Request) {
	var u category
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.createCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) getCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	u := category{ID: id}
	if err := u.getCategory(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Categoria not found press f")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) updateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var u category
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Categoria ID")
		return
	}

	u := category{ID: id}
	if err := u.deleteCategory(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//model category
type category struct {
	ID            int    `json:"id_category"`
	Name_Category string `json:"name_category"`
	Type_Category string `json:"type_category "`
}

func (u *category) getCategory(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT Name_Category, Type_Category FROM categorys WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.Name_Category, &u.Type_Category)
}

func (u *category) updateCategory(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE categorys SET Name_Category='%s', Type_Category='%s' WHERE id=%d", u.Name_Category, u.Type_Category, u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *category) deleteCategory(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM categorys WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *category) createCategory(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO categorys(Name_Category, Type_Category) VALUES('%s', '%s')", u.Name_Category, u.Type_Category)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getCategorys(db *sql.DB, start, count int) ([]category, error) {
	statement := fmt.Sprintf("SELECT id, Name_Category, Type_Category FROM categorys LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	categorys := []category{}

	for rows.Next() {
		var u category
		if err := rows.Scan(&u.ID, &u.Name_Category, &u.Type_Category); err != nil {
			return nil, err
		}
		categorys = append(categorys, u)
	}

	return categorys, nil
}

//product
func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var u product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := u.createProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, u)
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	u := product{ID: id}
	if err := u.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Producto not found press f")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var u product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&u); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resquest payload")
		return
	}
	defer r.Body.Close()
	u.ID = id

	if err := u.updateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, u)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Producto ID")
		return
	}

	u := product{ID: id}
	if err := u.deleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

//model product
type product struct {
	ID                  int    `json:"id_product"`
	Id_Category         int    `json:"id_category"`
	Name_Product        string `json:"name_product"`
	Description_Product string `json:"description_product"`
}

func (u *product) getProduct(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT Id_Category, Name_Product, Description_Product FROM products WHERE id=%d", u.ID)
	return db.QueryRow(statement).Scan(&u.Id_Category, &u.Name_Product, &u.Description_Product)
}

func (u *product) updateProduct(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE products SET Id_Category=%d, Name_Product='%s', Description_Product='%s' WHERE id=%d", u.Id_Category, u.Name_Product, u.Description_Product, u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *product) deleteProduct(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM products WHERE id=%d", u.ID)
	_, err := db.Exec(statement)
	return err
}

func (u *product) createProduct(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO products(Id_Category, Name_Product, Description_Product) VALUES(%d, '%s', '%s')", u.Id_Category, u.Name_Product, u.Description_Product)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&u.ID)

	if err != nil {
		return err
	}

	return nil
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	statement := fmt.Sprintf("SELECT id, Id_Category, Name_Product, Description_Product FROM products LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var u product
		if err := rows.Scan(&u.ID, &u.Id_Category, &u.Name_Product, &u.Description_Product); err != nil {
			return nil, err
		}
		products = append(products, u)
	}

	return products, nil
}

