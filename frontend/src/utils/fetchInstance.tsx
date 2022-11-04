import { cookies as cookie } from './Cookies';
import { Token } from './Types';
import { Config } from './Interfaces';

const BASE_URL: string = "http://127.0.0.1:8000";

const sendRequest = async (url: string,config: Config): Promise<any> => {
	console.log(`${BASE_URL}${url}`)
	let response = await fetch(`${BASE_URL}${url}`, config)
	let data: object = await response.json()
	console.log("Requesting",data)
	return {response,data}
}

export const refreshToken = async (refresh_token: string): Promise<any> => {
	let response: any = await fetch(`${BASE_URL}/account/get-access-token`,{
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			"refresh_token": refresh_token
		})
	})
	let data: Token = await response.json()
	cookie.remove("authTokens")
	cookie.set("authTokens",JSON.stringify({access_token:data.access_token,refresh_token:refresh_token}))

	return data
}

const customFetcher = async (url: string, config: Config): Promise<any> => {
	let authTokens = cookie.get("authTokens")
	console.log("authTokens",authTokens)

	config["headers"] = {
		Authorization: `Bearer ${authTokens.access_token}`
	}

	let {response,data} = await sendRequest(url,config)

	console.log(response.status)
	if (response.status === 401) {
		authTokens = await refreshToken(authTokens.refresh_token)

		config["headers"] = {
			Authorization: `Bearer ${authTokens.access_token}`
		}

		let newResponse = await sendRequest(url,config)

		response = newResponse.response
		data = newResponse.data
	}
	return {response,data}
}

export default customFetcher;