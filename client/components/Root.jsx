import React, { PropTypes } from 'react';
import { Provider } from 'react-redux';
import { Router, Route, IndexRoute, browserHistory } from 'react-router';
// import App from './App';
import AppContainer from '../containers/AppContainer';
import PlacesRoot from './PlacesRoot';

import ConnectedIntlProvider from '../containers/ConnectedIntlProvider';

const Root = ({ store }) => (
  <Provider store={store}>
    <ConnectedIntlProvider>
      <Router history={browserHistory}>
        <Route path="/manage" component={AppContainer}>
          <IndexRoute component={PlacesRoot} />
        </Route>
      </Router>
    </ConnectedIntlProvider>
  </Provider>
);

Root.propTypes = {
  store: PropTypes.shape({}).isRequired,
};

export default Root;
