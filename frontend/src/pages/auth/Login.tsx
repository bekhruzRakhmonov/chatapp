import React from 'react';
import { cookies } from '../../utils/Cookies';
import { User,Token } from '../../utils/Types';
import {Navigate} from 'react-router-dom';

class Login extends React.Component<any,any> {
	constructor(props: any) {
		super(props);

		console.log(props)

		this.handleSubmit = this.handleSubmit.bind(this)
		this.sendRequet = this.sendRequet.bind(this)
		this.setCookie = this.setCookie.bind(this)

		this.state = {
			error: {
				display: "none",
				message: null,
			},
			username:null,
			password:null,
			success: false,
		}
	}

	setCookie(token:Token) {
		let minutesToAdd = 2000;
		let currentDate = new Date();
		let futureDate = new Date(currentDate.getTime() + minutesToAdd*581212165);
		try {
			cookies.set("authTokens",JSON.stringify(token),{path: "/",expires: futureDate})
			// cookies.set("refresh_token",token.refresh_token,{path: "/"})
		} catch(error) {
			console.log("error",error)
		}
	}	

	async handleSubmit(event: any): Promise<any>{
		event.preventDefault();

		let userData: User = {
			username: this.state.username,
			password: this.state.password,
		}

		try {
			await this.sendRequet(userData)
			this.setState({success: true})
		} catch (status) {
			if (status === 404) {
				this.setState({error: {display:"",message:"User not found"}})
			} else {
				this.setState({error: {display:"",message:"Password is incorrect"}})
			}
		}

	}

	async sendRequet(userData: any): Promise<number>{
		let response = await fetch("http://127.0.0.1:8000/account/login",{
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				username: userData.username,
				password: userData.password,
			})
		})

		let result = await response.json()
		console.log(response)
		if (response.ok) {
			this.setCookie(result)
			return Promise.resolve(response.status)
		} else {
			return Promise.reject(response.status)
		}
	}

	render() {
		return (
			<div className="container">
				<div className="login-form">
					<div className="form">
					    <form onSubmit={this.handleSubmit}>
					        <div className="input-container">
					          	<label>Username </label>
					          	<input type="text" name="uname" onChange={(e)=> {this.setState({username: e.target.value});this.setState({error:{}})}} required />
					          	<p className="errorMessage" style={{display: this.state.error.display}}>{this.state.error.message}</p>
					          	{/*{renderErrorMessage("uname")}*/}
					        </div>
					        <div className="input-container">
					          	<label>Password </label>
					          	<input type="password" name="pass" onChange={(e)=> {this.setState({password: e.target.value});this.setState({error:{}})}} required />
					          	{/*{renderErrorMessage("pass")}*/}
					        </div>
					        <div className="button-container">
					          	<input type="submit" value="Login"/>
					        </div>
					    </form>
					    {this.state.success ? <Navigate to="/chat" /> : ""}
				    </div>
				</div>
			</div>
		)
	}
}

export default Login