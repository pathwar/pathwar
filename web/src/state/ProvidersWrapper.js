/* eslint-disable react/display-name */
import React from "react";
import { Provider } from "react-redux";
import thunk from "redux-thunk";
import {
  createStore as reduxCreateStore,
  compose,
  applyMiddleware,
} from "redux";
import { IntlProvider } from "react-intl";
import rootReducer from ".";

import messages_fr from "../translations/fr.json";
import messages_en from "../translations/en.json";

const messages = {
  fr: messages_fr,
  en: messages_en,
};

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

let lang;
const browser = typeof window !== "undefined" && window;
if (browser) {
  const storedLang = window.localStorage.getItem("pw.lang");
  if (storedLang) {
    lang = storedLang;
  } else {
    window.localStorage.setItem("pw.lang", "en");
    lang = "en";
  }
}

// eslint-disable-next-line react/display-name
export default ({ element }) => (
  <Provider store={createStore()}>
    <IntlProvider locale={lang} messages={messages[lang]}>
      {element}
    </IntlProvider>
  </Provider>
);
