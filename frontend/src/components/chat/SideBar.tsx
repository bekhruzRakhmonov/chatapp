import * as React from "react";
import SearchComponent from "./SearchComponent";
import { handleClick } from "../../pages/chat/Chat";
import customFetcher  from "../../utils/fetchInstance";
import "../styles/SideBar.css";


export let myDivRef = React.createRef<HTMLDivElement>();

const getChats = async (): Promise<any> => {
	let response = await customFetcher("/chat/get-chats", {
		method: "GET"
	})
	console.log(response)
}

const SideBar = () => {
	// const node = myDivRef.current!.focus()
	React.useEffect(() => {
		(async () => {
			await getChats();
		})();
	},[])
	return (
		<div className="side-bar">
			<SearchComponent />
			<br/>
			<div className="side-bar-child">
				<div className="bar" onClick={handleClick} id="bexruz" style={{ backgroundColor: "" }}>
					<div className="bar-chat" id="bexruz">
						<b id="bexruz">Bexruz</b>
						<br/>
						<small id="bexruz">
							You: hello
						</small>
						<hr/>
					</div>
				</div>
				<div className="bar" onClick={handleClick} id="feruz">
					<div className="bar-chat" id="feruz">
						<b id="feruz">Feruz</b>
						<br/>
						<small id="feruz">
							You: hello
						</small>
						<hr/>
					</div>
				</div>
				<div className="bar" onClick={handleClick} id="shoxruz">
					<div className="bar-chat" id="shoxruz">
						<b id="shoxruz">Shoxruz</b>
						<br/>
						<small id="shoxruz">
							You: hello
						</small>
						<hr/>
					</div>
				</div>
			</div>		
		</div>
	)
}

export default SideBar;