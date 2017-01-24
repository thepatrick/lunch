import { connect } from 'react-redux';
import PlaceList from '../components/PlaceList';

const getPlaces = (places, placeIds) => (
  placeIds.map(placeId => places[placeId]).map(place => (
    Object.assign({}, place, { editUrl: `/manage/places/${place.id}` })
  ))
);

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
