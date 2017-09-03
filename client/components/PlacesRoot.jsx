import React from 'react';
import { Link } from 'react-router';
import UserWelcome from '../containers/UserWelcome';
import VisiblePlaces from '../containers/VisiblePlaces';

function PlacesRoot() {
  return (
    <div className="row">
      <main className="col-sm-12 col-md-12 pt-5">

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
