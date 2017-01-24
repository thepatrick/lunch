import React from 'react';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import { render } from 'react-dom';
import { createStore, applyMiddleware } from 'redux';
import todoApp from './reducers';
import { fetchCurrentUser, fetchAllPlaces } from './actions';

import Root from './components/Root';

require('bootstrap-loader');
require('whatwg-fetch');

require('./dashboard.css');

const loggerMiddleware = createLogger();

const store = createStore(
  todoApp,
  applyMiddleware(
    thunkMiddleware,
    loggerMiddleware,
  ),
);

store.dispatch(fetchCurrentUser());
store.dispatch(fetchAllPlaces());

const root = document.createElement('div');
document.body.appendChild(root);

render(
  <Root store={store} />,
  root,
);
