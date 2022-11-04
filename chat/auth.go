package chat

import dbutils "example.com/chatapp/db/utils"

func IsUserRegistered(client *Client) bool {
	is_registered := GetChat(client.peer.outbound,client.peer.inbound)

	if !is_registered {
		return false
	}

	return true

}

func RegisterUser(client *Client) bool {
	_,exists := dbutils.GetUser(client.peer.inbound)

	if !exists {
		return false
	}

	CreateChat(client.peer.outbound,client.peer.inbound)
	return true
}
