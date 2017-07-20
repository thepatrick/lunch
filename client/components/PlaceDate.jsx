import React from 'react';
import PropTypes from 'prop-types';
import { FormattedDate } from 'react-intl';

const PlaceDate = ({ date, defaultString }) => {
  if (date) {
    return (<FormattedDate
      value={date}
      day="numeric"
      month="long"
      year="numeric"
    />);
  }

  return <span>{defaultString}</span>;
};

PlaceDate.propTypes = {
  defaultString: PropTypes.string.isRequired,
  date: PropTypes.instanceOf(Date),
};

PlaceDate.defaultProps = {
  date: undefined,
};

export default PlaceDate;
