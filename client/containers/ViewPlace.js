import { connect } from 'react-redux';
import ViewPlaceDetails from '../components/ViewPlaceDetails';
import { updatePlaceName } from '../actions';

const mapStateToProps = (state, ownProps) => {
  const { lastVisited, lastSkipped, skipCount, visitCount } = state.places.placesById[ownProps.placeId] || {};
  return {
    isFetching: state.places.isFetching,
    lastVisited,
    lastSkipped,
    skipCount,
    visitCount,
  };
};

const mapDispatchToProps = (dispatch, ownProps) => (
  {
    onSubmit: (form) => {
      dispatch(updatePlaceName(ownProps.placeId, form.name));
    },
  }
);

const EditPlace = connect(
  mapStateToProps,
  mapDispatchToProps,
)(ViewPlaceDetails);

export default EditPlace;
