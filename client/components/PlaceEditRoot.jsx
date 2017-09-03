import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router';
import EditPlace from '../containers/EditPlace';
import ViewPlace from '../containers/ViewPlace';
import DeletePlaceContainer from '../containers/DeletePlaceContainer';

function PlaceEditRoot({ params }) {
  return (
    <div className="row">
      <main className="col-sm-12 col-md-12 pt-5">
        <DeletePlaceContainer placeId={params.id}>Delete</DeletePlaceContainer>
        <h1>Place</h1>
        <div className="table-responsive">
          <EditPlace placeId={params.id} />
          <ViewPlace placeId={params.id} />
        </div>
      </main>
    </div>
  );
}

PlaceEditRoot.propTypes = {
  params: PropTypes.shape({
    id: PropTypes.string.isRequired,
  }).isRequired,
};

export default PlaceEditRoot;
