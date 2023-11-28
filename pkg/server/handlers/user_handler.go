package server

// type userStore interface {
// 	GetReviewers() ([]users.User, error)
// 	GetReviewerById(id int) (users.User, error)
// }

// type UserHandler struct {
// 	store userStore
// }

// func NewUserHandler(s userStore) *UserHandler {
// 	return &UserHandler{
// 		store: s,
// 	}
// }

// func (h *UserHandler) GetReviewers(w http.ResponseWriter, r *http.Request) {
// 	reviewers, err := h.store.GetReviewers()
// 	if err != nil {
// 		slog.Error(err.Error())
// 		utils.ErrorJSON(w, err)
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, reviewers)
// }

// func (h *UserHandler) GetReviewerById(w http.ResponseWriter, r *http.Request) {
// 	userId, err := getAuthUserId(r)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		utils.ErrorJSON(w, err)
// 		return
// 	}
// 	reviewer, err := h.store.GetReviewerById(userId)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		utils.ErrorJSON(w, err)
// 		return
// 	}

// 	utils.WriteJSON(w, http.StatusOK, reviewer)
// }

// func SignUp() {
// 	// Hash the users password

// 	// Generate a salt

// 	// Hash the salt and the password together

// 	// Join the hashed result and the salt together

// 	// Create a new user and save it

// 	// return the user
// }

// func getAuthUserId(r *http.Request) (int, error) {
// 	authHeader := r.Header.Get("Authorization")
// 	splits := strings.Split(authHeader, " ")
// 	var token string
// 	if len(splits) > 1 {
// 		token = splits[1]
// 		userId, err := strconv.Atoi(token)
// 		if err != nil {
// 			return 0, errors.New("invalid token")
// 		}
// 		return userId, nil
// 	} else {
// 		return 0, errors.New("no token provided")
// 	}
// }
