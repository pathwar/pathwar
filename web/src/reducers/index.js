import { combineReducers } from 'redux';
import session from './sessionReducer';
import teams from './teamsReducer';

const rootReducer = combineReducers({
  session,
  teams
});

export default rootReducer;
