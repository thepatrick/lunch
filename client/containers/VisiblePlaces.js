import { connect } from 'react-redux';
import PlaceList from '../components/PlaceList';

const getPlaces = (places, placeIds) => {
  return placeIds.map(placeId => places[placeId]);
};

const mapStateToProps = state => (
  {
    places: getPlaces(state.places.placesById, state.places.placeList),
    isFetching: state.places.isFetching,
  }
);

// const mapDispatchToProps = dispatch => (
//   {
//     onTodoClick: (id) => {
//       dispatch(toggleTodo(id));
//     },
//   }
// );

const VisiblePlaceList = connect(
  mapStateToProps,
  // mapDispatchToProps,
)(PlaceList);

export default VisiblePlaceList;
