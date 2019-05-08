import { combineReducers } from 'redux';
import session from './sessionReducer';
import teams from './teamsReducer';
import competitions from "./competitionsReducer";

const rootReducer = combineReducers({
  session,
  teams,
  competitions
});

export default rootReducer;
