import { combineReducers } from 'redux';
import userSession from './userSessionReducer';
import teams from './teamsReducer';
import seasons from "./seasonReducer";

const rootReducer = combineReducers({
  userSession,
  teams,
  seasons
});

export default rootReducer;
