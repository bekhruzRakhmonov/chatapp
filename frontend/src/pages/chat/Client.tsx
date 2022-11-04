import React from 'react';
// @ts-ignore
import { w3cwebsocket as WebSocket } from 'websocket';
import { cookies } from '../../utils/Cookies';
import { useWinLocation } from '../../utils/CustomHooks';
import { Token } from '../../utils/Types';
import { refreshToken } from '../../utils/fetchInstance';
import { getUsername } from '../../utils/Cookies';

// I also should remove question mark from pathname

interface Props {
	pathname: string;
}

export let sendMessages=(e:any)=>{};

export default class Client {
	private chatMessageArea: any;
	private chatMessages: any;

	constructor(props: Props) {
		this.connectToWebsocket = this.connectToWebsocket.bind(this)
		if (props.pathname !== ""){
			this.connectToWebsocket(props.pathname)
		}
	}


	private connectToWebsocket(username: string) {
		const authTokens = cookies.get("authTokens")
		const accessToken: string = authTokens.access_token
		const encodedAccessToken: string = window.btoa(accessToken)
		const client = new WebSocket(`ws://127.0.0.1:8000${username}?accessToken=${encodedAccessToken}`)

		client.onopen = () => {
			console.log("Socket connected successfully")
			function sendMessage() {
			    if (client.readyState === client.OPEN) {
			        client.send("Hello");
			        setTimeout(sendMessage, 1000);
			        console.log("Message sent")
			    }
			}
			// client.send(JSON.stringify({message:"hello"}))
		}

		client.onclose = (e: any) => {
			console.log("Socket closed unexpectedly",e)
		}

		client.onmessage = (e: any) => {
			let response: any = JSON.parse(e.data)
			console.log("RESPONSE",response)
			switch (response.Send.Status) {
				case 401:
					console.log("Unauthorized user")
					let authTokens: Token = cookies.get("authTokens")
					refreshToken(authTokens.refresh_token)
					this.connectToWebsocket(username)
					break
				case 404:
					console.log("User not found")
					client.close()
					break
				case 200:
					console.log("CREATED")
					if (this.chatMessages !== undefined) {
						if (response.Outbound === getUsername()) {
							this.chatMessages.innerHTML += `<div class="outbound-message" key=4><p id="message">${response.Send.Message}</p></div>`
						} else if (response.Inbound === getUsername()) {
							this.chatMessages.innerHTML += `<div class="inbound-message" key=4><p id="message">${response.Send.Message}</p></div>`
						}
						this.chatMessageArea.scrollTop = this.chatMessageArea.scrollHeight
					}
					break
			}
		}

		sendMessages = (e) => {
			e.preventDefault();
			let text: string = e.target.messageText.value;

			if (this.chatMessages === undefined) {
				let parentChatArea = e.nativeEvent.path[2]
				let childChatArea = parentChatArea.childNodes[0]
				this.chatMessageArea = childChatArea.childNodes[2]
				this.chatMessages = this.chatMessageArea.childNodes[0]
			}
			client.send(JSON.stringify({message: text}))
			console.log(e.target.messageText.value)
			e.target.messageText.value = "";
		}

	}
}