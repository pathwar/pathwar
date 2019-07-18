/* eslint-disable no-unused-vars */
import Cookies from "js-cookie";
import { 
	PERFORM_LOGIN, 
	LOGIN_FAILED,
	SET_USER_SESSION,
	PING_USER_SUCCESS,
	PING_USER_FAILED
} from "../constants/actionTypes"
import { USER_SESSION_TOKEN_NAME } from "../constants/userSession";
import { performLogin, pingUser } from "../api/userSession"
import { setActiveTeam as setActiveTeamAction } from "./teams";
import { setActiveTournament as setActiveTournamentAction } from "./tournaments"


export const performLoginAction = (email, password) => async dispatch => {
	dispatch({
		type: PERFORM_LOGIN
	});

	try {
		const response = await performLogin(email, password);
		const { userSession, token, lastActiveTeam, defaultTournament } = response.data;
		
		dispatch({
			type: SET_USER_SESSION,
			payload: { activeUser: userSession }
		});

		Cookies.set(USER_SESSION_TOKEN_NAME, token)

		dispatch(setActiveTeamAction(lastActiveTeam))
		dispatch(setActiveTournamentAction(defaultTournament));


	} catch (error) {
		dispatch({ type: LOGIN_FAILED, payload: { error } });
	}
};

export const pingUserAction = () => async dispatch => {

	try {
		const response = await pingUser();
		const { isAuthenticated, token, userSession, lastActiveTeam, defaultTournament } = response.data;
		dispatch({
			type: PING_USER_SUCCESS,
			payload: { 
				isAuthenticated: isAuthenticated,
				activeUser: userSession
			}
		});

		Cookies.set(USER_SESSION_TOKEN_NAME, token)

		dispatch(setActiveTeamAction(lastActiveTeam))
		dispatch(setActiveTournamentAction(defaultTournament));

	} catch (error) {
		dispatch({ type: PING_USER_FAILED, payload: { error } });
	}
};