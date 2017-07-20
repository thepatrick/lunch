import { connect } from 'react-redux';
import EditPlaceForm from '../components/EditPlaceForm';
import { updatePlaceName } from '../actions';

const mapStateToProps = (state, ownProps) => {
  const place = state.places.placesById[ownProps.placeId];
  return {
    initialValues: place,
    isFetching: state.places.isFetching,
    isSaving: !!(place && place.isSaving),
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
)(EditPlaceForm);

export default EditPlace;
