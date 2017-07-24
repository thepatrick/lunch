const placeList = (state = [], action) => {
  switch (action.type) {
    case 'FETCH_PLACES_SUCCESS':
      return action.response.map(place => place.id);
    case 'DELETE_PLACE_SUCCESS':
      return state.filter(place => place !== action.id);
    default:
      return state;
  }
};

export default placeList;
