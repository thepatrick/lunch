import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router';
import PlaceDate from './PlaceDate';

const ViewPlaceDetails = ({ isFetching, lastVisited, lastSkipped, visitCount, skipCount }) => {
  if (isFetching) {
    return (
      <div>Loading...</div>
    );
  }

  return (
    <form>
      <div className="form-group row">
        <label htmlFor="name-input-field" className="col-sm col-md-2 col-form-label">Last Visited</label>
        <div className="col-sm col-md-10">
          <p className="form-control-static">
            <PlaceDate
              date={lastVisited}
              defaultString="Never"
            /> ({visitCount} times)
          </p>
        </div>
      </div>
      <div className="form-group row">
        <label htmlFor="name-input-field" className="col-sm col-md-2 col-form-label">Last Skipped</label>
        <div className="col-sm col-md-10">
          <p className="form-control-static">
            <PlaceDate
              date={lastSkipped}
              defaultString="Never"
            /> ({skipCount} times)
          </p>
        </div>
      </div>
    </form>
  );
};

ViewPlaceDetails.propTypes = {
  isFetching: PropTypes.bool.isRequired,
  lastVisited: PropTypes.instanceOf(Date),
  lastSkipped: PropTypes.instanceOf(Date),
  skipCount: PropTypes.number.isRequired,
  visitCount: PropTypes.number.isRequired,
};

ViewPlaceDetails.defaultProps = {
  lastVisited: undefined,
  lastSkipped: undefined,
  visitCount: 0,
  skipCount: 0,
};

export default ViewPlaceDetails;
