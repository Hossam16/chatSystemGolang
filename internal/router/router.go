package router

import (
    "github.com/gorilla/mux"
    "github.com/Hossam16/go-chat-creation-api/internal/handlers"
    "github.com/Hossam16/go-chat-creation-api/configs"
)

func InitRouter() *mux.Router {
    router := mux.NewRouter().StrictSlash(true)
    
    router.HandleFunc(configs.ChatsRoute, handlers.CreateChat).Methods("POST")
    router.HandleFunc(configs.MessagesRoute, handlers.CreateMessage).Methods("POST")

    return router
}
