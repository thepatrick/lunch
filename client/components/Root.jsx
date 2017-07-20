import React from 'react';
import PropTypes from 'prop-types';
import { Provider } from 'react-redux';
import { Router, Route, IndexRoute, browserHistory } from 'react-router';
// import App from './App';
import AppContainer from '../containers/AppContainer';
import PlacesRoot from './PlacesRoot';
import PlaceEditRoot from './PlaceEditRoot';

import ConnectedIntlProvider from '../containers/ConnectedIntlProvider';

const Root = ({ store }) => (
  <Provider store={store}>
    <ConnectedIntlProvider>
      <Router history={browserHistory}>
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
};

export default Root;
