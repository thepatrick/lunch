import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import UserWelcome from '../containers/UserWelcome';
import VisiblePlaces from '../containers/VisiblePlaces';

function PlacesRoot() {
  return (
    <div className="row">
      <nav className="col-sm-3 col-md-2 hidden-xs-down bg-faded sidebar">
        <ul className="nav nav-pills flex-column">
          <li className="nav-item">
            <Link className="nav-link active" to="/manage">
              Overview <span className="sr-only">(current)</span>
            </Link>
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
  );
}

PlacesRoot.propTypes = {

};

export default PlacesRoot;
