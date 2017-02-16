const placeList = (state = [], action) => {
  switch (action.type) {
    case 'FETCH_PLACES_SUCCESS':
      return action.response.map(place => place.id);
    default:
      return state;
  }
};

export default placeList;
