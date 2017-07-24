import React from 'react';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import { render } from 'react-dom';
import { createStore, applyMiddleware } from 'redux';

import { browserHistory } from 'react-router';
import { routerMiddleware, syncHistoryWithStore } from 'react-router-redux';

import reducers from './reducers';
import { fetchCurrentUser, fetchAllPlaces } from './actions';

import Root from './components/Root';

require('bootstrap-loader');
require('whatwg-fetch');

require('./dashboard.css');

const loggerMiddleware = createLogger();

const middleware = [thunkMiddleware, routerMiddleware(browserHistory)];
if (process.env.NODE_ENV === 'development') {
  middleware.push(loggerMiddleware);
}

const store = createStore(
  reducers,
  applyMiddleware(...middleware),
);

// Create an enhanced history that syncs navigation events with the store
const history = syncHistoryWithStore(browserHistory, store);

store.dispatch(fetchCurrentUser());
store.dispatch(fetchAllPlaces());

const root = document.createElement('div');
document.body.appendChild(root);

render(
  <Root store={store} history={history} />,
  root,
);
