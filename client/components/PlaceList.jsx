import React from 'react';
import PropTypes from 'prop-types';
import Place from './Place';

const PlaceList = ({ places, isFetching }) => {
  let body;

  if (isFetching) {
    body = (
      <tr>
        <td colSpan="5">Loading...</td>
      </tr>
    );
  } else if (places.length === 0) {
    body = (
      <tr>
        <td colSpan="5">
          None yet, use <code>/lunch add My Place</code> from slack to add your first place.
        </td>
      </tr>
    );
  } else {
    body = places.map(place =>
      <Place
        key={place.id}
        {...place}
      />,
    );
  }

  return (<table className="table table-striped">
    <thead>
      <tr>
        <th>Name</th>
        <th colSpan="2">Last Visited</th>
        <th colSpan="2">Last Skipped</th>
      </tr>
    </thead>
    <tbody>
      {body}
    </tbody>
  </table>);
};

PlaceList.propTypes = {
  isFetching: PropTypes.bool.isRequired,
  places: PropTypes.arrayOf(PropTypes.shape({
    id: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
    lastVisited: PropTypes.instanceOf(Date),
    lastSkipped: PropTypes.instanceOf(Date),
  }).isRequired).isRequired,
};

export default PlaceList;

/* number, shape, string, bool, arrayOf, func */
/* onClick={() => onTodoClick(todo.id)} */
