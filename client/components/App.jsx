import React, { PropTypes } from 'react';
import Nav from './Nav';

const App = ({ userFetching, userError, children }) => {
  if (userFetching) {
    return <div>Loading...</div>;
  }
  if (userError) {
    return <div>{userError.message}</div>;
  }
  return (<div>
    <Nav />

    <div className="container-fluid">
      {children}
    </div>

  </div>);
};


App.propTypes = {
  userError: PropTypes.instanceOf(Error),
  userFetching: PropTypes.bool.isRequired,
  children: PropTypes.node.isRequired,
};

App.defaultProps = {
  userError: undefined,
};

export default App;
