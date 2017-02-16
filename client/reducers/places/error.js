const error = (state = null, action) => {
  switch (action.type) {
    case 'FETCH_PLACES':
      return null;
    case 'FETCH_PLACES_FAILURE':
      return action.error;
    default:
      return state;
  }
};

export default error;
