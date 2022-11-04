import { createContext } from "react";
import { TiArrowUpThick } from 'react-icons/ti';
import { sendMessages } from "../../pages/chat/Client";
import "../styles/ChatInputArea.css";

const ChatInputArea = () => {

	return (
		<div className="chat-input-area">
			<form onSubmit={sendMessages}>
				<div className="form">
					<div className="form-controller">
						<input type="text" placeholder="Type something..." name="messageText"/>
						<TiArrowUpThick fontSize="26px" style={{ marginTop: "5px", color: "blue" }} onClick={sendMessages } />
					</div>
				</div>
			</form>
		</div>
	)
}

export default ChatInputArea;