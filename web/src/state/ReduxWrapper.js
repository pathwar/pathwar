import React from 'react';
import { Provider } from 'react-redux';
import thunk from 'redux-thunk';
import { createStore as reduxCreateStore, compose, applyMiddleware } from 'redux';
import rootReducer from '.';

const windowExist = typeof window === 'object';

const composeEnhancers = windowExist && window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__
? window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__
: compose;

const middlewares = [
    thunk
  ];

const createStore = () => reduxCreateStore(
    rootReducer,
    composeEnhancers(applyMiddleware(...middlewares))
);

export default ({ element }) => (
  <Provider store={createStore()}>{element}</Provider>
);
