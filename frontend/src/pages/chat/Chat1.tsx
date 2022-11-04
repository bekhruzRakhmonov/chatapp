import { useState,useEffect } from 'react';
import { useLocation, useMatch, Navigate } from 'react-router-dom';
import Client from './Client';
import { EffectCallback } from '../../utils/Types';
import customFetcher from '../../utils/fetchInstance';
import SideBar from '../../components/chat/SideBar'
import ChatArea from '../../components/chat/ChatArea'

let uname:string;

const getLastChat = async (): Promise<any> => {
	let response = await customFetcher("/chat/get-last-chat",{
		method: "GET"
	})

	console.log("getLastChat",response.data)
	return response
}

const Chat = () => {
	const [pathname,setPathname] = useState<string>("")
	const [username,setUsername] = useState<string>("")
	const location = useLocation()

	async function run(){
		let res = await getLastChat()
		setUsername(res.data.username)
	}
	run()

	console.log(username)

	let nextURL = `#${username}`;
	let nextTitle = 'New chat';
	let nextState = { next: false };


	useEffect(() => {
		setPathname(location.pathname)
		console.log(window.history)
	},window.history.state.next)
	
	console.log("uname",uname)

	setInterval(()=>{
		window.history.pushState(nextState, nextTitle, nextURL);
		nextURL += ""
		window.history.state.next = !window.history.state.next 
	},2000)

	return (
		<div className="chat-app">
			<SideBar />
			<ChatArea />
		</div>
	)
}

export default Chat;

export const handleClick = (e:any) => {
	console.log(e.target.id)
	uname = e.target.id
}












// console.log(username)
	// let routeMatch = useMatch("/chat/:username")
	// let redirectUrl: string;
	// if (pathname !== "") {
	// 	if (routeMatch === null) {
	// 		if (typeof pathname.match("^/chat$") === "object") {
	// 			redirectUrl= `/chat/${username}`
	// 			return <Navigate to={redirectUrl}/>
	// 		} else if (typeof pathname.match("^/chat/$") === "object") {
	// 			redirectUrl= `/chat${username}`
	// 			return <Navigate to={redirectUrl}/>
	// 		}
	// 	}
	// }

	// if (routeMatch !== null && pathname !== "") {
	// 	const client = new Client({pathname:pathname})
	// }