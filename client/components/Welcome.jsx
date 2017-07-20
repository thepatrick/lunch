import React from 'react';
import PropTypes from 'prop-types';

const Welcome = ({ teamName, isFetching, error }) => {
  if (isFetching) {
    return (
      <h1>Loading...</h1>
    );
  }
  if (error) {
    return (
      <h1>{error.message || error}</h1>
    );
  }
  return (
    <h1>{teamName}&#8217;s Lunch Places</h1>
  );
};

Welcome.propTypes = {
  teamName: PropTypes.string.isRequired,
  isFetching: PropTypes.bool.isRequired,
  error: PropTypes.instanceOf(Error),
};

Welcome.defaultProps = {
  error: undefined,
};

export default Welcome;
