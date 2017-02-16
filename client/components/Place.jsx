import React, { PropTypes } from 'react';
import { Link } from 'react-router';
import PlaceDate from './PlaceDate';

const Place = ({ editUrl, name, lastVisited, lastSkipped }) => (
  <tr>
    <td>
      <Link to={editUrl}>
        {name}
      </Link>
    </td>
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
  editUrl: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
  lastVisited: PropTypes.instanceOf(Date),
  lastSkipped: PropTypes.instanceOf(Date),
};

Place.defaultProps = {
  lastVisited: undefined,
  lastSkipped: undefined,
};

export default Place;
