import { combineReducers } from 'redux';
import userSession from './userSessionReducer';
import teams from './teamsReducer';
import competition from "./competitionReducer";

const rootReducer = combineReducers({
  userSession,
  teams,
  competition
});

export default rootReducer;
