package views

import (
	"net/http"
	"go_ws/models"
	"go_ws/tools"
	"encoding/json"
	"log"
)

func APIMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	room := r.FormValue("room_id")
	signed, useId := tools.SingleSign(r)
	if signed == false {
		http.Error(w, "Please sign in firstly!", http.StatusOK)
		return
	}

	//todo fixme improve it
	messageRows, err := models.SelectQuery(
		"select message.id as mid , message.content as content, " +
			"message.room_id as rid, message.user_id as from_uid, message.create_date, " +
			"user.username as from_name, guser.avatar_image_small as image " +
			"from web_chatmessage as message " +
			"inner join auth_user as user on user.id = message.user_id and user.is_active = 1 " +
			"inner join web_ggacuser as guser on guser.user_ptr_id = user.id" +
			" where message.room_id = ? order by message.create_date desc limit 30", room)
	if err != nil {
		log.Printf("error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err1 := models.UpdateQuery(
		"UPDATE web_chatmessage INNER JOIN web_chatroom on web_chatroom.id = web_chatmessage.room_id" +
			" set unread = 0 WHERE web_chatroom.id = ? AND web_chatmessage.user_id != ?", room, useId)
	if err1 != nil {
		log.Printf("error: %v", err1)
		http.Error(w, err1.Error(), http.StatusInternalServerError)
		return
	}
	reverseList, _ := tools.Reverse(messageRows)
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(tools.ApiJsonNormalization(reverseList, 0, "success"))
	return
}