const isFetching = (state = false, action) => {
  switch (action.type) {
    case 'FETCH_PLACES':
      return true;
    case 'FETCH_PLACES_FAILURE':
    case 'FETCH_PLACES_SUCCESS':
      return false;
    default:
      return state;
  }
};

export default isFetching;
