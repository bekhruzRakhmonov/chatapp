import { useEffect, useState, createContext, useContext } from "react";
import { usernameContext } from "../../pages/chat/Chat";
import ChatInputArea from "./ChatInputArea"
import ChatMessageArea from "./ChatMessageArea"
import "../styles/ChatArea.css";

let newUsernameContext = createContext<string>("")

const ChatArea = () => {

	let username = useContext(usernameContext)

	return (
		<div className="chat-area">
			<div className="chat-area-child">
				<div className="chat-details">
					<h2 id="chat-inbound-username">{username}</h2>
				</div>
				<hr/>
				<ChatMessageArea inboundUsername={username}/>
			</div>
			<ChatInputArea/>
		</div>
	)
}

export {newUsernameContext};
export default ChatArea;