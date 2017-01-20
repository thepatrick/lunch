import React, { PropTypes } from 'react';
import UserWelcome from '../containers/UserWelcome';
import VisiblePlaces from '../containers/VisiblePlaces';

const App = ({ userFetching, userError }) => {
  if (userFetching) {
    return <div>Loading...</div>;
  }
  if (userError) {
    return <div>{userError.message}</div>;
  }
  return (<div>
    <nav className="navbar navbar-toggleable-md navbar-inverse fixed-top bg-inverse">
      <button
        className="navbar-toggler navbar-toggler-right hidden-lg-up"
        type="button"
        data-toggle="collapse"
        data-target="#navbarsExampleDefault"
        aria-controls="navbarsExampleDefault"
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span className="navbar-toggler-icon" />
      </button>
      <a className="navbar-brand" href="/">Lunch Bot</a>

      <div className="collapse navbar-collapse" id="navbarsExampleDefault">
        <ul className="navbar-nav mr-auto">
          <li className="nav-item active">
            <a className="nav-link" href="/manage/">
              Places <span className="sr-only">(current)</span>
            </a>
          </li>
          <li className="nav-item">
            <a className="nav-link" href="/manage/api/logout">Logout</a>
          </li>
        </ul>
      </div>
    </nav>

    <div className="container-fluid">
      <div className="row">
        <nav className="col-sm-3 col-md-2 hidden-xs-down bg-faded sidebar">
          <ul className="nav nav-pills flex-column">
            <li className="nav-item">
              <a className="nav-link active" href="/manage/">
                Overview <span className="sr-only">(current)</span>
              </a>
            </li>
          </ul>
        </nav>

        <main className="col-sm-9 offset-sm-3 col-md-10 offset-md-2 pt-3">

          <UserWelcome />
          <div className="table-responsive">
            <VisiblePlaces />
          </div>
        </main>
      </div>
    </div>

  </div>);
};


App.propTypes = {
  userError: PropTypes.instanceOf(Error),
  userFetching: PropTypes.bool.isRequired,
};

export default App;
