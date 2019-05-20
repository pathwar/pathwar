import { combineReducers } from 'redux';
import userSession from './userSessionReducer';
import teams from './teamsReducer';
import tournaments from "./tournamentReducer";

const rootReducer = combineReducers({
  userSession,
  teams,
  tournaments
});

export default rootReducer;
