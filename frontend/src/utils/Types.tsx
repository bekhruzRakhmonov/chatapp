export type User = {
	username:string;
	password:string;
}

export type Token = {
	access_token:string;
	refresh_token:string;
}

export type EffectCallback = () => (void | Promise<void> | (() => void));