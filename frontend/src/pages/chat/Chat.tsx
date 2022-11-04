import { useState,useEffect,useRef,createContext } from 'react';
import { useLocation, useMatch, Navigate } from 'react-router-dom';
import Client from './Client';
import customFetcher from '../../utils/fetchInstance';
import SideBar from '../../components/chat/SideBar'
import ChatArea from '../../components/chat/ChatArea'

const getLastChat = async (): Promise<any> => {
	let response = await customFetcher("/chat/get-last-chat",{
		method: "GET"
	})

	return response
}


export let handleClick=(e:any)=>{};

let usernameContext = createContext<string>("")

const Chat = () => {
	const [username,setUsername] = useState<string>("")
	const [divProp,setDivProp] = useState<any>(null)
	const prevUsername = useRef("")

	let [hash,uname] = window.location.hash.split("#")

	useEffect(()=>{
		prevUsername.current = username
	},[username])

	handleClick = (e) => {
		setUsername(e.target.id)
		if (divProp !== null) {
			divProp.classList.remove("selected")		
		}
		setDivProp(e.target)

		let nextURL = `#${e.target.id}`;
		let nextTitle = 'New chat';
		let nextState = { next: false };

		e.target.classList.add("selected")
		
		// let client = new Client({pathname:`/chat/${e.target.id}`})

		window.history.pushState(nextState, nextTitle, nextURL)
	}

	let client = new Client({pathname:`/chat/${uname}`})
	console.log(client)

	return (
		<div className="chat-app">
			<SideBar />
			<usernameContext.Provider value={uname}>
				<ChatArea />
			</usernameContext.Provider>
		</div>
	)
}

export {usernameContext}
export default Chat;