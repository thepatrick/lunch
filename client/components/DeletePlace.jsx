import React from 'react';
import PropTypes from 'prop-types';

const DeletePlace = ({ active, onClick, children }) => (
  <button
    className="btn btn-danger float-right"
    onClick={(e) => {
      e.preventDefault();
      onClick();
    }}
    disabled={!active}
  >
    {children}
  </button>
  );

DeletePlace.propTypes = {
  active: PropTypes.bool.isRequired,
  children: PropTypes.node.isRequired,
  onClick: PropTypes.func.isRequired,
};

export default DeletePlace;
