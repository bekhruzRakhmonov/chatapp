import { useEffect, useState } from 'react';


export function useWinLocation(location:any): string {
	const [loc,setLoc] = useState(null);

	useEffect(()=>{
		setLoc(location.pathname)
	},[location]);

	return location
}