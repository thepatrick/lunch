import React, { PropTypes } from 'react';
import NavBar from '../containers/NavBar';

const App = ({ userFetching, userError, children }) => {
  if (userFetching) {
    return (<div>
      <NavBar />
    </div>);
  }
  if (userError) {
    return <div>{userError.message}</div>;
  }
  return (<div>
    <NavBar />

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
