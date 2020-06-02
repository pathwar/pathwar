/* eslint-disable react/display-name */
import React from "react";
import { Provider } from "react-redux";
import thunk from "redux-thunk";
import {
  createStore as reduxCreateStore,
  compose,
  applyMiddleware,
} from "redux";
import { ThemeProvider } from "emotion-theming";
import { lightTheme, darkTheme } from "../styles/themes";
import rootReducer from ".";

const windowExist = typeof window === "object";

const composeEnhancers =
  windowExist && window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__
    ? window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__
    : compose;

const middlewares = [thunk];

const createStore = () =>
  reduxCreateStore(
    rootReducer,
    composeEnhancers(applyMiddleware(...middlewares))
  );

if (!window.activeTheme) {
  window.activeTheme = "light";
}

const themeToUse = window.activeTheme === "dark" ? darkTheme : lightTheme;

// eslint-disable-next-line react/display-name
export default ({ element }) => (
  <ThemeProvider theme={themeToUse}>
    <Provider store={createStore()}>{element}</Provider>
  </ThemeProvider>
);
