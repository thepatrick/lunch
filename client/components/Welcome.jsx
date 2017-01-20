import React, { PropTypes } from 'react';

const Welcome = ({ name, teamName, isFetching, error }) => {
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
    <h1>{name} / {teamName}</h1>
  );
};

Welcome.propTypes = {
  name: PropTypes.string.isRequired,
  teamName: PropTypes.string.isRequired,
  isFetching: PropTypes.bool.isRequired,
};

export default Welcome;
