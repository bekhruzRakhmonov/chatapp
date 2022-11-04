import React from 'react';
import { cookies } from '../utils/Cookies'
import { isExpired, decodeToken } from "react-jwt";
import { useNavigate } from "react-router-dom";

export default class Home extends React.Component<any,any> {
	constructor(props: any) {
		super(props);

		this.decode_token = this.decode_token.bind(this)
		this.handleClick = this.handleClick.bind(this)

		this.state = {
			isOnline: false,
		}
	}

	private decode_token(): string {
		let token: any = cookies.get("authTokens")
		let decodedToken: object | any = decodeToken(token.access_token)
		let username: string = decodedToken.iss
		return username
	}

	componentWillUnmount() {
		console.log("componentWillUnmount")
		this.setState({isOnline: true})
	}

	componentDidUpdate(prevProps: any, prevState: any) {
	    console.log('Prev state', prevState); // Before update
	    console.log('New state', this.state); // After update 
	    console.log('Prev props', prevProps)
	}

	componentDidMount() {
		console.log("componentDidMount")
	}

	handleClick() {
		//window.location.replace("/")
		console.log(window.location)
		console.log(window.history)
		// Current URL: https://my-website.com/page_a
		const nextURL = '#sds';
		const nextTitle = 'My new page title';
		const nextState = { next: false };

		// This will create a new entry in the browser's history, without reloading
		window.history.pushState(nextState, nextTitle, nextURL);
		console.log(window.location.pathname)

		// This will replace the current entry in the browser's history, without reloading
		// window.history.replaceState(nextState, nextTitle, nextURL);
		console.log(window.location.hash)
		console.log(window.history)

		console.log(window.history.state.next === false)
		if (window.history.state.next === false){
			console.log("next state")
		}

		window.addEventListener("popstate",()=>{
			console.log("works")

		})

	}


	render() {
		const username:string = this.decode_token()
		return (
			<div className="home">
				<h1>Home page,{username}</h1>
				<button type="button" onClick={this.handleClick}>Click me</button>
			</div>
		)
	}
}