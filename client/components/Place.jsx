import React, { PropTypes } from 'react';
import PlaceDate from './PlaceDate';

const Place = ({ name, lastVisited, lastSkipped }) => (
  <tr>
    <td>{name}</td>
    <td>
      <PlaceDate
        date={lastVisited}
        defaultString="Never"
      />
    </td>
    <td>
      <PlaceDate
        date={lastSkipped}
        defaultString="Never"
      />
    </td>
  </tr>
);

Place.propTypes = {
  name: PropTypes.string.isRequired,
  lastVisited: PropTypes.instanceOf(Date),
  lastSkipped: PropTypes.instanceOf(Date),
};

export default Place;
