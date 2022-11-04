import React from 'react';
import {Outlet,Navigate} from 'react-router-dom';
import {cookies} from './Cookies';

export default class PrivateRoutes extends React.Component<any,any> {
	constructor(props: any){
		super(props);
		console.log(props)
	}

	render() {
		let token: string = cookies.get('authTokens')
		return (
			token ? <Outlet/> : <Navigate to="/login"/>
		)
	}
}