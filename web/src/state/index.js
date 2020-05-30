import { combineReducers } from "redux";
import userSession from "./userSessionReducer";
import organizations from "./organizationsReducer";
import seasons from "./seasonReducer";

const rootReducer = combineReducers({
  userSession,
  organizations,
  seasons
});

export default rootReducer;
