import { connect } from 'react-redux';
import { deletePlace } from '../actions';
import DeletePlace from '../components/DeletePlace';

const mapStateToProps = (state, ownProps) => {
  const place = state.places.placesById[ownProps.placeId];
  return {
    active: (place && !place.isDeleting) || false,
  };
};

const mapDispatchToProps = (dispatch, ownProps) => (
  {
    onClick: () => {
      dispatch(deletePlace(ownProps.placeId));
    },
  }
);

const DeletePlaceContainer = connect(
  mapStateToProps,
  mapDispatchToProps,
)(DeletePlace);

export default DeletePlaceContainer;
