const updateIndividualPlace = (state = {}, action) => {
  switch (action.type) {
    case 'SAVE_PLACE':
      return Object.assign({}, state, { isSaving: true, error: null });
    case 'SAVE_PLACE_SUCCESS':
      return Object.assign({}, state, {
        isSaving: false,
        error: null,
        name: action.placeName,
      });
    case 'SAVE_PLACE_FAILURE':
      return Object.assign({}, state, { isSaving: false, error: action.error });
    case 'DELETE_PLACE':
      return Object.assign({}, state, { isDeleting: true });
    case 'DELETE_PLACE_FAILURE':
      return Object.assign({}, state, { isDeleting: false, error: action.error });
    default:
      return state;
  }
};

const placesById = (state = {}, action) => {
  switch (action.type) {
    case 'FETCH_PLACES_SUCCESS':
      return action.response.reduce((previous, current) => {
        const newItem = {};
        newItem[current.id] = {
          id: current.id,
          name: current.name,
          lastSkipped: current.last_skipped && new Date(current.last_skipped),
          lastVisited: current.last_visited && new Date(current.last_visited),
          skipCount: current.skip_count,
          visitCount: current.visit_count,
          isSaving: false,
        };
        return Object.assign({}, previous, newItem);
      }, {});
    case 'SAVE_PLACE':
    case 'SAVE_PLACE_SUCCESS':
    case 'SAVE_PLACE_FAILURE':
    case 'DELETE_PLACE':
    case 'DELETE_PLACE_FAILURE': {
      const updatedPlace = {};
      updatedPlace[action.id] = updateIndividualPlace(state[action.id], action);
      return Object.assign({}, state, updatedPlace);
    }
    case 'DELETE_PLACE_SUCCESS': {
      const updatedPlace = {};
      updatedPlace[action.id] = undefined;
      return Object.assign({}, state, updatedPlace);
    }
    default:
      return state;
  }
};

export default placesById;
