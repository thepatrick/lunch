const defaultState = {
  isFetching: false,
  error: null,
  placesById: {},
  placeList: [],
};

const placesById = places => (
  places.reduce((previous, current) => {
    const newItem = {};
    newItem[current.id] = {
      id: current.id,
      name: current.name,
      lastSkipped: current.last_skipped && new Date(current.last_skipped),
      lastVisited: current.last_visited && new Date(current.last_visited),
    };
    return Object.assign({}, previous, newItem);
  }, {})
);

const placeList = places => (
  places.map(place => place.id)
);

const user = (state = defaultState, action) => {
  switch (action.type) {
    case 'FETCH_PLACES':
      return Object.assign({}, state, {
        isFetching: true,
        error: null,
      });
    case 'FETCH_PLACES_SUCCESS':
      return Object.assign({}, state, {
        isFetching: false,
        error: null,
        placesById: placesById(action.response),
        placeList: placeList(action.response),
      });
    case 'FETCH_PLACES_FAILURE':
      return Object.assign({}, state, {
        isFetching: false,
        error: action.error,
      });
    default:
      return state;
  }
};

export default user;

// const todo = (state = {}, action) => {
//   switch (action.type) {
//     case 'ADD_TODO':
//       return {
//         id: action.id,
//         text: action.text,
//         completed: false,
//       }
//     case 'TOGGLE_TODO':
//       if (state.id !== action.id) {
//         return state;
//       }

//       return Object.assign({}, state, {
//         completed: !state.completed,
//       });

//     default:
//       return state;
//   }
// };

// const todos = (state = [], action) => {
//   switch (action.type) {
//     case 'ADD_TODO':
//       return [
//         ...state,
//         todo(undefined, action),
//       ];
//     case 'TOGGLE_TODO':
//       return state.map(t => todo(t, action));
//     default:
//       return state;
//   }
// };

// export default todos;
