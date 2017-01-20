import React from 'react';
import thunkMiddleware from 'redux-thunk';
import createLogger from 'redux-logger';
import { render } from 'react-dom';
import { Provider } from 'react-redux';
import { createStore, applyMiddleware } from 'redux';
import todoApp from './reducers';
import AppContainer from './containers/AppContainer';
import { fetchCurrentUser, fetchAllPlaces } from './actions';

import ConnectedIntlProvider from './containers/ConnectedIntlProvider';

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
  <Provider store={store}>
    <ConnectedIntlProvider>
      <AppContainer />
    </ConnectedIntlProvider>
  </Provider>,
  root,
);
