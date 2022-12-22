package server

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"os"
// 	"se/server/backend"
// 	"strconv"
// 	"strings"

// 	"github.com/gorilla/mux"
// )

// func defaultFunc(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	fmt.Fprintln(w, API)
// }

// func fetchBid(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	uid, valid := sessionValid(w, r)
// 	if !valid {
// 		return
// 	}

// 	path := mux.Vars(r)
// 	args := r.URL.Query()

// 	switch path["key"] {
// 	case "help":
// 		fmt.Fprint(w, BidAPI)

// 	case "get": //For single bid product
// 		pdid := args.Get("pdid")

// 		res, err := backend.GetProductBidInfo(pdid)
// 		response(w, res, err)

// 	case "set":
// 		pdid := args.Get("pdid")
// 		money := args.Get("money")

// 		res, err := backend.SetBidForBuyer(uid, pdid, money)
// 		response(w, res, err)

// 	case "delete":
// 		pdid := args.Get("pdid")

// 		res, err := backend.DeleteBid(pdid)
// 		response(w, res, err)

// 	default:
// 		http.NotFound(w, r)
// 	}
// }

// func fetchCart(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	uid, valid := sessionValid(w, r)
// 	if !valid {
// 		return
// 	}

// 	path := mux.Vars(r)
// 	args := r.URL.Query()

// 	switch path["key"] {
// 	case "help":
// 		fmt.Fprint(w, CartAPI)

// 	case "add": //For single product
// 		pdid := args.Get("pdid")
// 		amount := args.Get("amount")

// 		res, err := backend.AddProductToCart(uid, pdid, amount)
// 		response(w, res, err)

// 	case "remo":
// 		pdid := args.Get("pdid")

// 		res, err := backend.RemoveProduct(uid, pdid)
// 		response(w, res, err)

// 	case "modf":
// 		pdid := args.Get("pdid")
// 		amount := args.Get("amount")

// 		res, err := backend.ModifyProductAmount(uid, pdid, amount)
// 		response(w, res, err)

// 	case "geps": //拿商品
// 		fmt.Fprint(w, backend.GetProducts(uid))

// 	default:
// 		http.NotFound(w, r)
// 	}
// }

// func fetchHistory(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	uid, valid := sessionValid(w, r)
// 	if !valid {
// 		return
// 	}

// 	path := mux.Vars(r)
// 	args := r.URL.Query()

// 	switch path["key"] {
// 	case "help":
// 		fmt.Fprint(w, HistoryAPI)

// 	case "get":
// 		amount := args.Get("amount")
// 		newest := args.Get("newest")

// 		res, err := backend.GetHistory(uid, amount, newest)
// 		response(w, res, err)

// 	case "delete":
// 		pdid := args.Get("pdid")

// 		res, err := backend.DeleteHistory(uid, pdid)
// 		response(w, res, err)

// 	case "deleteall":
// 		fmt.Fprint(w, backend.DeleteAllHistory(uid))

// 	case "deletespec":
// 		pdids := args.Get("pdids")
// 		res, err := backend.DeleteSpecificHistory(uid, pdids)
// 		response(w, res, err)

// 	case "add":
// 		pdid := args.Get("pdid")

// 		res, err := backend.AddHistory(uid, pdid)
// 		response(w, res, err)

// 	default:
// 		http.NotFound(w, r)
// 	}
// }

// func fetchMessage(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	uid, valid := sessionValid(w, r)
// 	if !valid {
// 		return
// 	}

// 	path := mux.Vars(r)
// 	args := r.URL.Query()

// 	switch path["key"] {
// 	case "help":
// 		fmt.Fprint(w, MessageAPI)

// 	case "all":
// 		fmt.Fprint(w, backend.GetAllMessages())

// 	case "send":
// 		ruid := args.Get("remoteUID")
// 		txt := args.Get("text")

// 		res, err := backend.AddMessage(uid, ruid, txt)
// 		response(w, res, err)

// 	case "get":
// 		ruid := args.Get("remoteUID")
// 		asc := args.Get("ascend")

// 		res, err := backend.GetMessages(uid, ruid, asc)
// 		response(w, res, err)

// 	default:
// 		http.NotFound(w, r)
// 	}
// }

// func fetchOrder(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	uid, valid := sessionValid(w, r)
// 	if !valid {
// 		return
// 	}

// 	path := mux.Vars(r)
// 	args := r.URL.Query()

// 	switch path["key"] {
// 	case "help":
// 		fmt.Fprint(w, OrderAPI)

// 	case "get":
// 		fmt.Fprint(w, backend.GetOrders(uid))

// 	case "add":
// 		pdid := args.Get("pdid")
// 		amount := args.Get("amount")

// 		res, err := backend.AddOrder(uid, pdid, amount)
// 		response(w, res, err)

// 	case "del":
// 		pdid := args.Get("pdid")

// 		res, err := backend.DeleteOrder(uid, pdid)
// 		response(w, res, err)

// 	default:
// 		http.NotFound(w, r)
// 	}

// }

// func fetchProduct(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	uid, valid := sessionValid(w, r)
// 	if !valid {
// 		return
// 	}

// 	if r.Method == "POST" {
// 		if mux.Vars(r)["key"] == "postadd" {
// 			r.ParseMultipartForm(32 << 20)

// 			name := r.FormValue("name")
// 			price := r.FormValue("price")
// 			des := r.FormValue("description")
// 			amount := r.FormValue("amount")
// 			bid := r.FormValue("bid")
// 			date := r.FormValue("date")

// 			pdid, err := backend.AddProduct(uid, name, price, des, amount, bid, date)
// 			if err == nil {
// 				fmt.Fprint(w, "ok")
// 			} else {
// 				http.Error(w, "failed", http.StatusBadRequest)
// 			}

// 			file, handler, err := r.FormFile("uploadfile")
// 			if err != nil {
// 				log.Println("at upload file:", err)
// 			}
// 			defer file.Close()

// 			spt := strings.Split(handler.Filename, ".")
// 			subName := spt[len(spt)-1]

// 			fmt.Fprint(w, handler.Header)
// 			f, err := os.Create("webpage/img/" + strconv.Itoa(pdid) + "." + subName)
// 			if err != nil {
// 				log.Println("at creating file:", err)
// 			}
// 			defer f.Close()

// 			io.Copy(f, file)
// 		}
// 		return
// 	}

// 	path := mux.Vars(r)
// 	args := r.URL.Query()

// 	switch path["key"] {
// 	case "help":
// 		fmt.Fprint(w, ProductAPI)

// 	case "get":
// 		pdid := args.Get("pdid")

// 		res, err := backend.GetProductInfo(pdid)
// 		response(w, res, err)

// 	case "delete":
// 		val := args.Get("pdid")

// 		res := backend.DeleteProduct(uid, val)
// 		fmt.Fprint(w, res)

// 	case "newest":
// 		amount := args.Get("amount")

// 		res, err := backend.GetNewestProduct(amount)
// 		response(w, res, err)

// 	case "search":
// 		name := args.Get("name")
// 		min := args.Get("minprice")
// 		max := args.Get("maxprice")
// 		eval := args.Get("eval")

// 		res, err := backend.SearchProducts(name, min, max, eval)
// 		response(w, res, err)

// 	case "urproduct":
// 		fmt.Fprint(w, backend.GetSellerProduct(uid))

// 	default:
// 		http.NotFound(w, r)
// 	}
// }

// func fetchUser(w http.ResponseWriter, r *http.Request) {
// 	if !validation(w, r) {
// 		return
// 	}

// 	// POST
// 	path := mux.Vars(r)
// 	r.ParseForm()

// 	switch path["key"] {
// 	case "help":
// 		fmt.Fprint(w, UserAPI)

// 	case "login":
// 		uid, valid := backend.Login(r.FormValue("account"), r.FormValue("password"))

// 		if valid {
// 			// set session to maintain login condition
// 			login(w, r, uid)
// 			fmt.Fprintln(w, "登入成功!")
// 		} else {
// 			fmt.Fprint(w, "登入失敗")
// 		}

// 	case "regist":
// 		fmt.Fprint(w, backend.Regist(r.FormValue("account"), r.FormValue("password"), r.FormValue("name")))

// 	default:
// 		uid, valid := sessionValid(w, r)
// 		if valid {
// 			switch path["key"] {
// 			case "delete":
// 				fmt.Fprint(w, backend.DeleteUser(uid, r.FormValue("password")))

// 			case "changePassword":
// 				fmt.Fprint(w, backend.ChangePassword(uid, r.FormValue("oldPassword"), r.FormValue("newPassword")))

// 			case "changeName":
// 				fmt.Fprint(w, backend.ChangeName(uid, r.FormValue("newName")))

// 			case "logout":
// 				logout(w, r)
// 				fmt.Fprintln(w, "已登出")

// 			default:
// 				http.NotFound(w, r)
// 			}
// 		}
// 	}
// }

// func response(w http.ResponseWriter, str string, err error) {
// 	if err == nil {
// 		fmt.Fprint(w, str)
// 	} else {
// 		log.Println(err)
// 		http.Error(w, str, http.StatusBadRequest)
// 	}
// }
