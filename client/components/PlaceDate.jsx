import React, { PropTypes } from 'react';
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

export default PlaceDate;
