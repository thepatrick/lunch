import React from 'react';
import PropTypes from 'prop-types';
import { Provider } from 'react-redux';
import { Router, Route, IndexRoute } from 'react-router';
import AppContainer from '../containers/AppContainer';
import PlacesRoot from './PlacesRoot';
import PlaceEditRoot from './PlaceEditRoot';

import ConnectedIntlProvider from '../containers/ConnectedIntlProvider';

const Root = ({ store, history }) => (
  <Provider store={store}>
    <ConnectedIntlProvider>
      <Router history={history}>
        <Route path="/manage" component={AppContainer}>
          <IndexRoute component={PlacesRoot} />
          <Route path="places/:id" component={PlaceEditRoot} />
        </Route>
      </Router>
    </ConnectedIntlProvider>
  </Provider>
);

Root.propTypes = {
  store: PropTypes.shape({}).isRequired,
  history: PropTypes.shape({}).isRequired,
};

export default Root;
