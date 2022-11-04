import { useState, useEffect, useRef, useContext } from "react";
import { isExpired, decodeToken } from "react-jwt";
import { usernameContext } from "../../pages/chat/Chat";
import customFetcher from "../../utils/fetchInstance";
import { cookies } from "../../utils/Cookies";
import "../styles/ChatMessageArea.css";
import { getUsername } from "../../utils/Cookies";

interface Props {
	inboundUsername: string;
}

interface Message {
	ID: number,
	Outbound: string,
	Inbound: string,
	Message: string,
	Date: string,	
}

const getMessages = async (username: string): Promise<any> => {
	console.log(`/chat/get-messages/${username}`)
	let response = await customFetcher(`/chat/get-messages/${username}`,{
		method: "GET"
	})
	let {resp,data} = response

	return data
}

const ChatMessageArea = ({inboundUsername}:Props) => {
	const chatMessageArea = useRef<HTMLDivElement>(null)

	let message: Message = {ID:0,Outbound:"",Inbound:"",Message:"",Date:""};
	const [messages,setMessages] = useState<Message[]>([message])

	const scrollToBottom = () => {
	    chatMessageArea.current?.scrollTo({ top: chatMessageArea.current?.scrollHeight, behavior: "smooth" })
	}

	useEffect(()=>{
		async function getMsgs(){
			let data = await getMessages(inboundUsername)
			setMessages(data.messages)
		}
		getMsgs();
		scrollToBottom()
	},[chatMessageArea.current,inboundUsername])

	if (messages.length === 0) {
		return <h3 className="no-messages-yet">No messages yet</h3>
	}

	return (
		<div className="chat-message-area" ref={chatMessageArea}>
			<div className="chat-messages">
				{
					messages !== null ? messages.map(msg=> {
						if (msg.ID === 0) {
							return  <div className="loading" key={msg.ID}><div></div><div></div><div></div></div>
						}

						if (msg.Outbound === getUsername()){
							return  <div key={msg.ID} className="outbound-message">
										<p id="message">{msg.Message}</p>
									</div>
						} else {
							return  <div key={msg.ID} className="inbound-message">
										<p id="message">{msg.Message}</p>
									</div>
						}
					})
					: <div className="loading"><div></div><div></div><div></div></div>

				}
			</div>
		</div>
	)
}

export default ChatMessageArea;