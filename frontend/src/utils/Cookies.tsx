import Cookies from 'universal-cookie';
import { decodeToken } from "react-jwt";

export const cookies = new Cookies()

export const getUsername = (): string => {
    let token: any = cookies.get("authTokens")
    let decodedToken: object | any = decodeToken(token.access_token)
    let username: string = decodedToken.iss
    return username
}