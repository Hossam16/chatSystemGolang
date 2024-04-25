package handlers

import (
	"context"
	"log"
	"time"
	"strconv"
	"strings"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/Hossam16/go-chat-creation-api/pkg/sidekiq"
	"github.com/Hossam16/go-chat-creation-api/pkg/redis"
	"github.com/Hossam16/go-chat-creation-api/pkg/network"
	"github.com/Hossam16/go-chat-creation-api/configs"
)

type chatResponse struct {
	Number        int64   `json:"number"`
	AccessToken   string  `json:"access_token"`
}

type appsApiResponse struct {
	Name          string  `json:"name"`
	AccessToken   string  `json:"access_token"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	ChatCount     int64   `json:"chat_count"`
}

func CreateChat(w http.ResponseWriter, r *http.Request) {
	// Read in request
	accessToken := mux.Vars(r)["access_token"]

	// Get next number
	redisClient, err := redis.GetRedis()
	if err != nil {
		network.RespondErr(w, err)
		return
	}

	redisLocker := redis.GetLocker()
	key := configs.RedisChatKeyPrefix + accessToken
	ctx := context.Background()
	// Begin critical section
	lock, err := redisLocker.Obtain(ctx,key + "_LOCK", 1000*time.Millisecond, nil)
	if err != nil {
		defer lock.Release(ctx)
		network.RespondErr(w, err)
		return
	}

	exists, err := redisClient.Exists(ctx,key).Result()
	if err != nil {
		defer lock.Release(ctx)
		network.RespondErr(w, err)
		return
	} else if exists == 0 {
		log.Println("Key " + key + " not found in Redis, requsting chat count from Rails instead")
		appsResp, err := RequestChats(accessToken)
		if err != nil {
			defer lock.Release(ctx)
			network.RespondErr(w, err)
			return
		}
		redisClient.Set(ctx,key, appsResp.ChatCount, time.Hour)
	}

	nextNum, err := redisClient.Incr(ctx,key).Result()
	defer lock.Release(ctx)
	// End critical section
	if err != nil {
		network.RespondErr(w, err)
		return
	}

	// Push to Sidekiq
	err = sidekiq.Push(configs.RedisChatQueue, configs.ChatWorkerClass, accessToken, strconv.FormatInt(nextNum, 10))
	if err != nil {
		network.RespondErr(w, err)
		return
	}

	// Send response
	resp := chatResponse{Number: nextNum, AccessToken: accessToken}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		network.RespondErr(w, err)
		return
	}

	network.Respond(w, respBytes, http.StatusCreated)
}

func RequestChats(accessToken string) (appsApiResponse, error) {
	var resp appsApiResponse
	r, err := http.Get(strings.Replace(configs.AppEndpoint + configs.ApplicationsRoute, "{access_token}", accessToken, 1))
	if err != nil {
		return resp, err
	}
	json.NewDecoder(r.Body).Decode(&resp)
	return resp, nil
}

